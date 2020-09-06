package testutils

import (
	"testing"

	"github.com/abemedia/vpnswitch"
	"github.com/abemedia/vpnswitch/internal"
)

func IntegrationTest(t *testing.T, c *vpnswitch.Client) {
	oldIP, err := internal.GetIP()
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Start(); err != nil {
		t.Fatal(err)
	}

	newIP, err := internal.GetIP()
	if err != nil {
		t.Fatal(err)
	}

	if newIP == oldIP {
		t.Fatal("didn't change IP")
	}

	if err = c.Next(); err != nil {
		t.Fatal(err)
	}

	nextIP, err := internal.GetIP()
	if err != nil {
		t.Fatal(err)
	}

	if nextIP == newIP || nextIP == oldIP {
		t.Fatal("didn't change IP")
	}

	if err = c.Stop(); err != nil {
		t.Fatal(err)
	}

	if nextIP, err = internal.GetIP(); err != nil {
		t.Fatal(err)
	}

	if nextIP != oldIP {
		t.Fatal("should have original IP")
	}
}
