package report

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"strconv"
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
		var pingResponse, filename string
		for range ticker.C {
			for key, element := range receivedReports {
				if time.Now().After(element.Add(60 * time.Minute)) {
					fmt.Printf("Haven't received a report from %s in 60 minutes. Attempting to ping...\n", key)
					filename = time.Now().Format("2006-01-02_1504") + "_Ping.txt"
					icmp := probers.NewICMPProbe(key, replyCh)
					icmp.Start()
					p = <-replyCh
					if p.Reachable == true {
						element = time.Now()
						pingResponse = "Tower " + key + " reached at " + element.String() + "\nLost percentage: " + strconv.Itoa(int(p.LostPercentage)) + " Avg rtt: " + strconv.Itoa(int(p.AvgRtt))
						fmt.Printf("\nTower %s was reached at %s\n", key, element)
						LogPingReport(filename, pingResponse)
					} else {
						fmt.Printf("\nTower %s is unreachable!\n", key)
						pingResponse = "Tower " + key + "is unreachable at " + element.String()
						LogPingReport(filename, pingResponse)
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

			var filename = time.Now().Format("2006-01-02_150405") + "_Report.txt"
			err2 := LogReport(filename, m.String())
			if err2 != nil {
				fmt.Printf("Error in LogReport: %s", err2)
				time.Sleep(60 * time.Second)
				log.Fatal(err2)
			}
		}
	}
}

// LogReport This function writes report data to a file
func LogReport(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %s", err)
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		fmt.Printf("Error writing data to file: %s", err)
		return err
	}
	fmt.Printf("Report logged in %s\n", filename)
	return file.Sync()
}

// LogPingReport This function writes report data to a file
func LogPingReport(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %s", err)
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		fmt.Printf("Error writing data to file: %s", err)
		return err
	}
	return file.Sync()
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
