package report

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"strings"
	"testbed-monitor/measure"
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
	fmt.Printf("Starting Report Receiver\n")
	go func() {
		receiver.receive()
	}()
}

func (receiver *Receiver) receive() {
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

		fmt.Printf("Received report from %s. Report: %s\n", addr, m.String())
		if m.Strings == nil {
			// some entries don't have a string map set
			m.Strings = make(map[string]string)
		}
		m.Strings[SENSOR_IP] = addr.IP.String()

		receiver.measureCh <- &m
	}
}
