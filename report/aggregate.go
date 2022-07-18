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
				DailyAggregator(reportAggregate, &templateData)
				subject := "11pm Daily Status Report for " + templateData.TowerIP
				Mail(subject, templateData)
			}
		case msg := <-aggregate.statusChan:
			reportAggregate[time.Now()] = msg
		}
	}
}

func DailyAggregator(reports map[time.Time]*StatusReport, templateData *TemplateData) {
	var CPUAvg, CPUAvgCounter int64 = 0, 0
	compareTime := time.Time{}
	for key, element := range reports {
		if key.Day() == time.Now().Day() {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.TowerIP = element.TowerIP
				templateData.LastArduinoReachableTimestamp = element.LastArduinoReachableTimestamp
				templateData.LastTowerReachableTimestamp = element.LastTowerReachableTimestamp
				templateData.BootTimestamp = element.BootTimestamp
				templateData.RebootsCurrentDay = strconv.FormatInt(element.RebootsCurrentDay, 10)
				templateData.LastRamReadMB = strconv.FormatInt(element.LastRamReadMB, 10)
				templateData.LastDiskReadGB = strconv.FormatInt(element.LastDiskReadGB, 10)
				compareTime = element.Timestamp.AsTime()
			}

			CPUAvg = CPUAvg + element.LastCPUAvg
			CPUAvgCounter++
		}
	}
	if CPUAvgCounter > 0 {
		CPUAvg = CPUAvg / CPUAvgCounter
	}
	templateData.LastCPUAvg = strconv.FormatInt(CPUAvg, 10)
	templateData.Timestamp = time.Now().Format(time.RFC822)
}
