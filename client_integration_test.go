// +build integration

package vpnswitch_test

import (
	"os"
	"testing"

	"github.com/abemedia/vpnswitch"
	"github.com/abemedia/vpnswitch/internal/testutils"
)

var auth = &vpnswitch.Credentials{
	Username: os.Getenv("NORDVPN_USER"),
	Password: os.Getenv("NORDVPN_PASSWORD"),
}

func TestClientIntegration(t *testing.T) {
	configURLs := []string{
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/uk1780.nordvpn.com.tcp443.ovpn",
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/au570.nordvpn.com.udp1194.ovpn",
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/al19.nordvpn.com.udp1194.ovpn",
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/uk1890.nordvpn.com.udp1194.ovpn",
	}

	c := new(vpnswitch.Client)
	for _, url := range configURLs {
		s, err := vpnswitch.NewServerFromURL(url, auth)
		if err != nil {
			t.Fatal(err)
		}
		c.Add(s)
	}

	testutils.IntegrationTest(t, c)
}
