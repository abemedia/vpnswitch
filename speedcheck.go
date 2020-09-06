package vpnswitch

import (
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/tatsushid/go-fastping"
)

type serverList struct {
	servers []Server
	ipMap   map[string]int
	latency map[string]time.Duration
}

func newServerList() *serverList {
	return &serverList{
		servers: make([]Server, 0),
		ipMap:   make(map[string]int),
		latency: make(map[string]time.Duration),
	}
}

func (sl serverList) Len() int {
	return len(sl.servers)
}

func (sl serverList) Less(i, j int) bool {
	return sl.latency[sl.servers[i].Host()] < sl.latency[sl.servers[j].Host()]
}

func (sl serverList) Swap(i, j int) {
	sl.servers[i], sl.servers[j] = sl.servers[j], sl.servers[i]
}

func speedcheck(servers []Server, timeout time.Duration) ([]Server, error) {
	sl := newServerList()

	p := fastping.NewPinger()
	p.MaxRTT = timeout
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		s := servers[sl.ipMap[addr.String()]]
		sl.latency[s.Host()] = rtt
		sl.servers = append(sl.servers, s)
	}

	for i, s := range servers {
		if err := p.AddIP(s.Host()); err != nil {
			// couldn't parse IP - try resolving hostname
			var ip *net.IPAddr
			if ip, err = net.ResolveIPAddr("ip4:icmp", s.Host()); err != nil {
				fmt.Println(err)
				continue
			}
			sl.ipMap[ip.String()] = i // add IP to map (we currently have a hostname)
			p.AddIPAddr(ip)
		} else {
			sl.ipMap[s.Host()] = i
		}
	}

	if err := p.Run(); err != nil {
		return nil, err
	}

	if len(sl.servers) == 0 {
		return nil, fmt.Errorf("no servers with latency below %s", p.MaxRTT)
	}

	sort.Sort(sl) // sort servers by latency

	return sl.servers, nil
}
