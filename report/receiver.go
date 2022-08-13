package report

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
	"strings"
	"testbed-monitor/measure"
	pb "testbed-monitor/pinger"
	"testbed-monitor/probers"
	"time"
)

type Receiver struct {
	target     string
	connection *net.UDPConn
	measureCh  chan *measure.Measure
	statusCh   chan *StatusReport
}

func NewReportReceiver(measureCh chan *measure.Measure, statusCh chan *StatusReport) (*Receiver, error) {
	receiver := new(Receiver)
	receiver.measureCh = measureCh
	receiver.statusCh = statusCh
	s, err := net.ResolveUDPAddr("udp4", ":"+os.Getenv("RECEIVE_PORT"))
	if err != nil {
		return nil, err
	}
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		return nil, err
	}
	receiver.connection = connection
	return receiver, nil
}

func (receiver *Receiver) Start(towers *[]string) {
	ticker := time.NewTicker(60 * time.Minute)
	receivedReports := map[string]time.Time{}
	// Go func to receive reports
	go func() {
		receiver.receive(receivedReports, towers)
	}()
	// Go func to ping a host when no report in 60 minutes
	go func() {
		time.Sleep(60 * time.Minute)
		replyCh := make(chan *pb.PingReply)
		var p *pb.PingReply
		for range ticker.C {
			for key, element := range receivedReports {
				if time.Now().After(element.Add(60 * time.Minute)) {
					fmt.Printf("Haven't received a report from %s in 60 minutes. Attempting to ping...\n", key)
					icmp := probers.NewICMPProbe(key, replyCh)
					icmp.Start()
					p = <-replyCh
					if p.Reachable == true {
						element = time.Now()
						fmt.Printf("\nTower %s was reached at %s\n", key, element)
					} else {
						fmt.Printf("\nTower %s is unreachable!\n", key)
					}
				}
			}
		}
	}()
}

func (receiver *Receiver) receive(receivedReports map[string]time.Time, towers *[]string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal error while receiving: %s, will sleep 5 seconds before attempting to proceed\n", err)
			// since it will try again do not terminate
			time.Sleep(5 * time.Second)
		} else {
			err := receiver.connection.Close()
			if err != nil {
				fmt.Printf("Fatal error: %s, while closing the udp connection\n", err)
			}
		}
	}()

	buffer := make([]byte, 65000)
	for {
		n, addr, err := receiver.connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Fatal error: %s, while reading from udp connection\n", err)
			continue
		}

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		m := measure.Measure{}
		err = proto.Unmarshal(buffer[0:n], &m)
		if err != nil {
			fmt.Printf("Fatal error %s while unmarshalling message from %s!\n", err, addr)
			continue
		}

		fmt.Printf("Received report from %s.\nReport: %s\n", addr, m.String())
		receivedReports[addr.IP.String()] = time.Now()

		// Check to see if this tower is new (aggregate will look at towers[])
		towerKnown := false
		for i := 0; i < len(*towers); i++ {
			if (*towers)[i] == addr.IP.String() {
				towerKnown = true
			}
		}
		if towerKnown == false {
			*towers = append(*towers, addr.IP.String())
		}

		if m.Strings == nil {
			// some entries don't have a string map set
			m.Strings = make(map[string]string)
		}

		if m.Strings["host_id"] != "Hello" {
			s := &StatusReport{}
			GetStatusFromMeasure(addr.IP.String(), &m, s)
			receiver.statusCh <- s
		}
	}
}

// GetStatusFromMeasure reads a Measure Report and prepares a Status Report for the app to use later
func GetStatusFromMeasure(ip string, m *measure.Measure, s *StatusReport) {
	s.TowerIP = ip
	s.LastArduinoReachableTimestamp = time.Now().Add(time.Duration(m.Integers["LastArduinoReachableTimestamp"])).Format(time.RFC822)
	s.LastTowerReachableTimestamp = time.Now().Format(time.RFC822)
	s.BootTimestamp = m.Strings["bootTime"]
	s.RebootsCurrentDay = m.Integers["reboots"]
	s.RAMUsed = m.Integers["vm_used"]
	s.RAMTotal = m.Integers["vm_total"]
	s.DiskUsed = m.Integers["DISK_USED"]
	s.DiskTotal = m.Integers["DISK_TOTAL"]
	s.CPUAvg = m.Integers["CPU_AVG"]
	s.Timestamp = m.Timestamp
}
