package agent

import (
	"infra/game/commons"
	"infra/game/message"
)

type Trust interface {
	CompileTrustMessage(agentMap map[commons.ID]Agent) message.Trust
	HandleTrustMessage(message message.TaggedMessage)
}
