package vpnswitch

import (
	"errors"
	"regexp"

	"github.com/abemedia/vpnswitch/internal"
)

type server struct {
	config string
	host   string
	auth   *Credentials
}

func NewServerFromURL(configURL string, auth *Credentials, options ...Option) (Server, error) {
	config, err := internal.GetString(configURL)
	if err != nil {
		return nil, err
	}

	return NewServer(config, auth)
}

var remoteRegex = regexp.MustCompile(`(?m)^remote (.*) \d+$`)

func NewServer(config string, auth *Credentials, options ...Option) (Server, error) {
	if options != nil {
		config = applyConfigOptions(config, options)
	}

	match := remoteRegex.FindStringSubmatch(config)
	if len(match) == 0 {
		return nil, errors.New("config contains no remote")
	}

	return &server{auth: auth, config: config, host: match[1]}, nil
}

func (s *server) Config() string {
	return s.config
}

func (s *server) Name() string {
	return s.Host()
}

func (s *server) Host() string {
	return s.host
}

func (s *server) Credentials() Credentials {
	return *s.auth
}
