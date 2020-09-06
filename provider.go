package vpnswitch

import (
	"strconv"
)

type providerServer struct {
	p    *Provider
	name string
	host string
}

func (s *providerServer) Config() string {
	return s.p.config + "\nremote " + s.host + " " + strconv.Itoa(s.p.port)
}

func (s *providerServer) Credentials() Credentials {
	return *s.p.auth
}

func (s *providerServer) Host() string {
	return s.host
}

func (s *providerServer) Name() string {
	return s.name
}

type Provider struct {
	auth    *Credentials
	config  string
	port    int
	servers []Server
}

func NewProvider(config string, port int, auth *Credentials, options ...Option) *Provider {
	if options != nil {
		config = applyConfigOptions(config, options)
	}

	return &Provider{
		auth:   auth,
		config: config,
		port:   port,
	}
}

func (p *Provider) Servers() []Server {
	return p.servers
}

func (p *Provider) Add(name, host string) {
	server := &providerServer{
		p:    p,
		name: name,
		host: host,
	}
	p.servers = append(p.servers, server)
}
