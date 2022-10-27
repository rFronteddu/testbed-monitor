package mqtt

import (
	"encoding/json"
	"log"
	"time"
)

type Message struct {
	Location   Location  `json:"location"`
	Reading    []Reading `json:"reading"`
	SensorType string    `json:"sensorType"`
	Time       string    `json:"time"`
	Tower      string    `json:"tower"`
}
type Reading struct {
	Measurement string  `json:"measurement"`
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
}
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func ContainsTemperature(data []byte) (containsTemperature bool) {
	var message Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Println(err)
	}
	containsTemperature = false
	for i := range message.Reading {
		if message.Reading[i].Name == "temperature" {
			containsTemperature = true
		}
	}
	return containsTemperature
}

func Parse(data []byte, m map[string]string) (tower string, temperature int, timestamp time.Time) {
	var message Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Println(err)
	}

	towerName := message.Tower
	tower = m[towerName]

	timestamp = time.Now()

	var temperatureF float64
	for i := range message.Reading {
		if message.Reading[i].Name == "temperature" {
			temperatureF = message.Reading[i].Value
		}
	}
	temperature = int(temperatureF)

	return tower, temperature, timestamp
}
