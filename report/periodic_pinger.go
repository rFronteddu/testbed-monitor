package report

import (
	"log"
	pb "testbed-monitor/pinger"
	"testbed-monitor/probers"
	"time"
)

type Monitor struct {
	lastReachable time.Time
}

func NewMonitor() *Monitor {
	monitor := new(Monitor)
	monitor.lastReachable = time.Time{}
	return monitor
}

func (monitor *Monitor) Start(ip string, period int) {
	log.Printf("Starting periodic pinger...\n")
	dailyFlag := false
	replyCh := make(chan *pb.PingReply)
	var p *pb.PingReply
	var emailDataN NotificationTemplate
	emailDataN.TowerIP = ip
	ticker := time.NewTicker(time.Duration(period) * time.Minute)
	dailyCheck := time.NewTicker(60 * time.Minute)
	go func() {
		for range ticker.C {
			log.Printf("\nPinging testbed @ %s...\n", ip)
			icmp := probers.NewICMPProbe(ip, replyCh)
			icmp.Start()
			p = <-replyCh
			if p.Reachable == true {
				monitor.lastReachable = time.Now()
				log.Printf("Testbed %s was reached at %s\n", ip, monitor.lastReachable.String())
			} else {
				log.Printf("Testbed %s could not be reached.\n", ip)
				setNotification(ip, "Tower", "not reachable")
				if !dailyFlag {
					subjectN := emailDataN.TowerIP + " could not be reached at " + emailDataN.Timestamp
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
