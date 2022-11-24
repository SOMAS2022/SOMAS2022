package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/message"
)

type Communication struct {
	receipt <-chan message.TaggedMessage
	peer    immutable.Map[commons.ID, chan<- message.TaggedMessage]
}

func NewCommunication(receipt <-chan message.TaggedMessage, peer immutable.Map[commons.ID, chan<- message.TaggedMessage]) Communication {
	return Communication{receipt: receipt, peer: peer}
}
