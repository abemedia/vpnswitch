package vpnswitch

type Server interface {
	Config() string
	Host() string
	Name() string
	Credentials() Credentials
}

type Credentials struct {
	Username string
	Password string
}
