# vpnswitch: A VPN Switcher Library for Go

Uses OpenVPN to connect to and switch VPNs.

## Caveats

Requires sudo.

## Example

```go
package main

import (
	"log"
	"time"

	"github.com/abemedia/vpnswitch"
	"github.com/abemedia/vpnswitch/providers/nordvpn"
)

func main() {
	s, err := nordvpn.New("user", "password", "udp")
	if err != nil {
		log.Fatal(err)
	}

	c := new(vpnswitch.Client)
	c.Add(s...)
    
	err = c.Start() // connect to VPN
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)
    
	err = c.Next()// switch to next VPN server
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)
    
	err = c.Stop()// disconnect
	if err != nil {
		log.Fatal(err)
	}
}
```
