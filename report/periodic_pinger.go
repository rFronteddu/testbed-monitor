package report

import (
	"fmt"
	"net/http"
	pb "testbed-monitor/pinger"
	"testbed-monitor/probers"
	"time"
)

type Monitor struct {
	lastReachable time.Time
	myIP          string
	apiPort       string
}

type NotificationTemplate struct {
	TestbedIP string
	Timestamp string
}

func NewMonitor(apiPort string) *Monitor {
	monitor := new(Monitor)
	monitor.lastReachable = time.Time{}
	monitor.myIP = string(GetOutboundIP())
	monitor.apiPort = apiPort
	return monitor
}

func (monitor *Monitor) Start(ip string, period int) {
	fmt.Printf("Starting periodic pinger...\n")
	dailyFlag := false
	replyCh := make(chan *pb.PingReply)
	var p *pb.PingReply
	var emailDataN NotificationTemplate
	emailDataN.TestbedIP = ip
	ticker := time.NewTicker(time.Duration(period) * time.Minute)
	dailyCheck := time.NewTicker(60 * time.Minute)
	go func() {
		for range ticker.C {
			fmt.Printf("\nPinging testbed @ %s...\n", ip)
			icmp := probers.NewICMPProbe(ip, replyCh)
			icmp.Start()
			p = <-replyCh
			if p.Reachable == true {
				monitor.lastReachable = time.Now()
				fmt.Printf("Testbed %s was reached at %s\n", ip, monitor.lastReachable.String())
			} else {
				fmt.Printf("Testbed %s could not be reached.\n", ip)
				emailDataN.Timestamp = time.Now().Format("Jan 02 2006 15:04:05")
				if !dailyFlag {
					subjectN := emailDataN.TestbedIP + " could not be reached at " + emailDataN.Timestamp
					MailNotification(subjectN, emailDataN)
					monitor.towerAlertInApp(ip)
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

func (monitor *Monitor) towerAlertInApp(alertIP string) {
	apiAddress := monitor.myIP + ":" + monitor.apiPort
	_, err := http.Get("http://" + apiAddress + "/alert/" + alertIP)
	if err != nil {
		fmt.Println(err)
	}
}
