package task

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "testbed-monitor/pinger"
	"time"
)

type Ping struct {
}

func NewPingTask() *Ping {
	task := new(Ping)
	return task
}

func (task *Ping) Start() {
	go task.execute()
}

func (task *Ping) execute() {
	conn, err := grpc.Dial("192.168.101.2:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPingerClient(conn)
	// Contact the server and print out its response.
	log.Printf("Contacting server, this may take a while...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := c.Ping(ctx, &pb.PingRequest{TargetAddress: "192.168.250.99"})
	if err != nil {
		// host monitor was not reachable!
		log.Fatalf("could not greet: %v", err)
	}

	if r.GetReachable() {
		log.Printf("Packet Loss: %v Rtt: %v", r.GetLostPercentage(), r.GetAvgRtt())
	} else {
		log.Printf("Endpoint was not reachable")
	}
}
