package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"net/http"
	"testbed-monitor/db"
	"testbed-monitor/graph"
	"testbed-monitor/graph/generated"
	"testbed-monitor/measure"
	"testbed-monitor/report"
)

func main() {
	//pingTask := task.NewPingTask()
	//pingTask.Start()
	generatedConf, resolver := graph.NewResolver()

	measureCh := make(chan *measure.Measure)
	receiver, err := report.NewReportReceiver(measureCh)
	if err != nil {
		fmt.Printf("Fatal error %s while creating the report.Receiver, aborting\n", err)
	}
	receiver.Start()

	proxy, _ := db.NewProxy(resolver, measureCh)
	proxy.Start()

	go gqlServer(":8081", generatedConf)

	quitCh := make(chan int)
	<-quitCh
}

func gqlServer(serveAddr string, conf generated.Config) {
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000", "*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(conf))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				fmt.Printf("Host: %s\n", r.Host)
				//return r.Host == "http://localhost:3000"
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("Testbed Manager Playground", "/query"))
	router.Handle("/query", srv)
	fmt.Printf("connect to https://%s/ for playground\n", serveAddr)
	err := http.ListenAndServe(serveAddr, router)
	if err != nil {
		panic(err)
	}
}
