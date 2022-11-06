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
