package mqtt

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
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

func Parse(data []byte) (tower string, temperature int, timestamp time.Time) {
	var message Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Println(err)
	}

	tower = message.Tower
	timeToken := strings.Split(message.Time, ":")
	utime, errt := strconv.Atoi(timeToken[1])
	if errt != nil {
		log.Println(errt)
	}
	timestamp = time.Unix(int64(utime), 0)

	var temperatureF float64
	for i := range message.Reading {
		if message.Reading[i].Name == "temperature" {
			temperatureF = message.Reading[i].Value
		}
	}
	temperature = int(temperatureF)

	return tower, temperature, timestamp
}
