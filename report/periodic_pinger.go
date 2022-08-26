package report

import (
	"fmt"
	pb "testbed-monitor/pinger"
	"testbed-monitor/probers"
	"time"
)

type Monitor struct {
	lastReachable time.Time
}

type NotificationTemplate struct {
	TestbedIP string
	Timestamp string
}

var dailyFlag = false

func NewMonitor() *Monitor {
	monitor := new(Monitor)
	monitor.lastReachable = time.Time{}
	return monitor
}

func (monitor *Monitor) Start(IP string, period int) {
	fmt.Printf("Starting periodic pinger...\n")
	replyCh := make(chan *pb.PingReply)
	var p *pb.PingReply
	var emailDataN NotificationTemplate
	emailDataN.TestbedIP = IP
	ticker := time.NewTicker(time.Duration(period) * time.Minute)
	dailyCheck := time.NewTicker(60 * time.Minute)
	go func() {
		for range ticker.C {
			fmt.Printf("\nPinging testbed @ %s...\n", IP)
			icmp := probers.NewICMPProbe(IP, replyCh)
			icmp.Start()
			p = <-replyCh
			if p.Reachable == true {
				monitor.lastReachable = time.Now()
				fmt.Printf("Testbed %s was reached at %s\n", IP, monitor.lastReachable.String())
			} else {
				fmt.Printf("Testbed %s could not be reached.\n", IP)
				if dailyFlag == false {
					subjectN := emailDataN.TestbedIP + " could not be reached at " + emailData.Timestamp
					MailNotification(subjectN, emailDataN)
					dailyFlag = true
				}
			}
		}
	}()
	// Ensure only one notification email is sent per day
	go func() {
		for range dailyCheck.C {
			if dailyFlag == true {
				if time.Now().Hour() == 0 {
					dailyFlag = false
				}
			}
		}
	}()
}
