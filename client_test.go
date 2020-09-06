package vpnswitch

import (
	"testing"
)

func TestClientRefresh(t *testing.T) {
	c := &Client{}

	configURLs := []string{
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/uk1780.nordvpn.com.tcp443.ovpn",
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/au570.nordvpn.com.udp1194.ovpn",
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/al19.nordvpn.com.udp1194.ovpn",
		"https://downloads.nordcdn.com/configs/files/ovpn_legacy/servers/uk1890.nordvpn.com.udp1194.ovpn",
	}

	hosts := make([]string, len(configURLs))

	for i, url := range configURLs {
		s, err := NewServerFromURL(url, nil)
		if err != nil {
			t.Fatal(err)
		}
		c.Add(s)
		hosts[i] = s.Host()
	}

	err := c.Refresh()
	if err != nil {
		t.Fatal(err)
	}

	// check if order changed
	var hasChanged bool
	for i, s := range c.activeServers {
		if s.Host() != hosts[i] {
			hasChanged = true
		}
	}
	if !hasChanged {
		t.Error("refresh didn't change server order")
	}
}
