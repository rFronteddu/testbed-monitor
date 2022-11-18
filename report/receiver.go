package report

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"strings"
	"testbed-monitor/measure"
	pb "testbed-monitor/pinger"
	"testbed-monitor/probers"
	"time"
)

type Receiver struct {
	target               string
	connection           *net.UDPConn
	measureCh            chan *measure.Measure
	statusCh             chan *StatusReport
	gqlCh                chan *StatusReport
	expectedReportPeriod int
}

func NewReportReceiver(measureCh chan *measure.Measure, statusCh chan *StatusReport, gqlCh chan *StatusReport, receivePort string, expectedReportPeriod int) (*Receiver, error) {
	receiver := new(Receiver)
	receiver.measureCh = measureCh
	receiver.statusCh = statusCh
	receiver.gqlCh = gqlCh
	s, err := net.ResolveUDPAddr("udp4", ":"+receivePort)
	if err != nil {
		log.Panicf("Unable to resolve address %s\n%s\n", s, err)
		return nil, err
	}
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		log.Panicf("Unable to listen on %s\n%s\n", s, err)
		return nil, err
	}
	receiver.connection = connection
	receiver.expectedReportPeriod = expectedReportPeriod
	return receiver, nil
}

func (receiver *Receiver) Start(towers *[]string) {
	ticker := time.NewTicker(time.Duration(receiver.expectedReportPeriod) * time.Minute)
	replyCh := make(chan *pb.PingReply)
	var p *pb.PingReply
	receivedReports := map[string]time.Time{}

	// Go func to receive reports
	go func() {
		receiver.receive(receivedReports, towers)
	}()
	// Go func to ping a host when no report in expected time
	go func() {
		for range ticker.C {
			for key, element := range receivedReports {
				if time.Now().After(element.Add(time.Duration(receiver.expectedReportPeriod) * time.Minute)) {
					log.Printf("Haven't received a report from %s in %v minutes. Attempting to ping...\n", key, receiver.expectedReportPeriod)
					icmp := probers.NewICMPProbe(key, replyCh)
					icmp.Start()
					p = <-replyCh
					if p.Reachable == true {
						element = time.Now()
						log.Printf("\nTower %s was reached at %s\n", key, element)
					} else {
						log.Printf("\nTower %s is unreachable!\n", key)
						s := &StatusReport{}
						s.Tower = key
						s.Timestamp = &timestamp.Timestamp{Seconds: time.Now().Unix()}
						s.Reachable = false
						receiver.statusCh <- s
					}
				}
			}
		}
	}()
}

func (receiver *Receiver) receive(receivedReports map[string]time.Time, towers *[]string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Fatal error while receiving: %s, will sleep 5 seconds before attempting to proceed\n", err)
			// since it will try again do not terminate
			time.Sleep(5 * time.Second)
		} else {
			err := receiver.connection.Close()
			if err != nil {
				log.Fatalf("Fatal error: %s, while closing the udp connection\n", err)
			}
		}
	}()

	buffer := make([]byte, 65000)
	for {
		n, addr, err := receiver.connection.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Fatal error: %s, while reading from udp connection\n", err)
			continue
		}

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			log.Println("Exiting UDP server!")
			return
		}

		m := measure.Measure{}
		err = proto.Unmarshal(buffer[0:n], &m)
		if err != nil {
			log.Printf("Fatal error %s while unmarshalling message from %s!\n", err, addr)
			continue
		}

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
			receiver.gqlCh <- s
		}

		log.Printf("Received report from %s.\nReport: %s\n", addr.IP.String(), m.String())
	}
}

// GetStatusFromMeasure reads a Measure Report and prepares a Status Report for the app to use later
func GetStatusFromMeasure(ip string, m *measure.Measure, s *StatusReport) {
	s.Tower = ip
	s.ArduinoReached = time.Now().Add(time.Duration(m.Integers["arduinoReached"])).Format(time.RFC822)
	s.TowerReached = time.Now().Format(time.RFC822)
	s.BootTime = time.Now().Add(time.Duration(m.Integers["bootTime"] * -1)).Format(time.RFC822)
	s.Reboots = m.Integers["reboot_sensor"]
	s.UsedRAM = m.Integers["vm_used"]
	s.TotalRAM = m.Integers["vm_total"]
	s.UsedDisk = m.Integers["DISK_USED"]
	s.TotalDisk = m.Integers["DISK_TOTAL"]
	s.Cpu = m.Integers["CPU_AVG"]
	s.Timestamp = m.Timestamp
	s.Reachable = true
	s.Temperature = m.Integers["MQTT_Temp"]
}
