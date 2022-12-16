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
	"math/rand"
	"sort"
	"time"
)

// This type will make it easier to extract from map, sort, and retrieve agent ID
type SocialCapInfo struct {
	ID  string
	Arr [4]float64
}

func GetSortedAgentSubset(selfID string, socialCapital map[string][4]float64) []SocialCapInfo {
	// select agents randomly to get a subset of agents
	if len(socialCapital) == 0 {
		return make([]SocialCapInfo, 0)
	}
	listSC := make([]SocialCapInfo, 0)
	for k, sc := range socialCapital {
		if k == selfID { // Exclude self
			continue
		}
		sci := SocialCapInfo{ID: k, Arr: sc}
		listSC = append(listSC, sci)
	}
	rand.Seed(time.Now().Unix())
	permutation := rand.Perm(len(listSC))

	numSubset := 10
	if len(listSC) < numSubset {
		numSubset = len(listSC)
	}
	sortedSC := make([]SocialCapInfo, numSubset)
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
