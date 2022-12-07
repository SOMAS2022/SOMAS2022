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
	if w < len(sortedAgentIDs) {
		peerAgent := sortedAgentIDs[w]
		sci := currentAgent.socialCapital[peerAgent]
		sci[1] = 0.8
		currentAgent.socialCapital[peerAgent] = sci
	}

	return false
}

func connectAgents(agentMap map[commons.ID]agent.Agent) {
	agents = agentMap
	numAgents := len(agentMap)
	// Create a grid graph of all the agents. Most agents would
	// be connected to 4 other agents
	gridHeight := 10
	gridWidth := numAgents / gridHeight
	if numAgents%gridHeight != 0 {
		gridWidth++
	}
	gridWidth += 20
	grid := build.Grid(gridHeight, gridWidth)
	// Create a complete bipartite graph. Meaning that 10 agents
	// are each connected to another 40 agents
	// veryConnectedAgents := int(0.05 * float32(numAgents))
	// quiteConnectedAgents := int(0.2 * float32(numAgents))
	// g1 := build.Kmn(veryConnectedAgents, quiteConnectedAgents)
	// combined := grid.Union(g1)
	// Uncomment to print out graph:
	// fmt.Println(combined.String())
	combined := grid
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
