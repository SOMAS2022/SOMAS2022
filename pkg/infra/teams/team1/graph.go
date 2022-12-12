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
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
	"os"
	"sort"
	"strconv"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/yourbasic/graph/build"
)

var currentAgent *SocialAgent
var sortedAgentIDs []string
var agents map[commons.ID]agent.Agent

var gridWidth int

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
	gridWidth = numAgents / gridHeight
	if numAgents%gridHeight != 0 {
		gridWidth++
	}
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

	os.RemoveAll("./pkg/infra/teams/team1/graph/pics/")
	os.MkdirAll("./pkg/infra/teams/team1/graph/pics/", os.ModePerm)
}

func printGraph(agentMap map[commons.ID]agent.Agent, state *state.State) {
	numInitialAgents := len(sortedAgentIDs)
	aliveAgents := map[int]bool{}
	for i := 0; i < numInitialAgents; i++ {
		aliveAgents[i] = false
	}

	g := graph.New(graph.IntHash, graph.Directed())
	g.AddVertex(-1, graph.VertexAttribute("pos", "0,-1!"), graph.VertexAttribute("shape", "box"), graph.VertexAttribute("label", "Level: "+strconv.Itoa(int(state.CurrentLevel))))
	// add alive agents to graph
	for k, a := range agentMap {
		healthG := int(float64(state.AgentState[k].Hp) / 1000 * 255)
		if healthG > 255 {
			healthG = 255
		}
		healthR := 255 - healthG
		healthStr := "#" + fmt.Sprintf("%02X", healthR) + fmt.Sprintf("%02X", healthG) + "00"
		sa := a.Strategy.(*SocialAgent)
		pos := strconv.Itoa(sa.graphID/gridWidth) + "," + strconv.Itoa(sa.graphID%gridWidth) + "!"
		if k == state.CurrentLeader {
			g.AddVertex(sa.graphID, graph.VertexAttribute("style", "filled"), graph.VertexAttribute("fillcolor", "yellow"), graph.VertexAttribute("pos", pos), graph.VertexAttribute("color", healthStr))
		} else {
			g.AddVertex(sa.graphID, graph.VertexAttribute("pos", pos), graph.VertexAttribute("color", healthStr))
		}
		aliveAgents[sa.graphID] = true
	}
	// and dead agents
	for n, alive := range aliveAgents {
		if !alive {
			pos := strconv.Itoa(n/gridWidth) + "," + strconv.Itoa(n%gridWidth) + "!"
			g.AddVertex(n, graph.VertexAttribute("style", "filled"), graph.VertexAttribute("fillcolor", "red"), graph.VertexAttribute("pos", pos), graph.VertexAttribute("color", "#FF0000"))
		}
	}

	// add edges
	for _, a := range agentMap {
		sa := a.Strategy.(*SocialAgent)
		addEdge(g, sa.socialCapital, agentMap, sa.graphID)
	}

	filename := "./pkg/infra/teams/team1/graph/pics/graph" + strconv.Itoa(int(state.CurrentLevel)) + ".gv"

	file, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	draw.DOT(g, file)
	file.Close()
}

func addEdge(g graph.Graph[int, int], sc map[string][4]float64, agentMap map[commons.ID]agent.Agent, id int) {
	networkThreshold := 0.5
	for peer, sc := range sc {
		if sc[1] >= networkThreshold {
			if p, ok := agentMap[peer]; ok { // if alive
				peerSa := p.Strategy.(*SocialAgent)
				g.AddEdge(id, peerSa.graphID)
			}
		}
	}
}
