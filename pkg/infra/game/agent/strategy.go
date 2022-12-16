package agent

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

type Strategy interface {
	Fight
	Election
	Loot
	HPPool
	Trade
	// HandleUpdateWeapon return the index of the weapon you want to use in AgentState.weapons
	HandleUpdateWeapon(baseAgent BaseAgent) decision.ItemIdx
	// HandleUpdateShield return the index of the shield you want to use in AgentState.Shields
	HandleUpdateShield(baseAgent BaseAgent) decision.ItemIdx

	UpdateInternalState(baseAgent BaseAgent, fightResult *commons.ImmutableList[decision.ImmutableFightResult], voteResult *immutable.Map[decision.Intent, uint], logChan chan<- logging.AgentLog)
}
