package task

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"strconv"
	pb "testbed-monitor/pinger"
	"time"
)

type Ping struct {
}

func NewPingTask() *Ping {
	task := new(Ping)
	return task
}

func (task *Ping) Start(HOST_IP string) {
	go task.execute(HOST_IP)
}

func (task *Ping) execute(HOST_IP string) {
	conn, err := grpc.Dial(HOST_IP+":"+os.Getenv("GRPC_PINGER_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPingerClient(conn)
	// Contact the server and print out its response.
	log.Printf("Contacting server, this may take a while...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := c.Ping(ctx, &pb.PingRequest{TargetAddress: HOST_IP})
	if err != nil {
		// host monitor was not reachable!
		log.Println("Could not reach the host!")
		var filename = time.Now().Format("2006-01-02_1504")
		filename += "_Ping.txt"
		err := LogPing(filename, "Could not reach the host: "+HOST_IP)
		log.Fatalf("Could not greet: %v", err)
	}

	if r.GetReachable() {
		pingResponse := "Reached host: " + HOST_IP + "  Packet Loss: " + strconv.Itoa(int(r.GetLostPercentage())) + "  Rtt: " + strconv.Itoa(int(r.GetAvgRtt()))
		var filename = time.Now().Format("2006-01-02_1504")
		filename += "_Ping.txt"
		LogPing(filename, pingResponse)
		log.Printf("%s/n", pingResponse)
	} else {
		var filename = time.Now().Format("2006-01-02_1504")
		filename += "_Ping.txt"
		LogPing(filename, "Reached the host: "+HOST_IP+", but could not reach endpoint.")
		log.Printf("Endpoint was not reachable")
	}
}

// LogPing This function writes report data to a file
func LogPing(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
		log.Printf("Error creating file: %s", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
		log.Printf("Error writing data to file: %s", err)
	}
	return file.Sync()
}
