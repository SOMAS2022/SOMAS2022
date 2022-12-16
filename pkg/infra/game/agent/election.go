package agent

import (
	"infra/game/decision"
)

type Election interface {
	CreateManifesto(baseAgent BaseAgent) *decision.Manifesto
	HandleConfidencePoll(baseAgent BaseAgent) decision.Intent
	HandleElectionBallot(baseAgent BaseAgent, params *decision.ElectionParams) decision.Ballot
}
