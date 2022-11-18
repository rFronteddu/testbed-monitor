package report

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testbed-monitor/traps"
	"time"
)

type Aggregate struct {
	statusChan      chan *StatusReport
	aggregatePeriod int
	aggregateHour   int
	apiAddress      string
	apiPort         string
	apiReport       string
	apiAlert        string
	traps           map[string]trigger
}

type reportData struct {
	Tower          string `json:"tower"`
	ArduinoReached string `json:"arduinoReached"`
	TowerReached   string `json:"towerReached"`
	BootTime       string `json:"bootTime"`
	Reboots        string `json:"reboots"`
	UsedRAM        string `json:"usedRAM"`
	UsedDisk       string `json:"usedDisk"`
	Cpu            string `json:"cpu"`
	Reachable      bool   `json:"reachable"`
	Temperature    string `json:"temperature"`
}

type reportTemplate struct {
	ReportType string       `json:"type"`
	Timestamp  string       `json:"timestamp"`
	Report     []reportData `json:"report"`
}

type trigger struct {
	field            string
	operator         string
	trigger          int64
	period           int
	lastNotification time.Time
	flag             map[string]bool
}

func NewAggregate(statusChan chan *StatusReport, aggregatePeriod int, aggregateHour int, apiAddress string, apiPort string, apiReport string, apiAlert string) *Aggregate {
	aggregate := new(Aggregate)
	aggregate.statusChan = statusChan
	aggregate.aggregatePeriod = aggregatePeriod
	aggregate.aggregateHour = aggregateHour
	aggregate.apiAddress = apiAddress
	aggregate.apiPort = apiPort
	aggregate.apiReport = apiReport
	aggregate.apiAlert = apiAlert
	aggregate.traps = make(map[string]trigger)
	return aggregate
}

var aggregatedReport reportData
var emailData reportTemplate
var subject string
var unreachableFlag bool
var notificationFlag bool

func (aggregate *Aggregate) Start(iPs *[]string) {
	dailyTicker := time.NewTicker(time.Duration(aggregate.aggregatePeriod) * time.Minute)
	periodTicker := time.NewTicker(time.Duration(aggregate.aggregatePeriod) * time.Hour)
	reportAggregate := make(map[time.Time]*StatusReport)
	fields := reflect.VisibleFields(reflect.TypeOf(struct{ reportData }{}))
	var i int
	// Go func to handle & aggregate reports
	go func() {
		for {
			select {
			case <-dailyTicker.C:
				if time.Now().Hour() == aggregate.aggregateHour {
					emailData.Report = nil
					emailData.ReportType = "24 Hour"
					unreachableFlag = false
					if time.Now().Weekday().String() == "Sunday" { // Weekly report on Sundays
						emailData.ReportType = "Weekly"
						for i = 0; i < len(*iPs); i++ {
							aggregator(reportAggregate, (*iPs)[i], &aggregatedReport, "Week")
							emailData.Report = append(emailData.Report, aggregatedReport)
						}
					} else { // Daily report
						for i = 0; i < len(*iPs); i++ {
							aggregator(reportAggregate, (*iPs)[i], &aggregatedReport, "Day")
							emailData.Report = append(emailData.Report, aggregatedReport)
						}
					}
					emailData.Timestamp = time.Now().Format("Jan 02 2006 15:04:05")
					subject = emailData.ReportType + " Testbed Status Report " + emailData.Timestamp
					if unreachableFlag {
						subject += ": 1 or more host monitors are down!"
					}
					MailReport(subject, emailData)
					aggregate.postStatusToApp(emailData)
				}
			case msg := <-aggregate.statusChan:
				reportAggregate[time.Now()] = msg
				if !msg.Reachable {
					aggregate.towerAlertInApp(msg.Tower)
				}
				notificationFlag = false
				for trapField := range aggregate.traps {
					for _, reportField := range fields {
						if strings.EqualFold(aggregate.traps[trapField].field, reportField.Name) {
							fieldValue := getReportValue(msg, aggregate.traps[trapField].field)
							if _, ok := aggregate.traps[trapField].flag[msg.Tower]; !ok {
								aggregate.traps[trapField].flag[msg.Tower] = false
							}
							if aggregate.traps[trapField].flag[msg.Tower] == false {
								switch aggregate.traps[trapField].operator {
								case ">":
									if checkGreater(fieldValue, aggregate.traps[trapField].trigger) {
										notificationFlag = true
									}
								case "<":
									if checkLess(fieldValue, aggregate.traps[trapField].trigger) {
										notificationFlag = true
									}
								case "=":
									if checkEqual(fieldValue, aggregate.traps[trapField].trigger) {
										notificationFlag = true
									}
								}
								if notificationFlag == true {
									if !aggregate.traps[trapField].flag[msg.Tower] {
										var notificationData NotificationTemplate
										notificationData = setNotification(msg.Tower, trapField, strconv.FormatInt(fieldValue, 10))
										subject = msg.Tower + " " + trapField + " " + aggregate.traps[trapField].operator + " " + strconv.FormatInt(aggregate.traps[trapField].trigger, 10)
										MailNotification(subject, notificationData)
										aggregate.traps[trapField].flag[msg.Tower] = true
									}
								}
							}
						}
					}
				}

			case <-periodTicker.C:
				for trap := range aggregate.traps {
					if time.Now().After((aggregate.traps[trap].lastNotification).Add(time.Duration(aggregate.traps[trap].period) * time.Hour)) {
						for i = 0; i < len(*iPs); i++ {
							aggregate.markFlag(trap, (*iPs)[i], false)
						}
					}
				}
			}
		}
	}()
}

