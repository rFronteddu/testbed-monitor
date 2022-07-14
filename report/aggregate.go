package report

import (
	"strconv"
	"time"
)

type Aggregate struct {
	statusChan chan *StatusReport
}

func NewAggregate(statusChan chan *StatusReport) *Aggregate {
	aggregate := new(Aggregate)
	aggregate.statusChan = statusChan
	return aggregate
}

var templateData TemplateData

func (aggregate *Aggregate) Start() {
	dailyTicker := time.NewTicker(1 * time.Hour)
	reportAggregate := make(map[time.Time]*StatusReport)
	for {
		select {
		case <-dailyTicker.C:
			if time.Now().Hour() == 23 { // every day at 11pm
				Aggregator(reportAggregate, &templateData)
				subject := "Status Report for " + templateData.TowerIP
				Mail(subject, templateData)
			}
		case msg := <-aggregate.statusChan:
			reportAggregate[time.Now()] = msg
		}
	}
}

func Aggregator(reports map[time.Time]*StatusReport, templateData *TemplateData) {
	var CPUAvg, CPUAvgCounter int64 = 0, 0
	for key, element := range reports {
		if key.Day() == time.Now().Day() {
			templateData.TowerIP = element.TowerIP
			templateData.LastArduinoReachableTimestamp = element.LastArduinoReachableTimestamp
			templateData.LastTowerReachableTimestamp = element.LastTowerReachableTimestamp
			templateData.BootTimestamp = element.BootTimestamp
			templateData.RebootsCurrentDay = strconv.FormatInt(element.RebootsCurrentDay, 10)
			templateData.LastRamReadMB = strconv.FormatInt(element.LastRamReadMB, 10)
			templateData.LastDiskReadGB = strconv.FormatInt(element.LastDiskReadGB, 10)
			CPUAvg = CPUAvg + element.LastCPUAvg
			CPUAvgCounter++
		}
	}
	CPUAvg = CPUAvg / CPUAvgCounter
	templateData.LastCPUAvg = strconv.FormatInt(CPUAvg, 10)
}
