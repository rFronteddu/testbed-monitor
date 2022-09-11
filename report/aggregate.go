package report

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"testbed-monitor/thresholds"
	"time"
)

type Aggregate struct {
	statusChan      chan *StatusReport
	aggregatePeriod int
	aggregateHour   int
	apiIP           string
	apiPort         string
}

type reportData struct {
	tower          string `json:"tower"`
	arduinoReached string `json:"arduinoReached"`
	towerReached   string `json:"towerReached"`
	bootTime       string `json:"bootTime"`
	reboots        string `json:"reboots"`
	usedRAM        string `json:"usedRAM"`
	usedDisk       string `json:"usedDisk"`
	cpu            string `json:"cpu"`
	reachable      bool   `json:"reachable"`
	temperature    string `json:"temperature"`
}

type reportTemplate struct {
	reportType string       `json:"type"`
	timestamp  string       `json:"timestamp"`
	report     []reportData `json:"report"`
}

type trigger struct {
	field    string
	operator string
	trigger  int
}

func NewAggregate(statusChan chan *StatusReport, aggregatePeriod int, aggregateHour int, apiIP string, apiPort string) *Aggregate {
	aggregate := new(Aggregate)
	aggregate.statusChan = statusChan
	aggregate.aggregatePeriod = aggregatePeriod
	aggregate.aggregateHour = aggregateHour
	aggregate.apiIP = apiIP
	aggregate.apiPort = apiPort
	return aggregate
}

var aggregatedReport reportData
var emailData reportTemplate
var subject string
var unreachableFlag bool

func (aggregate *Aggregate) Start(iPs *[]string) {
	dailyTicker := time.NewTicker(time.Duration(aggregate.aggregatePeriod) * time.Minute)
	reportAggregate := make(map[time.Time]*StatusReport)
	var i int
	for {
		select {
		case <-dailyTicker.C:
			if time.Now().Hour() == aggregate.aggregateHour {
				emailData.report = nil
				unreachableFlag = false
				if time.Now().Weekday().String() == "Sunday" { // Weekly report on Sundays
					for i = 0; i < len(*iPs); i++ {
						weeklyAggregator(reportAggregate, (*iPs)[i], &aggregatedReport)
						emailData.report = append(emailData.report, aggregatedReport)
					}
					emailData.reportType = "Weekly"
				} else { // Daily report
					for i = 0; i < len(*iPs); i++ {
						dailyAggregator(reportAggregate, (*iPs)[i], &aggregatedReport)
						emailData.report = append(emailData.report, aggregatedReport)
					}
					emailData.reportType = "Daily"
				}
				emailData.timestamp = time.Now().Format("Jan 02 2006 15:04:05")
				subject = emailData.reportType + " Testbed Status Report " + emailData.timestamp
				if unreachableFlag {
					subject += ": 1 or more towers is down!"
				}
				MailReport(subject, emailData)
				aggregate.postStatusToApp(emailData)
			}
		case msg := <-aggregate.statusChan:
			reportAggregate[time.Now()] = msg
			if !msg.Reachable {
				aggregate.towerAlertInApp(msg.TowerIP)
			}
			//for _, trigger := range aggregate.threshold {
			//	log.Println(trigger)
			//	// if field operator trigger {}
			//}
		}
	}
}

func (aggregate *Aggregate) SetTriggers(threshold []thresholds.Config) {
	thresholds := make(map[string]trigger)
	fields := reflect.VisibleFields(reflect.TypeOf(struct{ reportData }{}))
	log.Println("Thresholds:")
	for _, field := range fields {
		for _, t := range threshold {
			if field.Name == t.Field {
				setTrigger := trigger{}
				setTrigger.field = t.Field
				setTrigger.operator = t.Operator
				setTrigger.trigger, _ = strconv.Atoi(t.Trigger)
				thresholds[field.Name] = setTrigger
				log.Printf("%s %s %s\n", t.Field, t.Operator, t.Trigger)
			}
		}
	}
}

func (aggregate *Aggregate) postStatusToApp(emailData reportTemplate) {
	apiAddress := aggregate.apiIP + ":" + aggregate.apiPort
	jsonReport, errJ := json.Marshal(emailData.report)
	log.Println(string(jsonReport))
	if errJ != nil {
		log.Println("Error creating json object: ", errJ)
	}
	_, err := http.Post("http://"+apiAddress+"/towers", "application/json", bytes.NewBuffer(jsonReport))
	if err != nil {
		log.Println("Error posting to API", err)
	}
}

func (aggregate *Aggregate) towerAlertInApp(alertIP string) {
	apiAddress := aggregate.apiIP + ":" + aggregate.apiPort
	_, err := http.Get("http://" + apiAddress + "/alert/" + alertIP)
	if err != nil {
		log.Println("Error posting to API", err)
	}
}

