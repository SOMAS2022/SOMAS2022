package hppool

import (
	"sync"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
)

func UpdateHpPool(agentMap map[commons.ID]agent.Agent, globalState *state.State) {
	var wg sync.WaitGroup
	donationChan := make(chan decision.HpPoolDonation, len(agentMap))
	for id, a := range agentMap {
		id := id
		a := a
		aState := globalState.AgentState[id]
		wg.Add(1)
		go func(wait *sync.WaitGroup, donationChan chan decision.HpPoolDonation, agentState state.AgentState) {
			donation := a.HandleDonateToHpPool(agentState)
			donationChan <- decision.HpPoolDonation{AgentID: id, Donation: donation}
			wait.Done()
		}(&wg, donationChan, aState)
	}

	go func(wait *sync.WaitGroup) {
		wait.Wait()
		close(donationChan)
	}(&wg)

	sum := uint(0)
	for agentDonation := range donationChan {
		agentHp := globalState.AgentState[agentDonation.AgentID].Hp
		if agentDonation.Donation >= agentHp {
			agentDonation.Donation = agentHp
			delete(globalState.AgentState, agentDonation.AgentID)
			delete(agentMap, agentDonation.AgentID)
		}

		logging.Log(logging.Trace, logging.LogField{
			"Agent Donation": agentDonation,
			"Old Sum":        sum,
			"New Sum":        sum + agentDonation.Donation,
		}, "HP Pool Donation")

		sum += agentDonation.Donation
		if a, ok := globalState.AgentState[agentDonation.AgentID]; ok {
			a.Hp = agentHp - agentDonation.Donation
			globalState.AgentState[agentDonation.AgentID] = a
		}
	}

	logging.Log(logging.Info, logging.LogField{
		"Old HP Pool":           globalState.HpPool,
		"HP Donated This Round": sum,
		"New Hp Pool":           globalState.HpPool + sum,
	}, "HP Pool Donation")

	globalState.HpPool += sum
}
