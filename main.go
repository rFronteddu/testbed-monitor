package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"strconv"
	"testbed-monitor/measure"
	"testbed-monitor/report"
)

type Configuration struct {
	MonitorHosts         bool   `yaml:"MONITOR_HOSTS"`
	ReceivePort          string `yaml:"RECEIVE_PORT"`
	AggregatePeriod      string `yaml:"AGGREGATE_PERIOD"`
	AggregateHour        string `yaml:"AGGREGATE_HOUR"`
	ExpectedReportPeriod string `yaml:"EXPECTED_REPORT_PERIOD"`
	CriticalTemp         string `yaml:"CRITICAL_TEMP"`
	MonitorTestbed       bool   `yaml:"MONITOR_TESTBED"`
	TestbedIP            string `yaml:"TESTBED_IP"`
	PingPeriod           string `yaml:"PING_PERIOD"`
}

func loadConfiguration(path string) *Configuration {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not open %s error: %s\n", path, err)
		conf := &Configuration{true, "8758", "60", "23", "30", "200", false, "", "30"}
		fmt.Printf("Host Monitor will use default configuration: %v\n", conf)
		return conf
	}
	if yfile == nil {
		panic("There was no error but YFile was null")
	}
	conf := Configuration{}
	err2 := yaml.Unmarshal(yfile, &conf)
	if err2 != nil {
		fmt.Printf("Configuration file could not be parsed, error: %s\n", err2)
		panic(err2)
	}
	fmt.Printf("Found configuration: %v\n", conf)
	return &conf
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := loadConfiguration("configuration.yaml")
	if conf.MonitorHosts {
		measureCh := make(chan *measure.Measure)
		statusCh := make(chan *report.StatusReport)
		var towers []string
		expectedReportPeriod := 30
		if conf.ExpectedReportPeriod != "" {
			expectedReportPeriod, err = strconv.Atoi(conf.ExpectedReportPeriod)
			if err != nil {
				fmt.Printf("Error converting %s to integer, expected report period set to default (60)", conf.ExpectedReportPeriod)
			}
			fmt.Printf("Host will be pinged if no report in %v minutes\n", conf.ExpectedReportPeriod)
		}
		receiver, err := report.NewReportReceiver(measureCh, statusCh, conf.ReceivePort, expectedReportPeriod)
		if err != nil {
			fmt.Printf("Fatal error %s while creating the report.Receiver, aborting\n", err)
		}
		fmt.Printf("Starting report receiver on port %s...\n", conf.ReceivePort)
		receiver.Start(&towers)

		aggregatePeriod := 60
		if conf.AggregatePeriod != "" {
			aggregatePeriod, err = strconv.Atoi(conf.AggregatePeriod)
			if err != nil {
				fmt.Printf("Error converting %s to integer, aggregate period set to default (60)", conf.AggregatePeriod)
			}
			fmt.Printf("Aggregate will check to send email every %v minutes\n", conf.AggregatePeriod)
		}
		aggregateHour := 23
		if conf.AggregateHour != "" {
			aggregateHour, err = strconv.Atoi(conf.AggregateHour)
			if err != nil {
				fmt.Printf("Error converting %s to integer, aggregate hour set to default (23)", conf.AggregateHour)
			}
			fmt.Printf("Daily report will be emailed at hour %v\n", conf.AggregateHour)
		}
		criticalTemp := 200
		if conf.CriticalTemp != "" {
			criticalTemp, err = strconv.Atoi(conf.CriticalTemp)
			if err != nil {
				fmt.Printf("Error converting %s to integer, critical temperature set to default (200)", conf.CriticalTemp)
			}
			fmt.Printf("Program will notify user if temperature is above %v\n", conf.CriticalTemp)
		}

		aggregate := report.NewAggregate(statusCh, aggregatePeriod, aggregateHour, criticalTemp)
		aggregate.Start(&towers)
	}

	if conf.MonitorTestbed {
		testbedIP := conf.TestbedIP
		pingPeriod := 30
		if conf.PingPeriod != "" {
			pingPeriod, err = strconv.Atoi(conf.PingPeriod)
			if err != nil {
				fmt.Printf("Error converting %s to integer, ping period set to default (30)", conf.PingPeriod)
			}
			fmt.Printf("Program will ping testbed %s every %v minutes\n", conf.TestbedIP, conf.PingPeriod)
		}
		monitor := report.NewMonitor()
		monitor.Start(testbedIP, pingPeriod)
	}

	//generatedConf, resolver := graph.NewResolver()
	//proxy, _ := db.NewProxy(resolver, measureCh)
	//proxy.Start()
	//go gqlServer(":8081", generatedConf)

	quitCh := make(chan int)
	<-quitCh
}

//func gqlServer(serveAddr string, conf generated.Config) {
//	router := chi.NewRouter()
//	router.Use(cors.New(cors.Options{
//		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000", "*"},
//		AllowCredentials: true,
//		Debug:            false,
//	}).Handler)
//	srv := handler.NewDefaultServer(generated.NewExecutableSchema(conf))
//	srv.AddTransport(&transport.Websocket{
//		Upgrader: websocket.Upgrader{
//			CheckOrigin: func(r *http.Request) bool {
//				// Check against your desired domains here
//				fmt.Printf("Host: %s\n", r.Host)
//				//return r.Host == "http://localhost:3000"
//				return true
//			},
//			ReadBufferSize:  1024,
//			WriteBufferSize: 1024,
//		},
//	})
//
//	router.Handle("/", playground.Handler("Testbed Manager Playground", "/query"))
//	router.Handle("/query", srv)
//	fmt.Printf("Connect to https://%s/ for playground\n", serveAddr)
//	err := http.ListenAndServe(serveAddr, router)
//	if err != nil {
//		panic(err)
//	}
//}
