//go:build linux
// +build linux

package probers

import (
	"bytes"
	"fmt"
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
		fmt.Printf("Probe was not able to reach %s\n", target)
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

	tokens := strings.Split(out.String(), "\n")
	for _, token := range tokens {

		if strings.Contains(token, "packets transmitted") {
			fmt.Println("Parsing: " + token)
			// 3 packets transmitted, 0 received, +3 errors, 100% packet loss, time 2002ms
			loss = extractStats(token)
			if loss == 100 {
				break
			}
		}
		if strings.Contains(token, "rtt min/avg/max/mdev") {
			fmt.Println("Parsing: " + token)
			// rtt min/avg/max/mdev = 20.182/46.469/76.346/23.069 ms
			rtt = extractRTT(token)
		}
	}

	fmt.Printf("Loss: %v Avg rtt: %v\n", packetsReceived, rtt)
	replyCh <- &pb.PingReply{
		Reachable:      loss != 0,
		AvgRtt:         int32(rtt),
		LostPercentage: int32(loss),
	}
}

func extractReceived(token string) int {
	// 0 received
	s := strings.Replace(token, "received", "", -1)
	s1 := strings.TrimSpace(s)
	received, _ := strconv.Atoi(strings.TrimSpace(s1))
	return received
}

func extractRTT(token string) int {
	// rtt min/avg/max/mdev = 20.182/46.469/76.346/23.069 ms
	s := strings.Replace(token, "rtt min/avg/max/mdev =", "", -1)
	s1 := strings.Replace(s, "ms", "", -1)
	s2 := strings.TrimSpace(s1)
	fmt.Println("S2: " + s2)
	// 20.182/46.469/76.346/23.069
	tokens := strings.Split(s2, "/")
	fmt.Println("RTT: " + tokens[1])
	rtt, _ := strconv.ParseFloat(tokens[1], 32)
	return int(rtt)
}

func extractStats(token string) int {
	// not reachable
	// 3 packets transmitted, 0 received, +3 errors, 100% packet loss, time 2002ms
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
			// 100% packet loss
			s := strings.Replace(token, "packet loss", "", -1)
			s1 := strings.TrimSpace(s)
			loss, _ := strconv.Atoi(strings.TrimSpace(s1))
			return loss
		}
	}
	// something went wrong?
	return 100
}
