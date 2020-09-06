package vpnswitch

import (
	"strconv"
	"strings"
)

type Option interface{}

type ConfigOption func(*strings.Builder)

func applyConfigOptions(config string, options []Option) string {
	b := new(strings.Builder)
	b.WriteString(config)

	for _, o := range options {
		if opt, ok := o.(ConfigOption); ok {
			opt(b)
		}
	}

	return b.String()
}

func WithProto(proto string) Option {
	return func(b *strings.Builder) {
		b.WriteString("\nproto ")
		b.WriteString(proto)
	}
}

func WithRemote(host string, port int) Option {
	return func(b *strings.Builder) {
		b.WriteString("\nremote ")
		b.WriteString(host)
		b.WriteString(" ")
		b.WriteString(strconv.Itoa(port))
	}
}
