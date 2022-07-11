package report

import (
	"fmt"
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
	dailyTicker := time.NewTicker(1 * time.Minute)
	statusReport := &StatusReport{} // move this inside case msg
	//reportAggregate := make(map[string]*StatusReport)
	for {
		select {
		case <-dailyTicker.C:
			subject := "Status Report for" + statusReport.TowerIP
			templateData.TowerIP = statusReport.TowerIP
			templateData.LastArduinoReachableTimestamp = statusReport.LastArduinoReachableTimestamp
			templateData.LastTowerReachableTimestamp = statusReport.LastTowerReachableTimestamp
			templateData.BootTimestamp = statusReport.BootTimestamp
			templateData.RebootsCurrentDay = strconv.FormatInt(statusReport.RebootsCurrentDay, 10)
			templateData.LastRamReadMB = strconv.FormatInt(statusReport.LastRamReadMB, 10)
			templateData.LastDiskReadGB = strconv.FormatInt(statusReport.LastDiskReadGB, 10)
			templateData.LastCPUAvg = strconv.FormatInt(statusReport.LastCPUAvg, 10)
			templateData.Timestamp = statusReport.Timestamp.AsTime().Format("January 2, 2006")
			Mail(subject, templateData)
		case msg := <-aggregate.statusChan:
			statusReport = msg
			fmt.Printf("Status report received in channel: %s\n", msg.String())
			//reportAggregate[statusReport.TowerIP] = &statusReport
		}
	}

}
