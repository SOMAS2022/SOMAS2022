package team5

import "infra/game/decision"

func CopySlice(agents []*Agent) []*Agent {
	newAgents := make([]*Agent, len(agents))
	for i, item := range agents {
		newAgents[i] = &Agent{
			ID:      item.ID,
			Hp:      item.Hp,
			Attack:  item.Attack,
			Defense: item.Defense,
			Action:  item.Action,
		}
	}
	return newAgents
}

func BuildNewAgent(agents []*Agent, p float64) []*Agent {
	newAgents := make([]*Agent, len(agents))
	for i, item := range agents {
		attack := item.Attack
		defense := item.Defense
		newAgents[i] = &Agent{
			ID: item.ID,
			Hp: item.Hp,
		}
		if uint(float64(item.Attack)*p) >= defense {
			newAgents[i].Attack = attack
			newAgents[i].Defense = 0
			newAgents[i].Action = uint(decision.Attack)
		} else {
			newAgents[i].Attack = 0
			newAgents[i].Defense = defense
			newAgents[i].Action = uint(decision.Defend)
		}
	}
	return newAgents
}

func GetPopulation(customs uint) float64 {
	if customs <= 10 {
		return 0.9
	}
	k := -float64(3) / float64(5)
	b := 96
	return k*float64(customs) + float64(b)
}
