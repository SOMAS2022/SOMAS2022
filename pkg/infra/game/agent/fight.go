package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/message/proposal"

	"github.com/benbjohnson/immutable"
)

type Fight interface {
	HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction])
	HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform
	FightResolution(baseAgent BaseAgent, prop commons.ImmutableList[proposal.Rule[decision.FightAction]], proposedActions immutable.Map[commons.ID, decision.FightAction]) immutable.Map[commons.ID, decision.FightAction]
	HandleFightProposal(proposal message.Proposal[decision.FightAction], baseAgent BaseAgent) decision.Intent
	// HandleFightProposalRequest only called as leader
	HandleFightProposalRequest(proposal message.Proposal[decision.FightAction], baseAgent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) bool
	FightActionNoProposal(baseAgent BaseAgent) decision.FightAction
	FightAction(baseAgent BaseAgent, proposedAction decision.FightAction, acceptedProposal message.Proposal[decision.FightAction]) decision.FightAction
}
