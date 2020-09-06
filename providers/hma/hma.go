package hma

import (
	"net"
	"strconv"

	"github.com/abemedia/vpnswitch"
	"github.com/abemedia/vpnswitch/internal"
)

func New(username, password, proto string) ([]vpnswitch.Server, error) {
	var port int
	if proto == "tcp" {
		port = 443
	} else {
		proto = "udp"
		port = 553
	}

	auth := &vpnswitch.Credentials{
		Username: username,
		Password: password,
	}

	// config, err := internal.GetString("https://securenetconnection.com/vpnconfig/openvpn-template.ovpn")
	// if err != nil {
	// 	return nil, err
	// }

	p := vpnswitch.NewProvider(config, port, auth, vpnswitch.WithProto(proto))

	servers, err := internal.GetServers("https://vpn.hidemyass.com/vpn-config/l2tp/")
	if err != nil {
		return nil, err
	}
	for _, s := range servers {
		ips, _ := net.LookupHost(s[1])
		for i, ip := range ips {
			name := s[0] + " " + strconv.Itoa(i+1)
			p.Add(name, ip)
		}
	}

	return p.Servers(), nil
}
