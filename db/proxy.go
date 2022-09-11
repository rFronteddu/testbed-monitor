package db

import (
	"testbed-monitor/graph"
	"testbed-monitor/graph/model"
	"testbed-monitor/measure"
	"testbed-monitor/report"
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
	address := m.Strings[report.SENSOR_IP]
	raw := proxy.resolver.GetHost(address)
	if raw == nil {
		status := getStatusFromMeasure(address, nil, m)
		proxy.resolver.CommitHost(status)
	} else {
		x := raw.(*model.HostStatus)
		log.Printf("Got raw: %v\n", x)
		proxy.resolver.CommitHost(getStatusFromMeasure(address, x, m))
	}
}

func getStatusFromMeasure(ip string, old *model.HostStatus, m *measure.Measure) *model.HostStatus {
	var newStatus model.HostStatus
	if old == nil {
		newStatus = model.HostStatus{
			BootTime: "-",
			HostName: "-",
			HostID:   "-",
			Os:       "-",
			Platform: "-",
			Kernel:   "-",
		}

	} else {
		newStatus = *old
	}

	newStatus.Time = time.Now().Format(time.RFC1123)

	newStatus.ID = ip
	if cpu := int(m.Integers["CPU_AVG"]); cpu != 0 {
		newStatus.CPUAvg = cpu
	}
	if diskUsage := int(m.Integers["DISK_USAGE"]); diskUsage != 0 {
		newStatus.DiskUsagePercent = diskUsage
	}
	if diskFree := int(m.Integers["DISK_FREE"]); diskFree != 0 {
		newStatus.DiskFree = diskFree
	}
	if load1 := int(m.Integers["LOAD_1"]); load1 != 0 {
		newStatus.Load1 = load1
	}
	if load5 := int(m.Integers["LOAD_5"]); load5 != 0 {
		newStatus.Load5 = load5
	}
	if load15 := int(m.Integers["LOAD_15"]); load15 != 0 {
		newStatus.Load15 = load15
	}
	if usedPercent := int(m.Integers["vm_used_percent"]); usedPercent != 0 {
		newStatus.VirtualMemoryUsagePercent = usedPercent
	}
	if free := int(m.Integers["vm_free"]); free != 0 {
		newStatus.VirtualMemoryFree = free
	}
	if hostID := m.Strings["host_id"]; hostID != "" {
		newStatus.HostID = hostID
	}
	if hostName := m.Strings["host_name"]; hostName != "" {
		newStatus.HostName = hostName
	}
	if os := m.Strings["os"]; os != "" {
		newStatus.Os = os
	}
	if platform := m.Strings["platform"]; platform != "" {
		newStatus.Platform = platform
	}
	if kernel := m.Strings["kernelArch"]; kernel != "" {
		newStatus.Kernel = kernel
	}
	if bootTime := m.Strings["bootTime"]; bootTime != "" {
		newStatus.BootTime = bootTime
	}
	return &newStatus
}
