package smartfanunit

import (
	"github.com/xvzf/computeblade-agent/pkg/smartfanunit/proto"
)

const (
	Baudrate = 115200
)

func MatchCmd(cmd proto.Command) func(any) bool {
	return func(pktAny any) bool {
		pkt, ok := pktAny.(proto.Packet)
		if !ok {
			return false
		}
		if pkt.Command == cmd {
			return true
		}
		return false
	}
}
