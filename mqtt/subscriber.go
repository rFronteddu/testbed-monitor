package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
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
		if ContainsTemperature(msg.Payload()) {
			tower, temperature, timestamp := Parse(msg.Payload())
			r := &report.StatusReport{}
			r.Tower = tower
			r.Temperature = int64(temperature)
			r.Timestamp = timestamppb.New(timestamp)
			r.Reachable = true
			log.Println("Temperature on tower ", r.Tower, " is ", r.Temperature)
			mqttCh <- r
		}
	}
	options.SetDefaultPublishHandler(messagePubHandler)
	connectHandler := func(client mqtt.Client) {
		log.Println("Subscriber connected")
	}
	options.OnConnect = connectHandler
	connectionLostHandler := func(client mqtt.Client, err error) {
		log.Printf("Subscriber connection Lost: %s\n", err.Error())
	}
	options.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

}
