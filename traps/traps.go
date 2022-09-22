package traps

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Description string `json:"description"`
	Operator    string `json:"operator"`
	Field       string `json:"field"`
	Trigger     string `json:"trigger"`
	Period      string `json:"period"`
}

func LoadConfiguration(file string) []Config {
	var config []Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
