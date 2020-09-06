package vpnswitch

import (
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
)

type callbacks interface {
	openvpn3.Logger
	openvpn3.EventConsumer
	openvpn3.StatsConsumer
}

type loggingCallbacks struct {
	Events chan openvpn3.Event
}

func newHandler() *loggingCallbacks {
	h := new(loggingCallbacks)
	h.Events = make(chan openvpn3.Event)
	return h
}

func (lc *loggingCallbacks) Log(text string) {
	// lines := strings.Split(text, "\n")
	// for _, line := range lines {
	// 	fmt.Println("Openvpn log >>", line)
	// }
}

func (lc *loggingCallbacks) OnEvent(event openvpn3.Event) {
	// fmt.Println("OpenVPN Event", event)
	lc.Events <- event
}

func (lc *loggingCallbacks) OnStats(stats openvpn3.Statistics) {
	// fmt.Printf("Openvpn stats >> %+v\n", stats)
}

var _ callbacks = &loggingCallbacks{}
