package mqtt

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

// Message The structs are defined from schema.json
type Message struct {
	Location   Location      `json:"location"`
	Reading    []Measurement `json:"reading"`
	SensorType string        `json:"sensorType"`
	Time       string        `json:"time"`
	Tower      string        `json:"tower"`
}
type Measurement struct {
	Measurement string `json:"measurement"`
	Name        string `json:"name"`
	Value       string `json:"value"`
}
type Location struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func ContainsTemperature(data []byte) (containsTemperature bool) {
	var message Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Printf("Error unmarshalling json: %s\n", err)
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
		log.Printf("Error unmarshalling json: %s\n", err)
	}

	tower = message.Tower

	timeToken := strings.Split(message.Time, ":")
	utime, errt := strconv.Atoi(timeToken[1])
	if errt != nil {
		log.Println("Could not parse time\n", errt)
	}
	timestamp = time.Unix(int64(utime), 0)

	var temperatureS string
	for i := range message.Reading {
		if message.Reading[i].Name == "temperature" {
			temperatureS = message.Reading[i].Value
		}
	}
	temperatureF, errF := strconv.ParseFloat(temperatureS, 64)
	if errF != nil {
		log.Println("Could not parse temperature\n", errF)
	}
	temperature = int(temperatureF)

	return tower, temperature, timestamp
}
