package commons5

import (
	"infra/game/commons"
	"infra/game/state"
)

// import pool when its merged
type Item struct {
	id    commons.ItemID
	value uint
}

type Loot struct {
	weapons        *commons.ImmutableList[Item]
	shields        *commons.ImmutableList[Item]
	hpPotions      *commons.ImmutableList[Item]
	staminaPotions *commons.ImmutableList[Item]
}

type Agents struct {
}

func FetchLoot(view *state.View) Loot {
	var Loot Loot
	//read loot information into Loot struct
	return Loot
}

func FetchAllAgents(view *state.View) Agents {
	var allAgents Agents
	//read all agent information
	return allAgents
}

func AliveAgents(allAgents Agents) Agents {
	var aliveAgents Agents
	//gather alive agent info to shrink space
	return aliveAgents
}

func AgentsPerm(aliveAgents Agents) []Agents {
	var agentsperm []Agents
	//permutate through all agents combination
	return agentsperm
}
