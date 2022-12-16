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
	totalAgent uint // 所有存在的agent
)

// FindBestStrategy ...
func FindBestStrategy(view state.View) *immutable.Map[commons.ID, decision.FightAction] {
	level := view.CurrentLevel()
	monsterHealth := view.MonsterHealth()
	monsterAttack := view.MonsterAttack()
	agentState := view.AgentState()

	// iterator get agent value
	agents := getAgents(agentState)

	// 所有人
	once.Do(func() {
		totalAgent = uint(len(agents))
	})
	nowAgents := uint(len(agents))

	// get allAttack & allDefence
	allAttack := uint(0)
	allDefence := uint(0)
	for _, item := range agents {
		allAttack += item.Attack
		allDefence += item.Defense
	}

	// 进行初步检查，看看是否总攻击》怪物血，或者总防御大于怪物攻击
	firstResult := handleFirstResult(allAttack, monsterHealth, allDefence, monsterAttack, agents)
	if firstResult != nil {
		return firstResult
	}

	// sort by Hp
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Hp > agents[j].Hp
	})

	// build new agent (weight attack 0.5,0.75,1,1.5,2)
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
	// 根据攻击排序
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Attack > agents[j].Attack
	})
	sumAttack := 0
	var index int
	// 找出总攻击>怪物血量的那个人的位置
	for i, item := range agents {
		sumAttack += int(item.Attack)
		// 全部攻击
		agents[i].Action = uint(decision.Attack)
		// 找到指定位置的人了
		if sumAttack >= int(monsterHealth) {
			index = i
			break
		}
	}
	// 如果不是最后一个人
	for i, _ := range agents {
		if i <= index {
			continue
		}
		// 指定位置之后的人都逃跑
		agents[i].Action = uint(decision.Cower)
	}
	return ConvertToImmutable(agents, agents)
}

func Defense(allDefence, monsterAttack uint, agents []*Agent) *immutable.Map[commons.ID, decision.FightAction] {
	// 根据防御排序
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Defense > agents[j].Defense
	})
	sumDefence := 0
	var index int
	// 找出总防御>怪物攻击的那个人的位置
	for i, item := range agents {
		sumDefence += int(item.Defense)
		// 全部防御
		agents[i].Action = uint(decision.Defend)
		// 找到指定位置的人了
		if sumDefence >= int(monsterAttack) {
			index = i
			break
		}
	}
	// 如果不是最后一个人
	if index < len(agents)-1 {
		for i, _ := range agents {
			if i <= index {
				continue
			}
			// 指定位置之后的人都攻击
			agents[i].Action = uint(decision.Attack)
		}
		return ConvertToImmutable(agents, agents)
	}
	return nil
}

// 返回查找结果
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
