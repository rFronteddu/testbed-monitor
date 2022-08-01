package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"testbed-monitor/measure"
	"testbed-monitor/report"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	measureCh := make(chan *measure.Measure)
	statusCh := make(chan *report.StatusReport)
	var towers []string

	receiver, err := report.NewReportReceiver(measureCh, statusCh)
	if err != nil {
		fmt.Printf("Fatal error %s while creating the report.Receiver, aborting\n", err)
	}
	fmt.Println("Starting report receiver...")
	receiver.Start(towers)

	aggregate := report.NewAggregate(statusCh)
	aggregate.Start(towers)

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
