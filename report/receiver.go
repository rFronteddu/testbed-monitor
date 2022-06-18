package report

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testbed-monitor/measure"
	"testbed-monitor/task"
	"time"
)

const (
	SENSOR_IP = "sensor_ip"
)

type Receiver struct {
	target     string
	connection *net.UDPConn
	measureCh  chan *measure.Measure
}

func NewReportReceiver(measureCh chan *measure.Measure) (*Receiver, error) {
	receiver := new(Receiver)
	receiver.measureCh = measureCh
	s, err := net.ResolveUDPAddr("udp4", ":8758")
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

func (receiver *Receiver) Start() {
	fmt.Printf("Starting Report Receiver...\n")
	ticker := time.NewTicker(60 * time.Minute)
	receivedReports := map[string]time.Time{
		"127.0.0.1":     time.Time{},
		"123.456.789.0": time.Time{},
	}
	go func() {
		println("Receiving reports") //////////////////////////////
		receiver.receive(receivedReports)
	}()

	go func() {
		for _ = range ticker.C {
			println("Checking timestamps") ///////////////////////
			for key, element := range receivedReports {
				if time.Now().After(element.Add(60 * time.Minute)) {
					fmt.Printf("Haven't received a report for %s in 30 minutes. Will attempt to ping.\n", key)
					pingTask := task.NewPingTask()
					pingTask.Start(key)
				}
			}
		}
	}()
}

func (receiver *Receiver) receive(receivedReports map[string]time.Time) {
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
		}

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		m := measure.Measure{}
		err = proto.Unmarshal(buffer[0:n], &m)
		if err != nil {
			fmt.Printf("Fatal error %s while unmarshalling messagr from %s!\n", err, addr)
			continue
		}

		fmt.Printf("Received report from %s.\nReport: %s\n", addr, m.String())
		receivedReports[addr.IP.String()] = time.Now()

		if m.Strings == nil {
			// some entries don't have a string map set
			m.Strings = make(map[string]string)
		}
		m.Strings[SENSOR_IP] = addr.IP.String()

		receiver.measureCh <- &m

		var filename = time.Now().Format("2006-01-02_1504")
		filename += "_Report.txt"
		err2 := LogReport(filename, m.String())
		if err2 != nil {
			fmt.Printf("Error in LogReport: %s", err2)
			time.Sleep(60 * time.Second)
			log.Fatal(err2)
		}
		fmt.Printf("Report logged in %s\n", filename)
	}
}

// LogReport This function writes report data to a file
func LogReport(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
		fmt.Printf("Error creating file: %s", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
		fmt.Printf("Error writing data to file: %s", err)
	}
	return file.Sync()
}
