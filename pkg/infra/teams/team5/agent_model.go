package team5

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type Agent struct {
	ID      string
	Hp      uint
	Attack  uint
	Defense uint
	Action  uint
}

type Result struct {
	Death  int      // Number of Agent deaths
	Damage uint     // Total damage to Agents
	Agents []*Agent // Fighting agents and fight decision
}

// converting result to map
func ConvertToImmutable(agents []*Agent, agentsAll []*Agent) *immutable.Map[commons.ID, decision.FightAction] {
	b := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	// Agents with fight decision
	agentMap := make(map[string]struct{})
	for _, item := range agents {
		agentMap[item.ID] = struct{}{}
		b.Set(item.ID, decision.FightAction(item.Action))
	}
	// Agents without fight decision will cower
	for _, item := range agentsAll {
		if _, ok := agentMap[item.ID]; !ok {
			b.Set(item.ID, decision.Cower)
		}
	}
	return b.Map()
}
