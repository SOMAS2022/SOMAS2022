package team5

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"sort"
	"sync"

	"github.com/benbjohnson/immutable"
)

var (
	wg         sync.WaitGroup
	once       sync.Once
	totalAgent uint // All agents from the start of the game
)

// FindBestStrategy ...
func FindBestStrategy(view state.View) *immutable.Map[commons.ID, decision.FightAction] {
	level := view.CurrentLevel()
	monsterHealth := view.MonsterHealth()
	monsterAttack := view.MonsterAttack()
	agentState := view.AgentState()

	// iterator get agent value
	agents := getAgents(agentState)

	//all agents
	once.Do(func() {
		totalAgent = uint(len(agents))
	})
	nowAgents := uint(len(agents))

	// Sum all agents attack and defense
	allAttack := uint(0)
	allDefence := uint(0)
	for _, item := range agents {
		allAttack += item.Attack
		allDefence += item.Defense
	}

	// Checking the 2 criteria (total attack > monster health and total shield > monster attack)
	firstResult := handleFirstResult(allAttack, monsterHealth, allDefence, monsterAttack, agents)
	if firstResult != nil {
		return firstResult
	}

	// sort by Hp
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Hp > agents[j].Hp
	})

	// Scaling attack
	agent1 := BuildNewAgent(agents, 0.5)
	agent2 := BuildNewAgent(agents, 0.75)
	agent3 := BuildNewAgent(agents, 1)
	agent4 := BuildNewAgent(agents, 1.5)
	agent5 := BuildNewAgent(agents, 2)

	// handle service
	result1 := getAgentsHandleResult(monsterHealth, monsterAttack, agent1)
	result2 := getAgentsHandleResult(monsterHealth, monsterAttack, agent2)
	result3 := getAgentsHandleResult(monsterHealth, monsterAttack, agent3)
	result4 := getAgentsHandleResult(monsterHealth, monsterAttack, agent4)
	result5 := getAgentsHandleResult(monsterHealth, monsterAttack, agent5)

	resultAll := make([]*Result, 0)
	resultAll = append(resultAll, result1...)
	resultAll = append(resultAll, result2...)
	resultAll = append(resultAll, result3...)
	resultAll = append(resultAll, result4...)
	resultAll = append(resultAll, result5...)

	return getResult(nowAgents, level, resultAll, agents)
}

func getAgentsHandleResult(monsterHealth, monsterAttack uint, agent []*Agent) []*Result {
	result := make([]*Result, len(agent))
	for i := 1; i <= len(agent); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result[i-1] = handleStrategyService(monsterHealth, monsterAttack, agent[0:i])
		}(i)
	}
	wg.Wait()
	return result
}

func handleFirstResult(allAttack, monsterHealth, allDefence, monsterAttack uint, agents []*Agent) *immutable.Map[commons.ID, decision.FightAction] {
	// allAttack > monsterHealth ,success
	if allAttack > monsterHealth {
		return Attack(allAttack, monsterHealth, agents)
	}

	// allDefence > monsterAttack, success
	if allDefence > monsterAttack {
		return Defense(allDefence, monsterAttack, agents)
	}
	return nil
}

func Attack(allAttack, monsterHealth uint, agents []*Agent) *immutable.Map[commons.ID, decision.FightAction] {
	// Rank according to attack
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Attack > agents[j].Attack
	})
	sumAttack := 0
	var index int
	// Find the iteration where total attack > monster resilience
	for i, item := range agents {
		sumAttack += int(item.Attack)
		// all agents included in Sum will attack
		agents[i].Action = uint(decision.Attack)
		if sumAttack >= int(monsterHealth) {
			index = i
			break
		}
	}
	// If not the last agent
	for i, item := range agents {
		if i <= index {
			continue
		}
		// Agents not included in sum will cower
		item.Action = uint(decision.Cower)
		agents[i].Action = uint(decision.Cower)
	}
	return ConvertToImmutable(agents, agents)
}

func Defense(allDefence, monsterAttack uint, agents []*Agent) *immutable.Map[commons.ID, decision.FightAction] {
	// Rank according to shield
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Defense > agents[j].Defense
	})
	sumDefence := 0
	var index int
	// Find the iteration where total shield > monster attack
	for i, item := range agents {
		sumDefence += int(item.Defense)
		agents[i].Action = uint(decision.Defend)
		if sumDefence >= int(monsterAttack) {
			index = i
			break
		}
	}
	// If not last agent
	if index < len(agents)-1 {
		for i, item := range agents {
			if i <= index {
				continue
			}
			// all agents not in sum will attack
			item.Action = uint(decision.Attack)
			agents[i].Action = uint(decision.Attack)
		}
		return ConvertToImmutable(agents, agents)
	}
	return nil
}

// return map
func getAgents(agentState immutable.Map[commons.ID, state.HiddenAgentState]) []*Agent {
	agentIterator := agentState.Iterator()
	// iterator get agent value
	agents := make([]*Agent, 0)
	for !agentIterator.Done() {
		key, value, _ := agentIterator.Next()
		agent := &Agent{
			ID:      key,
			Hp:      uint(value.Hp),
			Attack:  value.Attack + value.BonusAttack,
			Defense: value.Defense + value.BonusDefense,
		}
		agents = append(agents, agent)
	}
	return agents
}
