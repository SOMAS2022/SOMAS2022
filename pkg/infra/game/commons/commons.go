package commons

import (
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
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
