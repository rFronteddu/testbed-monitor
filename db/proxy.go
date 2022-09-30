package db

import (
	"fmt"
	"log"
	"testbed-monitor/graph"
	"testbed-monitor/graph/model"
	"testbed-monitor/measure"
	"time"
)

type Proxy struct {
	resolver  *graph.Resolver
	measureCh chan *measure.Measure
}

func NewProxy(resolver *graph.Resolver, measureCh chan *measure.Measure) (*Proxy, error) {
	proxy := new(Proxy)
	proxy.resolver = resolver
	proxy.measureCh = measureCh
	return proxy, nil
}

func (proxy *Proxy) Start() {
	go func() {
		timer := time.NewTimer(time.Minute)
		for {
			select {
			case m := <-proxy.measureCh:
				proxy.processReport(m)
			case <-timer.C:
				timer = time.NewTimer(time.Minute)
				hosts, _ := proxy.resolver.GetHosts()
				proxy.processState(hosts)
			}
		}
	}()
}

func (proxy *Proxy) processState(hosts []*model.HostStatus) {

}

func (proxy *Proxy) processReport(m *measure.Measure) {
	fmt.Println("Processing report.")
	address := m.Strings["host_id"]
	raw := proxy.resolver.GetHost(address)
	if raw == nil {
		status := copyHostStatus(address, nil, m)
		proxy.resolver.CommitHost(status)
	} else {
		x := raw.(*model.HostStatus)
		log.Printf("Got raw: %v\n", x)
		proxy.resolver.CommitHost(copyHostStatus(address, x, m))
	}
}

func copyHostStatus(ip string, old *model.HostStatus, m *measure.Measure) *model.HostStatus {
	var newStatus model.HostStatus
	if old == nil {
		newStatus = model.HostStatus{
			BootTime:    "-",
			Reboots:     0,
			Reachable:   true,
			Temperature: "-",
			// define all strings here?????
		}

	} else {
		newStatus = *old
	}
	newStatus.ID = ip
	newStatus.Tower = ip
	if cpu := int(m.Integers["CPU_AVG"]); cpu != 0 {
		newStatus.CPU = string(cpu)
	}
	if diskUsage := int(m.Integers["DISK_USAGE"]); diskUsage != 0 {
		newStatus.UsedDisk = string(diskUsage)
	}
	if usedPercent := int(m.Integers["vm_used_percent"]); usedPercent != 0 {
		newStatus.UsedRAM = string(usedPercent)
	}
	if bootTime := m.Strings["bootTime"]; bootTime != "" {
		newStatus.BootTime = bootTime
	}
	fmt.Println("Host Status copied.")
	return &newStatus
}
