//go:build darwin

package probers

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
	pb "testbed-monitor/pinger"
)

func ping(target string, replyCh chan *pb.PingReply) {
	var cmd = exec.Command("ping", target, "-c 3")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		// was not able to ping
		log.Printf("Probe was not able to reach %s\n", target)
		replyCh <- &pb.PingReply{
			Reachable:      false,
			AvgRtt:         0,
			LostPercentage: 100,
		}
		return
	}

	var packetsReceived int
	var rtt int
	var loss int
	var reachable = false

	tokens := strings.Split(out.String(), "\n")
	for _, token := range tokens {
		if strings.Contains(token, "unreachable") {
			log.Printf("Destination %s unreachable\n", target)
		}
		if strings.Contains(token, "packets transmitted") {
			loss = extractStats(token)
			if loss == 100 {
				break
			} else {
				reachable = true
			}
		}
		if strings.Contains(token, "round-trip min/avg/max/mdev") {
			rtt = extractRTT(token)
		}
	}

	log.Printf("Loss: %v Avg rtt: %v\n", packetsReceived, rtt)

	replyCh <- &pb.PingReply{
		Reachable:      reachable,
		AvgRtt:         int32(rtt),
		LostPercentage: int32(loss),
	}
}

func extractReceived(token string) int {
	s := strings.Replace(token, "packets received", "", -1)
	s1 := strings.TrimSpace(s)
	received, _ := strconv.Atoi(strings.TrimSpace(s1))
	return received
}

func extractRTT(token string) int {
	s := strings.Replace(token, "round-trip min/avg/max/mdev =", "", -1)
	s1 := strings.Replace(s, "ms", "", -1)
	s2 := strings.TrimSpace(s1)
	log.Println("S2: " + s2)
	tokens := strings.Split(s2, "/")
	log.Println("RTT: " + tokens[1])
	rtt, _ := strconv.ParseFloat(tokens[1], 32)
	return int(rtt)
}

func extractStats(token string) int {
	tokens := strings.Split(token, ",")

	for _, token := range tokens {
		if strings.Contains(token, "received") {
			received := extractReceived(token)
			if received == 0 {
				// was not able to reach endpoint
				return 100
			}
		}
		if strings.Contains(token, "packet loss") {
			s := strings.Replace(token, "packet loss", "", -1)
			s1 := strings.TrimSpace(s)
			loss, _ := strconv.Atoi(strings.TrimSpace(s1))
			return loss
		}
	}
	// something went wrong?
	return 100
}
