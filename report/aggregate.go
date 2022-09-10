package report

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Aggregate struct {
	statusChan      chan *StatusReport
	aggregatePeriod int
	aggregateHour   int
	criticalTemp    int
	apiIP           string
	apiPort         string
}

type TemplateData struct {
	TowerIP                       string `json:"tower"`
	LastArduinoReachableTimestamp string `json:"arduinoReached"`
	LastTowerReachableTimestamp   string `json:"towerReached"`
	BootTimestamp                 string `json:"bootTime"`
	RebootsCurrentDay             string `json:"reboots"`
	RAMUsedAvgMB                  string `json:"usedRAM"`
	DiskUsedAvgGB                 string `json:"usedDisk"`
	CPUAvg                        string `json:"cpu"`
	Reachable                     bool   `json:"reachable"`
	Temperature                   string `json:"temperature"`
}

type reportTemplate struct {
	ReportType string         `json:"type"`
	Timestamp  string         `json:"timestamp"`
	Template   []TemplateData `json:"report"`
}

func NewAggregate(statusChan chan *StatusReport, aggregatePeriod int, aggregateHour int, criticalTemp int, apiIP string, apiPort string) *Aggregate {
	aggregate := new(Aggregate)
	aggregate.statusChan = statusChan
	aggregate.aggregatePeriod = aggregatePeriod
	aggregate.aggregateHour = aggregateHour
	aggregate.criticalTemp = criticalTemp
	aggregate.apiIP = apiIP
	aggregate.apiPort = apiPort
	return aggregate
}

var aggregatedReport TemplateData
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
				emailData.Template = nil
				unreachableFlag = false
				if time.Now().Weekday().String() == "Sunday" { // Weekly report on Sundays
					for i = 0; i < len(*iPs); i++ {
						weeklyAggregator(reportAggregate, (*iPs)[i], &aggregatedReport)
						emailData.Template = append(emailData.Template, aggregatedReport)
					}
					emailData.ReportType = "Weekly"
				} else { // Daily report
					for i = 0; i < len(*iPs); i++ {
						dailyAggregator(reportAggregate, (*iPs)[i], &aggregatedReport)
						emailData.Template = append(emailData.Template, aggregatedReport)
					}
					emailData.ReportType = "Daily"
				}
				emailData.Timestamp = time.Now().Format("Jan 02 2006 15:04:05")
				subject = emailData.ReportType + " Testbed Status Report " + emailData.Timestamp
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
		}
	}
}

func (aggregate *Aggregate) postStatusToApp(emailData reportTemplate) {
	apiAddress := aggregate.apiIP + ":" + aggregate.apiPort
	jsonReport, errJ := json.Marshal(emailData.Template)
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

func dailyAggregator(reports map[time.Time]*StatusReport, iP string, templateData *TemplateData) {
	var usedRAMAvg, ramCounter, usedDiskAvg, diskCounter, cpuAvg, cpuCounter, rebootCounter int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk, maxTemp int64 = 0, 0, 0
	compareTime := time.Time{}
	templateData.TowerIP = iP
	for key, element := range reports {
		if key.Day() == time.Now().Day() && element.TowerIP == iP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.Reachable = true
				if element.Reachable == true {
					templateData.LastArduinoReachableTimestamp = element.LastArduinoReachableTimestamp
					templateData.LastTowerReachableTimestamp = element.LastTowerReachableTimestamp
					templateData.BootTimestamp = element.BootTimestamp
				}
				if totalRAM == 0 {
					totalRAM = element.RAMTotal
				}
				if totalDisk == 0 {
					totalDisk = element.DiskTotal
				}
				if !element.Reachable {
					templateData.Reachable = false
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
	templateData.RAMUsedAvgMB = strconv.FormatInt(usedRAMAvg, 10) + "/" + strconv.FormatInt(totalRAM, 10) + " MB"
	if diskCounter > 0 {
		usedDiskAvg = usedDiskAvg / diskCounter
	}
	templateData.DiskUsedAvgGB = strconv.FormatInt(usedDiskAvg, 10) + "/" + strconv.FormatInt(totalDisk, 10) + " GB"
	if cpuCounter > 0 {
		cpuAvg = cpuAvg / cpuCounter
	}
	templateData.CPUAvg = strconv.FormatInt(cpuAvg, 10) + "%"
	templateData.RebootsCurrentDay = strconv.FormatInt(rebootCounter, 10)
	if templateData.Reachable == false {
		unreachableFlag = true
	}

}

func weeklyAggregator(reports map[time.Time]*StatusReport, iP string, templateData *TemplateData) {
	var usedRAMAvg, ramCounter, usedDiskAvg, diskCounter, cpuAvg, cpuCounter, rebootCounter int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk, maxTemp int64 = 0, 0, 0
	compareTime := time.Time{}
	templateData.TowerIP = iP
	for key, element := range reports {
		templateData.TowerIP = iP
		if key.After(time.Now().Add(-7*24*time.Hour)) && element.TowerIP == iP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.Reachable = true
				if element.Reachable == true {
					templateData.LastArduinoReachableTimestamp = element.LastArduinoReachableTimestamp
					templateData.LastTowerReachableTimestamp = element.LastTowerReachableTimestamp
					templateData.BootTimestamp = element.BootTimestamp
				}
				if totalRAM == 0 {
					totalRAM = element.RAMTotal
				}
				if totalDisk == 0 {
					totalDisk = element.DiskTotal
				}
				if !element.Reachable {
					templateData.Reachable = false
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
	templateData.RAMUsedAvgMB = strconv.FormatInt(usedRAMAvg, 10) + "/" + strconv.FormatInt(totalRAM, 10) + " MB"
	if diskCounter > 0 {
		usedDiskAvg = usedDiskAvg / diskCounter
	}
	templateData.DiskUsedAvgGB = strconv.FormatInt(usedDiskAvg, 10) + "/" + strconv.FormatInt(totalDisk, 10) + " GB"
	if cpuCounter > 0 {
		cpuAvg = cpuAvg / cpuCounter
	}
	templateData.CPUAvg = strconv.FormatInt(cpuAvg, 10) + "%"
	templateData.RebootsCurrentDay = strconv.FormatInt(rebootCounter, 10)
	if templateData.Reachable == false {
		unreachableFlag = true
	}
	templateData.Temperature = strconv.FormatInt(maxTemp, 10)
}
