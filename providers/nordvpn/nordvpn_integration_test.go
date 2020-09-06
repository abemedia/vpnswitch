// +build integration

package nordvpn_test

import (
	"os"
	"testing"

	"github.com/abemedia/vpnswitch"
	"github.com/abemedia/vpnswitch/internal/testutils"
	"github.com/abemedia/vpnswitch/providers/nordvpn"
)

func TestNordVPNIntegration(t *testing.T) {
	s, err := nordvpn.New(os.Getenv("NORDVPN_USER"), os.Getenv("NORDVPN_PASSWORD"), "udp")
	if err != nil {
		t.Fatal(err)
	}

	c := new(vpnswitch.Client)
	c.Add(s...)

	testutils.IntegrationTest(t, c)
}
