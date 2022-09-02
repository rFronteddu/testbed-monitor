package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"testbed-monitor/report"
)

type Subscriber struct {
	broker string
	Client mqtt.Client
	mqttCh chan *report.StatusReport
}

func NewSubscriber(broker string, topic string, mqttCh chan *report.StatusReport) {
	var s Subscriber
	s.broker = broker
	s.mqttCh = mqttCh
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID("Host-monitor")
	messagePubHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Subscriber received a message on topic %s\n", msg.Topic())
		tower, temperature, timestamp := Parse(msg.Payload())
		fmt.Printf("Temperature is %v at %s\n", temperature, timestamp)
		r := &report.StatusReport{}
		r.TowerIP = tower
		r.Temperature = int64(temperature)
		r.Timestamp = timestamppb.New(timestamp)
		mqttCh <- r
	}
	options.SetDefaultPublishHandler(messagePubHandler)
	connectHandler := func(client mqtt.Client) {
		fmt.Println("Subscriber connected")
	}
	options.OnConnect = connectHandler
	connectionLostHandler := func(client mqtt.Client, err error) {
		fmt.Printf("Subscriber connection Lost: %s\n", err.Error())
	}
	options.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}
