/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/
package internal

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"math/rand"
	"sort"
	"time"
)

// This type will make it easier to extract from map, sort, and retrieve agent ID
type SocialCapInfo struct {
	ID  string
	Arr [4]float64
}

func getBestDonation(selfID string, m message.TradeInfo) (uint, string) {
	bestDonation := uint(0)
	var bestDonationId string
	for negId, neg := range m.Negotiations {
		if neg.Agent2 == selfID {
			if offer, ok := neg.GetOffer(neg.Agent1); ok {
				if offer.ItemType == commons.Weapon && offer.Item.Value() > bestDonation {
					bestDonation = offer.Item.Value()
					bestDonationId = negId
				}
			}
		}
	}

	return bestDonation, bestDonationId
}

func ShouldAcceptOffer(BA agent.BaseAgent, m message.TradeInfo) (commons.TradeID, bool, int) {
	selfID := BA.ID()
	agentWeapons := BA.AgentState().Weapons
	// get currently used weapon
	equippedWeapon := BA.AgentState().WeaponInUse
	var currentWeaponAttack uint = 0
	it := agentWeapons.Iterator()
	for !it.Done() {
		_, w := it.Next()
		if w.Id() == equippedWeapon {
			currentWeaponAttack = w.Value()
			break
		}
	}

	bestDonation, bestDonationId := getBestDonation(selfID, m)

	if bestDonation != 0 && bestDonation > currentWeaponAttack {
		// accept offer
		return bestDonationId, true, 0
	}

	// check what the second best weapon held is
	var bestFreeWStats uint = 0
	var bestFreeWIdx int = -1
	it = agentWeapons.Iterator()
	for !it.Done() {
		i, w := it.Next()
		if w.Id() != equippedWeapon {
			if bestFreeWStats < w.Value() {
				bestFreeWStats = w.Value()
				bestFreeWIdx = i
			}
		}
	}

	return "", false, bestFreeWIdx
}

func GetSortedAgentSubset(selfID string, socialCapital map[string][4]float64) []SocialCapInfo {
	// select agents randomly to get a subset of agents
	listSC := make([]SocialCapInfo, 0, len(socialCapital)-1)
	for k, sc := range socialCapital {
		if k == selfID { // Exclude self
			continue
		}
		sci := SocialCapInfo{ID: k, Arr: sc}
		listSC = append(listSC, sci)
	}
	rand.Seed(time.Now().Unix())
	permutation := rand.Perm(len(socialCapital) - 1)

	numSubset := 20
	if len(listSC) < numSubset {
		numSubset = len(listSC)
	}
	sortedSC := make([]SocialCapInfo, 0, numSubset)
	for t, i := range permutation {
		if t == numSubset {
			break
		}

		sortedSC = append(sortedSC, listSC[i])
	}

	sort.Slice(sortedSC, func(i int, j int) bool {
		return (OverallPerception(sortedSC[i].Arr) > OverallPerception(sortedSC[j].Arr))
	})

	return sortedSC
}
