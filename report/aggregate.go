package report

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"testbed-monitor/traps"
	"time"
)

type Aggregate struct {
	statusChan         chan *StatusReport
	aggregatePeriod    int
	aggregateHour      int
	apiIP              string
	apiPort            string
	traps              map[string]trigger
	notificationFields map[string]string
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
	field            string
	operator         string
	trigger          int64
	period           int
	lastNotification time.Time
	flag             bool
}

func NewAggregate(statusChan chan *StatusReport, aggregatePeriod int, aggregateHour int, apiIP string, apiPort string) *Aggregate {
	aggregate := new(Aggregate)
	aggregate.statusChan = statusChan
	aggregate.aggregatePeriod = aggregatePeriod
	aggregate.aggregateHour = aggregateHour
	aggregate.apiIP = apiIP
	aggregate.apiPort = apiPort
	aggregate.traps = make(map[string]trigger)
	aggregate.notificationFields = make(map[string]string)
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
	for {
		select {
		case <-dailyTicker.C:
			if time.Now().Hour() == aggregate.aggregateHour {
				emailData.report = nil
				unreachableFlag = false
				if time.Now().Weekday().String() == "Sunday" { // Weekly report on Sundays
					for i = 0; i < len(*iPs); i++ {
						aggregator(reportAggregate, (*iPs)[i], &aggregatedReport, "Week")
						emailData.report = append(emailData.report, aggregatedReport)
					}
					emailData.reportType = "Weekly"
				} else { // Daily report
					for i = 0; i < len(*iPs); i++ {
						aggregator(reportAggregate, (*iPs)[i], &aggregatedReport, "Day")
						emailData.report = append(emailData.report, aggregatedReport)
					}
					emailData.reportType = "24 Hour"
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
				aggregate.towerAlertInApp(msg.Tower)
			}
			notificationFlag = false
			for trapField := range aggregate.traps {
				for _, reportField := range fields {
					if aggregate.traps[trapField].field == reportField.Name {
						fieldValue := getReportValue(msg, reportField.Name)
						fieldName := aggregate.notificationFields[reportField.Name]
						if aggregate.traps[trapField].flag == false {
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
								var notificationData NotificationTemplate
								notificationData = setNotification(msg.Tower, fieldName, strconv.FormatInt(fieldValue, 10))
								subject = msg.Tower + " " + fieldName + " Notification"
								MailNotification(subject, notificationData)
								if entry, ok := aggregate.traps[trapField]; ok {
									entry.flag = true
									entry.lastNotification = time.Now()
								}
							}
						}
					}
				}
			}
		case <-periodTicker.C:
			for trap := range aggregate.traps {
				if time.Now().After((aggregate.traps[trap].lastNotification).Add(time.Duration(aggregate.traps[trap].period) * time.Hour)) {
					if entry, ok := aggregate.traps[trap]; ok {
						entry.flag = false
					}
				}
			}
		}
	}
}

func (aggregate *Aggregate) SetTriggers(traps []traps.Config) {
	fields := reflect.VisibleFields(reflect.TypeOf(struct{ reportData }{}))
	log.Println("Thresholds:")
	for _, field := range fields {
		for _, t := range traps {
			if field.Name == t.Field {
				setTrigger := trigger{}
				setTrigger.field = t.Field
				setTrigger.operator = t.Operator
				setTrigger.trigger, _ = strconv.ParseInt(t.Trigger, 10, 64)
				setTrigger.period, _ = strconv.Atoi(t.Period)
				setTrigger.flag = false
				aggregate.traps[field.Name] = setTrigger
				log.Printf("%s %s %s\n", t.Field, t.Operator, t.Trigger)
			}
		}
	}
	aggregate.notificationFields["reboots"] = "Daily reboot count"
	aggregate.notificationFields["usedRAM"] = "MB RAM used"
	aggregate.notificationFields["usedDisk"] = "GB Disk used"
	aggregate.notificationFields["cpu"] = "CPU %"
	aggregate.notificationFields["temperature"] = "Temperature"
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

func aggregator(reports map[time.Time]*StatusReport, iP string, templateData *reportData, reportType string) {
	var usedRAMAvg, ramCounter, usedDiskAvg, diskCounter, cpuAvg, cpuCounter, rebootCounter int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk, maxTemp int64 = 0, 0, 0
	var interval time.Duration
	compareTime := time.Time{}
	templateData.tower = iP
	for key, element := range reports {
		templateData.tower = iP
		if reportType == "Week" {
			interval = -7 * 24 * time.Hour
		} else {
			interval = -24 * time.Hour
		}
		if key.After(time.Now().Add(interval)) && element.Tower == iP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.reachable = true
				if element.Reachable == true {
					templateData.arduinoReached = element.ArduinoReached
					templateData.towerReached = element.TowerReached
					templateData.bootTime = element.BootTime
				}
				if totalRAM == 0 {
					totalRAM = element.TotalRAM
				}
				if totalDisk == 0 {
					totalDisk = element.TotalDisk
				}
				if !element.Reachable {
					templateData.reachable = false
				}
				compareTime = element.Timestamp.AsTime()
			}
			rebootCounter = rebootCounter + element.Reboots
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
