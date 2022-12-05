package commons5

import (
	"infra/game/commons"
	"infra/game/state"
)

type Loot struct {
	Shields  map[commons.ID]uint
	Weapons  map[commons.ID]uint
	HPpotion []uint
	STpotion []uint
}

type Agents struct {
}

func FetchLoot(view *state.View) Loot {
	var Loot Loot
	//read loot infomation into Loot struct
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