func dailyAggregator(reports map[time.Time]*StatusReport, iP string, templateData *reportData) {
	var usedRAMAvg, ramCounter, usedDiskAvg, diskCounter, cpuAvg, cpuCounter, rebootCounter int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk, maxTemp int64 = 0, 0, 0
	compareTime := time.Time{}
	templateData.tower = iP
	for key, element := range reports {
		if key.Day() == time.Now().Day() && element.TowerIP == iP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.reachable = true
				if element.Reachable == true {
					templateData.arduinoReached = element.LastArduinoReachableTimestamp
					templateData.towerReached = element.LastTowerReachableTimestamp
					templateData.bootTime = element.BootTimestamp
				}
				if totalRAM == 0 {
					totalRAM = element.RAMTotal
				}
				if totalDisk == 0 {
					totalDisk = element.DiskTotal
				}
				if !element.Reachable {
					templateData.reachable = false
				}
				compareTime = element.Timestamp.AsTime()
			}
		}
		rebootCounter = rebootCounter + element.RebootsCurrentDay
		usedRAMAvg = usedRAMAvg + element.RAMUsed
		ramCounter++
		usedDiskAvg = usedDiskAvg + element.DiskUsed
		diskCounter++
		cpuAvg = cpuAvg + element.CPUAvg
		cpuCounter++
		if element.Timestamp.AsTime().Day() == time.Now().Day() && element.TowerIP == iP {
			if element.Temperature > 0 {
				if element.Temperature > maxTemp {
					maxTemp = element.Temperature
				}
			}
		}
	}
	if ramCounter > 0 {
		usedRAMAvg = usedRAMAvg / ramCounter
	}
	templateData.usedRAM = strconv.FormatInt(usedRAMAvg, 10) + "/" + strconv.FormatInt(totalRAM, 10) + " MB"
	if diskCounter > 0 {
		usedDiskAvg = usedDiskAvg / diskCounter
	}
	templateData.usedDisk = strconv.FormatInt(usedDiskAvg, 10) + "/" + strconv.FormatInt(totalDisk, 10) + " GB"
	if cpuCounter > 0 {
		cpuAvg = cpuAvg / cpuCounter
	}
	templateData.cpu = strconv.FormatInt(cpuAvg, 10) + "%"
	templateData.reboots = strconv.FormatInt(rebootCounter, 10)
	if templateData.reachable == false {
		unreachableFlag = true
	}

}

func weeklyAggregator(reports map[time.Time]*StatusReport, iP string, templateData *reportData) {
	var usedRAMAvg, ramCounter, usedDiskAvg, diskCounter, cpuAvg, cpuCounter, rebootCounter int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk, maxTemp int64 = 0, 0, 0
	compareTime := time.Time{}
	templateData.tower = iP
	for key, element := range reports {
		templateData.tower = iP
		if key.After(time.Now().Add(-7*24*time.Hour)) && element.TowerIP == iP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.reachable = true
				if element.Reachable == true {
					templateData.arduinoReached = element.LastArduinoReachableTimestamp
					templateData.towerReached = element.LastTowerReachableTimestamp
					templateData.bootTime = element.BootTimestamp
				}
				if totalRAM == 0 {
					totalRAM = element.RAMTotal
				}
				if totalDisk == 0 {
					totalDisk = element.DiskTotal
				}
				if !element.Reachable {
					templateData.reachable = false
				}
				compareTime = element.Timestamp.AsTime()
			}
			rebootCounter = rebootCounter + element.RebootsCurrentDay
			usedRAMAvg = usedRAMAvg + element.RAMUsed
			ramCounter++
			usedDiskAvg = usedDiskAvg + element.DiskUsed
			diskCounter++
			cpuAvg = cpuAvg + element.CPUAvg
			cpuCounter++
			if element.Timestamp.AsTime().After(time.Now().Add(-7*24*time.Hour)) && element.TowerIP == iP {
				if element.Temperature > 0 {
					if element.Temperature > maxTemp {
						maxTemp = element.Temperature
					}
				}
			}

		}
	}
	if ramCounter > 0 {
		usedRAMAvg = usedRAMAvg / ramCounter
	}
	templateData.usedRAM = strconv.FormatInt(usedRAMAvg, 10) + "/" + strconv.FormatInt(totalRAM, 10) + " MB"
	if diskCounter > 0 {
		usedDiskAvg = usedDiskAvg / diskCounter
	}
	templateData.usedDisk = strconv.FormatInt(usedDiskAvg, 10) + "/" + strconv.FormatInt(totalDisk, 10) + " GB"
	if cpuCounter > 0 {
		cpuAvg = cpuAvg / cpuCounter
	}
	templateData.cpu = strconv.FormatInt(cpuAvg, 10) + "%"
	templateData.reboots = strconv.FormatInt(rebootCounter, 10)
	if templateData.reachable == false {
		unreachableFlag = true
	}
	templateData.temperature = strconv.FormatInt(maxTemp, 10)
}
