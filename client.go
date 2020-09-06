package vpnswitch

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/abemedia/vpnswitch/internal"
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
)

var (
	DefaultTimeout = 10 * time.Second
)

type Client struct {
	Timeout    time.Duration
	MaxLatency time.Duration
	MaxRetries int

	currentServer int
	conn          *loggingCallbacks
	servers       []Server
	activeServers []Server
	lock          sync.Mutex
	session       *openvpn3.Session
}

// Add adds more VPN servers
func (c *Client) Add(s ...Server) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.servers = append(c.servers, s...)
}

func (c *Client) timeout() time.Duration {
	if c.Timeout == 0 {
		return DefaultTimeout
	}
	return c.Timeout
}

// Start connects to VPN
func (c *Client) Start() error {
	if err := c.Refresh(); err != nil {
		return err
	}

	err := c.connect()
	if err != nil {
		for i := 0; i < c.MaxRetries; i++ {
			if err = c.Next(); err != nil {
				fmt.Printf("Error: %s. Trying new server...\n", err.Error())
				continue
			}
			break // no errors
		}
	}
	return err
}

func (c *Client) connect() error {
	s := c.activeServers[c.currentServer]

	c.conn = newHandler()
	config := openvpn3.NewConfig(s.Config())
	c.session = openvpn3.NewSession(config, openvpn3.UserCredentials(s.Credentials()), c.conn)
	c.session.Start()

	for {
		select {
		case event := <-c.conn.Events:
			switch event.Name {
			case "CONNECTED":
				currentIP, err := internal.GetIP()
				if err != nil {
					return err
				}
				fmt.Printf("Connected to %s. New IP: %s\n", s.Name(), currentIP)
				return nil
			case "RESOLVE", "WAIT", "CONNECTING", "GET_CONFIG", "ASSIGN_IP":
				continue
			}
			if event.Fatal {
				c.session.Stop()
				return c.session.Wait()
			}
			if event.Error {
				fmt.Println("Error: ", event.Name, "(", event.Info, ")")
				continue
			}
			fmt.Println("vpn: ", event.Name, "(", event.Info, ")")
			continue
		case <-time.After(c.timeout()):
			c.session.Stop()
			return errors.New("vpn: connection timed out")
		}
	}
}

// Stop disconnects from VPN
func (c *Client) Stop() error {
	s := c.activeServers[c.currentServer]
	c.session.Stop()
	for {
		select {
		case event := <-c.conn.Events:
			switch {
			case event.Name == "DISCONNECTED":
				currentIP, err := internal.GetIP()
				if err != nil {
					return err
				}
				fmt.Printf("Disconnected from %s. New IP: %s\n", s.Name(), currentIP)
				return c.session.Wait()
			case event.Fatal:
				fmt.Println("vpn fatal error: ", event.Name, "(", event.Info, ")")
			case event.Error:
				fmt.Println("vpn error: ", event.Name, "(", event.Info, ")")
			default:
				fmt.Println("vpn event: ", event.Name, "(", event.Info, ")")
			}
		case <-time.After(c.timeout()):
			return errors.New("vpn: connection timed out")
		}
	}
}

// Next switches to the next VPN
func (c *Client) Next() error {
	if err := c.Stop(); err != nil {
		return err
	}

	if c.currentServer < len(c.activeServers)-1 {
		c.currentServer++
	} else {
		if err := c.Refresh(); err != nil {
			return err
		}
	}

	return c.connect()
}

// Refresh pings the servers to update the latency values
func (c *Client) Refresh() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	fmt.Printf("Speed-checking %d VPNs...\n", len(c.servers))

	servers, err := speedcheck(c.servers, c.timeout())
	if err != nil {
		return err
	}

	fmt.Printf("Found %d VPNs with latency below %s.\n", len(servers), c.timeout())

	c.currentServer = 0
	c.activeServers = servers
	return nil
}
