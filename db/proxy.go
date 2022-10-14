package db

import (
	"log"
	"strconv"
	"testbed-monitor/graph"
	"testbed-monitor/graph/model"
	"testbed-monitor/report"
)

type Proxy struct {
	resolver *graph.Resolver
	statusCh chan *report.StatusReport
}

func NewProxy(resolver *graph.Resolver, statusCh chan *report.StatusReport) (*Proxy, error) {
	proxy := new(Proxy)
	proxy.resolver = resolver
	proxy.statusCh = statusCh
	return proxy, nil
}

func (proxy *Proxy) Start() {
	go func() {
		for {
			select {
			case s := <-proxy.statusCh:
				proxy.processReport(s)
			}
		}
	}()
}

func (proxy *Proxy) processReport(s *report.StatusReport) {
	ip := s.Tower
	raw := proxy.resolver.GetHost(ip)
	if raw == nil {
		status := copyHostStatus(ip, nil, s)
		proxy.resolver.CommitHost(status)
	} else {
		x := raw.(*model.HostStatus)
		log.Printf("Got raw: %v\n", x)
		proxy.resolver.CommitHost(copyHostStatus(ip, x, s))
	}
}

func copyHostStatus(ip string, old *model.HostStatus, s *report.StatusReport) *model.HostStatus {
	var newStatus model.HostStatus
	if old == nil {
		newStatus = model.HostStatus{
			BoardReached: "-",
			TowerReached: "-",
			BootTime:     "-",
			Reboots:      0,
			UsedRAM:      "-",
			UsedDisk:     "-",
			CPU:          "-",
			Reachable:    true,
			Temperature:  "-",
		}

	} else {
		newStatus = *old
	}
	newStatus.ID = ip
	if boardReached := s.ArduinoReached; boardReached != "" {
		newStatus.BoardReached = boardReached
	}
	if towerReached := s.TowerReached; towerReached != "" {
		newStatus.TowerReached = towerReached
	}
	if bootTime := s.BootTime; bootTime != "" {
		newStatus.BootTime = bootTime
	}
	if reboots := int(s.Reboots); reboots != 0 {
		newStatus.Reboots = reboots
	}
	if usedRAM := s.UsedRAM; usedRAM > 0 {
		newStatus.UsedRAM = strconv.FormatInt(usedRAM, 10) + "/" + strconv.FormatInt(s.TotalRAM, 10)
	}
	if usedDisk := s.UsedDisk; usedDisk > 0 {
		newStatus.UsedDisk = strconv.FormatInt(usedDisk, 10) + "/" + strconv.FormatInt(s.TotalDisk, 10)
	}
	if cpu := s.Cpu; cpu > 0 {
		newStatus.CPU = strconv.FormatInt(cpu, 10)
	}
	if reachable := s.Reachable; reachable == false {
		newStatus.Reachable = false
	} else {
		newStatus.Reachable = true
	}
	if temperature := s.Temperature; temperature > 0 {
		newStatus.Temperature = strconv.FormatInt(temperature, 10)
	}
	return &newStatus
}
