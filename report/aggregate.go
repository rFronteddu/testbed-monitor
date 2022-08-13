package report

import (
	"strconv"
	"time"
)

type Aggregate struct {
	statusChan      chan *StatusReport
	aggregatePeriod int
	aggregateHour   int
}

type TemplateData struct {
	Name                          string
	TowerIP                       string
	LastArduinoReachableTimestamp string
	LastTowerReachableTimestamp   string
	BootTimestamp                 string
	RebootsCurrentDay             string
	RAMUsedAvgMB                  string
	DiskUsedAvgGB                 string
	CPUAvg                        string
	Timestamp                     string
}

type emailTemplate struct {
	ReportType string
	Template   []TemplateData
}

func NewAggregate(statusChan chan *StatusReport, aggregatePeriod int, aggregateHour int) *Aggregate {
	aggregate := new(Aggregate)
	aggregate.statusChan = statusChan
	aggregate.aggregatePeriod = aggregatePeriod
	aggregate.aggregateHour = aggregateHour
	return aggregate
}

var aggregatedReport TemplateData
var emailData emailTemplate

func (aggregate *Aggregate) Start(IPs *[]string) {
	dailyTicker := time.NewTicker(time.Duration(aggregate.aggregatePeriod) * time.Minute)
	reportAggregate := make(map[time.Time]*StatusReport)

	var i int

	for {
		select {
		case <-dailyTicker.C:
			if time.Now().Hour() == aggregate.aggregateHour { // Daily report
				emailData.Template = nil
				if time.Now().Weekday().String() == "Sunday" { // Weekly report on Sundays
					for i = 0; i < len(*IPs); i++ {
						WeeklyAggregator(reportAggregate, (*IPs)[i], &aggregatedReport)
						emailData.Template = append(emailData.Template, aggregatedReport)
					}
					emailData.ReportType = "Weekly"
					Mail("Weekly Testbed Status Report for "+aggregatedReport.Timestamp, emailData)
				} else {
					for i = 0; i < len(*IPs); i++ {
						DailyAggregator(reportAggregate, (*IPs)[i], &aggregatedReport)
						emailData.Template = append(emailData.Template, aggregatedReport)
					}
					emailData.ReportType = "Daily"
					Mail("Daily Testbed Status Report for "+aggregatedReport.Timestamp, emailData)
				}
			}
		case msg := <-aggregate.statusChan:
			reportAggregate[time.Now()] = msg
		}
	}
}

func DailyAggregator(reports map[time.Time]*StatusReport, IP string, templateData *TemplateData) {
	var usedRAMAvg, RAMCounter, usedDiskAvg, diskCounter, CPUAvg, CPUCounter, rebootCounter int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk int64 = 0, 0
	compareTime := time.Time{}
	for key, element := range reports {
		if key.Day() == time.Now().Day() && element.TowerIP == IP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.TowerIP = element.TowerIP
				templateData.LastArduinoReachableTimestamp = element.LastArduinoReachableTimestamp
				templateData.LastTowerReachableTimestamp = element.LastTowerReachableTimestamp
				templateData.BootTimestamp = element.BootTimestamp
				rebootCounter = rebootCounter + element.RebootsCurrentDay
				if totalRAM == 0 {
					totalRAM = element.RAMTotal
				}
				if totalDisk == 0 {
					totalDisk = element.DiskTotal
				}
				compareTime = element.Timestamp.AsTime()
			}
			usedRAMAvg = usedRAMAvg + element.RAMUsed
			RAMCounter++
			usedDiskAvg = usedDiskAvg + element.DiskUsed
			diskCounter++
			CPUAvg = CPUAvg + element.CPUAvg
			CPUCounter++
		}
	}
	if RAMCounter > 0 {
		usedRAMAvg = usedRAMAvg / RAMCounter
	}
	templateData.RAMUsedAvgMB = strconv.FormatInt(usedRAMAvg, 10) + "/" + strconv.FormatInt(totalRAM, 10) + " MB"
	if diskCounter > 0 {
		usedDiskAvg = usedDiskAvg / diskCounter
	}
	templateData.DiskUsedAvgGB = strconv.FormatInt(usedDiskAvg, 10) + "/" + strconv.FormatInt(totalDisk, 10) + " GB"
	if CPUCounter > 0 {
		CPUAvg = CPUAvg / CPUCounter
	}
	templateData.CPUAvg = strconv.FormatInt(CPUAvg, 10) + "%"
	templateData.RebootsCurrentDay = strconv.FormatInt(rebootCounter, 10)
	templateData.Timestamp = time.Now().Format("Jan 02 2006")
}

func WeeklyAggregator(reports map[time.Time]*StatusReport, IP string, templateData *TemplateData) {
	var usedRAMAvg, RAMCounter, usedDiskAvg, diskCounter, CPUAvg, CPUCounter, rebootCounter int64 = 0, 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk int64 = 0, 0
	compareTime := time.Time{}
	for key, element := range reports {
		if key.After(time.Now().Add(-7*24*time.Hour)) && element.TowerIP == IP {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.TowerIP = element.TowerIP
				templateData.LastArduinoReachableTimestamp = element.LastArduinoReachableTimestamp
				templateData.LastTowerReachableTimestamp = element.LastTowerReachableTimestamp
				templateData.BootTimestamp = element.BootTimestamp
				rebootCounter = rebootCounter + element.RebootsCurrentDay
				if totalRAM == 0 {
					totalRAM = element.RAMTotal
				}
				if totalDisk == 0 {
					totalDisk = element.DiskTotal
				}
				compareTime = element.Timestamp.AsTime()
			}
			usedRAMAvg = usedRAMAvg + element.RAMUsed
			RAMCounter++
			usedDiskAvg = usedDiskAvg + element.DiskUsed
			diskCounter++
			CPUAvg = CPUAvg + element.CPUAvg
			CPUCounter++
		}
	}
	if RAMCounter > 0 {
		usedRAMAvg = usedRAMAvg / RAMCounter
	}
	templateData.RAMUsedAvgMB = strconv.FormatInt(usedRAMAvg, 10) + "/" + strconv.FormatInt(totalRAM, 10) + " MB"
	if diskCounter > 0 {
		usedDiskAvg = usedDiskAvg / diskCounter
	}
	templateData.DiskUsedAvgGB = strconv.FormatInt(usedDiskAvg, 10) + "/" + strconv.FormatInt(totalDisk, 10) + " GB"
	if CPUCounter > 0 {
		CPUAvg = CPUAvg / CPUCounter
	}
	templateData.CPUAvg = strconv.FormatInt(CPUAvg, 10) + "%"
	templateData.RebootsCurrentDay = strconv.FormatInt(rebootCounter, 10)
	templateData.Timestamp = time.Now().Format("Jan 02 2006")
}
