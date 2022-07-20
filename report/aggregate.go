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
	dailyTicker := time.NewTicker(60 * time.Minute)
	reportAggregate := make(map[time.Time]*StatusReport)
	for {
		select {
		case <-dailyTicker.C:
			if time.Now().Hour() == 23 { // 11 pm daily report
				DailyAggregator(reportAggregate, &templateData)
				Mail("Daily Testbed Status Report", templateData)
			}
			//if time.Now().Day() == 7 { // Sunday night weekly report
			//	if time.Now().Hour() == 18 {
			//		WeeklyAggregator(reportAggregate, &templateData)
			//		Mail("Weekly Testbed Status Report", templateData)
			//	}
			//}
		case msg := <-aggregate.statusChan:
			reportAggregate[time.Now()] = msg
		}
	}
}

func DailyAggregator(reports map[time.Time]*StatusReport, templateData *TemplateData) {
	var usedRAMAvg, RAMCounter, usedDiskAvg, diskCounter, CPUAvg, CPUCounter int64 = 0, 0, 0, 0, 0, 0
	var totalRAM, totalDisk int64 = 0, 0
	compareTime := time.Time{}
	for key, element := range reports {
		if key.Day() == time.Now().Day() {
			if element.Timestamp.AsTime().After(compareTime) {
				templateData.TowerIP = element.TowerIP
				templateData.LastArduinoReachableTimestamp = element.LastArduinoReachableTimestamp
				templateData.LastTowerReachableTimestamp = element.LastTowerReachableTimestamp
				templateData.BootTimestamp = element.BootTimestamp
				templateData.RebootsCurrentDay = strconv.FormatInt(element.RebootsCurrentDay, 10)
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
	templateData.Timestamp = time.Now().Format("Jan 02 2006")
}
