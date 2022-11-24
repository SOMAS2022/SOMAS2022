package agent

import (
	"infra/game/commons"
	"infra/game/message"

	"github.com/benbjohnson/immutable"
)

type Communication struct {
	receipt <-chan message.TaggedMessage
	peer    immutable.Map[commons.ID, chan<- message.TaggedMessage]
}

func NewCommunication(receipt <-chan message.TaggedMessage, peer immutable.Map[commons.ID, chan<- message.TaggedMessage]) *Communication {
	return &Communication{receipt: receipt, peer: peer}
}
