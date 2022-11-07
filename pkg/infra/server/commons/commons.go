package commons

import (
	"infra/server/decision"
	"infra/server/message"
	"infra/server/state"
)

type Communication struct {
	Peer     []chan message.Message
	Receiver <-chan state.State
	Sender   chan<- decision.Decision
}

func SaturatingSub(x uint, y uint) uint {
	res := x - y
	var val uint
	if res <= x {
		val = 1
	}
	res &= -val
	return res
}