func (aggregate *Aggregate) SetTriggers(traps []traps.Config) {
	reportDataFields := reflect.VisibleFields(reflect.TypeOf(struct{ reportData }{}))
	log.Println("Thresholds:")
	for _, field := range reportDataFields {
		for _, t := range traps {
			if strings.EqualFold(field.Name, t.Field) {
				setTrigger := trigger{}
				structField, _ := reflect.ValueOf(reportData{}).Type().FieldByName(field.Name)
				jsonName := structField.Tag.Get("json")
				setTrigger.field = jsonName
				setTrigger.operator = t.Operator
				setTrigger.trigger, _ = strconv.ParseInt(t.Trigger, 10, 64)
				setTrigger.period, _ = strconv.Atoi(t.Period)
				setTrigger.flag = make(map[string]bool)
				aggregate.traps[field.Name] = setTrigger
				log.Printf("%s %s %s\n", setTrigger.field, setTrigger.operator, setTrigger.trigger)
			}
		}
	}
}

func (aggregate *Aggregate) markFlag(fieldName string, ip string, value bool) {
	for _, t := range aggregate.traps {
		reset := trigger{}
		reset.field = t.field
		reset.operator = t.operator
		reset.trigger = t.trigger
		reset.period = t.period
		reset.flag[ip] = value
		reset.lastNotification = time.Now()
		aggregate.traps[fieldName] = reset
	}
}

func (aggregate *Aggregate) postStatusToApp(emailData reportTemplate) {
	apiAddress := aggregate.apiAddress + ":" + aggregate.apiPort
	jsonReport, errJ := json.Marshal(emailData.Report)
	if errJ != nil {
		log.Println("Error creating json object: ", errJ)
	}
	_, err := http.Post("http://"+apiAddress+"/api/towerreport", "application/json", bytes.NewBuffer(jsonReport))
	if err != nil {
		log.Println("Error posting to API", err)
	}
}

func (aggregate *Aggregate) towerAlertInApp(alertIP string) {
	apiAddress := aggregate.apiAddress + ":" + aggregate.apiPort
	_, err := http.Get("http://" + apiAddress + "/api/alert/" + alertIP)
	if err != nil {
		log.Println("Error posting to API", err)
	}
}

func aggregator(reports map[time.Time]*StatusReport, iP string, templateData *reportData, reportType string) {
	var usedRAMAvg, ramCounter, usedDiskAvg, diskCounter, cpuAvg, cpuCounter, reboots int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk, maxTemp int64 = 0, 0, 0
	var interval time.Duration
	compareTime := time.Time{}
	templateData.Tower = iP
	if reportType == "Week" {
		interval = -7 * 24 * time.Hour
	} else {
		interval = -24 * time.Hour
	}
	for key, element := range reports {
		if key.After(time.Now().Add(interval)) && element.Tower == iP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.Reachable = true
				if element.Reachable == true {
					templateData.ArduinoReached = element.ArduinoReached
					templateData.TowerReached = element.TowerReached
					templateData.BootTime = element.BootTime
				}
				if totalRAM == 0 {
					totalRAM = element.TotalRAM
				}
				if totalDisk == 0 {
					totalDisk = element.TotalDisk
				}
				if !element.Reachable {
					templateData.Reachable = false
				}
				if element.Reboots > 0 {
					reboots++
				}
				compareTime = element.Timestamp.AsTime()
			}
			usedRAMAvg = usedRAMAvg + element.UsedRAM
			ramCounter++
			usedDiskAvg = usedDiskAvg + element.UsedDisk
			diskCounter++
			cpuAvg = cpuAvg + element.Cpu
			cpuCounter++
			if element.Timestamp.AsTime().After(time.Now().Add(interval)) && element.Tower == iP {
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
	templateData.UsedRAM = strconv.FormatInt(usedRAMAvg, 10) + "/" + strconv.FormatInt(totalRAM, 10) + " MB"
	if diskCounter > 0 {
		usedDiskAvg = usedDiskAvg / diskCounter
	}
	templateData.UsedDisk = strconv.FormatInt(usedDiskAvg, 10) + "/" + strconv.FormatInt(totalDisk, 10) + " GB"
	if cpuCounter > 0 {
		cpuAvg = cpuAvg / cpuCounter
	}
	templateData.Cpu = strconv.FormatInt(cpuAvg, 10) + "%"
	if templateData.Reachable == false {
		unreachableFlag = true
	}
	templateData.Reboots = strconv.FormatInt(reboots, 10)
	templateData.Temperature = strconv.FormatInt(maxTemp, 10)
}

func checkGreater(reportVal int64, triggerVal int64) bool {
	if reportVal > triggerVal {
		return true
	} else {
		return false
	}
}

func checkLess(reportVal int64, triggerVal int64) bool {
	if reportVal < triggerVal {
		return true
	} else {
		return false
	}
}

func checkEqual(reportVal int64, triggerVal int64) bool {
	if reportVal == triggerVal {
		return true
	} else {
		return false
	}
}

func getReportValue(msg *StatusReport, field string) int64 {
	defer func() {
		if r := recover(); r != nil {
			// Function will panic if the field is not in the status report so recover
		}
	}()
	m, _ := json.Marshal(msg)
	var x map[string]interface{}
	_ = json.Unmarshal(m, &x)
	return int64(x[field].(float64))
}

func setNotification(tower string, field string, value string) (emailData NotificationTemplate) {
	emailData.TowerIP = tower
	emailData.Field = field
	emailData.Value = value
	emailData.Timestamp = time.Now().Format("Jan 02 2006 15:04:05")
	return emailData
}
