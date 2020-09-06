package nordvpn_test

import (
	"testing"

	"github.com/abemedia/vpnswitch"
	"github.com/abemedia/vpnswitch/providers/nordvpn"
)

func TestNordVPN(t *testing.T) {
	s, err := nordvpn.New("", "", "")
	if err != nil {
		t.Fatal(err)
	}

	if len(s) == 0 {
		t.Fatal("no servers found")
	}

	c := new(vpnswitch.Client)
	c.Add(s...)

	if err := c.Refresh(); err != nil {
		t.Fatal(err)
	}
}
