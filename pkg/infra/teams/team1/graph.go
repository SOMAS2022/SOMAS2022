/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/

package team1

import (
	"infra/game/agent"
	"infra/game/commons"
	"sort"

	"github.com/yourbasic/graph/build"
)

var currentAgent *SocialAgent
var sortedAgentIDs []string
var agents map[commons.ID]agent.Agent

func createConnection(w int, _ int64) bool {
	peerAgent := sortedAgentIDs[w]
	sci := currentAgent.socialCapital[peerAgent]
	sci.arr[1] = 0.81
	currentAgent.socialCapital[peerAgent] = sci

	return false
}

func connectAgents(agentMap map[commons.ID]agent.Agent) {
	agents = agentMap
	// Create a grid graph of all the agents. Most agents would
	// be connected to 4 other agents
	grid := build.Grid(10, 20)
	// Create a complete bipartite graph. Meaning that 10 agents
	// are each connected to another 40 agents
	g1 := build.Kmn(10, 40)
	combined := grid.Union(g1)
	// Uncomment to print out graph:
	// fmt.Println(combined.String())

	agentIDs := make([]string, 0, len(agentMap))
	for k := range agentMap {
		agentIDs = append(agentIDs, k)
	}
	sort.Strings(agentIDs)
	sortedAgentIDs = agentIDs
	for i, k := range agentIDs {
		currentAgent = agents[k].Strategy.(*SocialAgent)
		currentAgent.graphID = i
		combined.Visit(i, createConnection)
	}
}
