package team5

import (
	"infra/game/commons"
	"infra/game/decision"
	"math"

	"github.com/benbjohnson/immutable"
)

func handleStrategyService(monsterHealthy, monsterAttack uint, agents []*Agent) *Result {
	newAgents := CopySlice(agents)
	totalAttack := uint(0)
	totalDefense := uint(0)
	for _, item := range newAgents {
		totalAttack += item.Attack
		totalDefense += item.Defense
	}

	if totalAttack == 0 {
		totalAttack = 1
	}
	// total rounds
	rounds := monsterHealthy / totalAttack
	if monsterHealthy%totalAttack != 0 {
		rounds += 1
	}

	// total Damage
	totalDamage := (monsterAttack - totalDefense) * rounds
	divideDamage := totalDamage / uint(len(newAgents))
	death := 0
	for _, item := range newAgents {
		if divideDamage >= item.Hp {
			death++
		}
	}
	res := &Result{
		Damage: divideDamage,
		Death:  death,
		Agents: agents,
	}
	return res
}

func getResult(nowAgents uint, level uint, resultAll []*Result, agents []*Agent) *immutable.Map[commons.ID, decision.FightAction] {
	minDeathRes := &Result{Death: math.MaxInt}
	minDamageRes := &Result{Damage: math.MaxUint}

	// 死亡人数和总人数不等时，找出最少死亡的
	// 找出最少伤害的
	suviveRate := float64(nowAgents) / float64(totalAgent)
	population := GetPopulation(level)
	if suviveRate >= population {
		for _, item := range resultAll {
			if item.Death != len(item.Agents) && item.Damage < minDamageRes.Damage {
				minDamageRes = item
			}
		}
		// 有符合条件的
		if len(minDamageRes.Agents) > 0 {
			return ConvertToImmutable(minDamageRes.Agents, agents)
		}
	} else {
		for _, item := range resultAll {
			if item.Death != len(item.Agents) && uint(item.Death) < uint(minDeathRes.Death) {
				minDeathRes = item
			}
		}
		// 有符合条件的
		if len(minDeathRes.Agents) > 0 {
			return ConvertToImmutable(minDeathRes.Agents, agents)
		}
	}
	for i := range agents {
		agents[i].Action = uint(decision.Attack)
	}
	return ConvertToImmutable(agents, agents)
}
