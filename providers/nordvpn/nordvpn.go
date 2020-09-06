package nordvpn

import (
	"regexp"
	"strings"

	"github.com/abemedia/vpnswitch"
	"github.com/abemedia/vpnswitch/internal"
)

func New(username, password, proto string) ([]vpnswitch.Server, error) {
	var port int
	if proto == "tcp" {
		port = 443
	} else {
		proto = "udp"
		port = 1194
	}

	auth := &vpnswitch.Credentials{
		Username: username,
		Password: password,
	}

	p := vpnswitch.NewProvider(config, port, auth, vpnswitch.WithProto(proto))

	servers, err := getServers()
	if err != nil {
		return nil, err
	}
	for _, s := range servers {
		p.Add(s[0], s[1])
	}

	return p.Servers(), nil
}

func getServers() ([][2]string, error) {
	s, err := internal.GetString("https://nordvpn.com/ovpn/")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[^/](([a-z]{2})-?([a-z]+)?(\d+)\.nordvpn.com)[^.]`)
	matches := re.FindAllStringSubmatch(s, -1)
	result := make([][2]string, len(matches))

	for i, m := range matches {
		b := new(strings.Builder)
		b.WriteString(countryFromAlpha2(m[2]))
		if m[3] != "" {
			b.WriteString(" - ")
			b.WriteString(countryFromAlpha2(m[3]))
		}
		b.WriteRune(' ')
		b.WriteString(m[4])
		result[i] = [2]string{b.String(), m[1]}
	}

	return result, nil
}
