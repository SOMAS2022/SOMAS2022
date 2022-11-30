package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

// Agent template for Team6
type Team6Agent struct {
	bravery int
}

func (r Team6Agent) Default() decision.FightAction {
	return decision.Attack
}

func (r Team6Agent) HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction {
	return decision.Attack
}

func (r *Team6Agent) CreateManifesto(view *state.View, baseAgent BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(true, false, 10, 50)
	return manifesto
}

func (r *Team6Agent) HandleConfidencePoll(view *state.View, baseAgent BaseAgent) decision.Intent {
	return decision.Abstain
}

func (r *Team6Agent) HandleFightInformation(_ message.TaggedMessage, _ *state.View, agent BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) {
	agent.Log(logging.Trace, logging.LogField{"bravery": r.bravery, "hp": agent.ViewState().Hp}, "Cowering")
}

func (r *Team6Agent) HandleFightRequest(_ message.TaggedMessage, _ *state.View, _ *immutable.Map[commons.ID, decision.FightAction]) message.Payload {
	return nil
}

func (r *Team6Agent) CurrentAction() decision.FightAction {
	return decision.Attack
}

func (r *Team6Agent) HandleElectionBallot(view *state.View, _ BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	return decision.Ballot{}
}

func NewTeam6Agent() Strategy {
	return &Team6Agent{bravery: 0}
}
