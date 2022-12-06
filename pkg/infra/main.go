package main

import (
	"flag"
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/example"
	gamemath "infra/game/math"
	"infra/game/message"
	"infra/game/stage/discussion"
	"infra/game/stage/fight"
	"infra/game/stage/hppool"
	"infra/game/stage/loot"
	"infra/game/stages"
	"infra/logging"
	"infra/teams/team1"
	"math"

	"github.com/benbjohnson/immutable"
)

var InitAgentMap = map[commons.ID]func() agent.Strategy{
	"RANDOM": example.NewRandomAgent,
	"TEAM1":  team1.CreateSocialAgent,
}

func main() {
	// define flags
	useJSONFormatter := flag.Bool("j", false, "Whether to output logs in JSON")
	debug := flag.Bool("d", false, "Whether to run in debug mode. If false, only logs with level info or above will be shown")
	flag.Parse()

	logging.InitLogger(*useJSONFormatter, *debug)
	initGame()
	startGameLoop()
}

func startGameLoop() {
	var decisionMap map[commons.ID]decision.FightAction
	var channelsMap map[commons.ID]chan message.TaggedMessage
	var termLeft uint
	channelsMap = addCommsChannels()
	*viewPtr = globalState.ToView()

	for globalState.CurrentLevel = 1; globalState.CurrentLevel < (gameConfig.NumLevels + 1); globalState.CurrentLevel++ {
		// Election Stage
		_, alive := agentMap[globalState.CurrentLeader]
		var votes map[decision.Intent]uint
		if termLeft == 0 || !alive {
			termLeft = runElection()
		} else {
			termLeft, votes = runConfidenceVote(termLeft)
		}

		checkHpPool()

		// allow agents to change the weapon and the shield in use
		globalState = loot.UpdateItems(*globalState, agentMap)
		*viewPtr = globalState.ToView()

		// Battle Rounds
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		fightResultSlice := make([]decision.ImmutableFightResult, 0)
		roundNum := uint(0)
		for globalState.MonsterHealth != 0 {
			// find out the maximum attack from alive agents
			maxAttack := uint(0)
			for _, agentState := range globalState.AgentState {
				if agentState.Hp > 0 {
					maxAttack += agentState.TotalAttack(*globalState)
				}
			}

			decisionMapView := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
			for u, action := range decisionMap {
				decisionMapView.Set(u, action)
			}
			fightTally := stages.AgentFightDecisions(*globalState, agentMap, *decisionMapView.Map(), channelsMap)
			fightActions := discussion.ResolveFightDiscussion(*globalState, agentMap, agentMap[globalState.CurrentLeader], globalState.LeaderManifesto, fightTally)
			globalState = fight.HandleFightRound(*globalState, gameConfig.StartingHealthPoints, &fightActions)
			*viewPtr = globalState.ToView()

			logging.Log(logging.Info, logging.LogField{
				"currLevel":     globalState.CurrentLevel,
				"monsterHealth": globalState.MonsterHealth,
				"monsterDamage": globalState.MonsterAttack,
				"numCoward":     len(fightActions.CoweringAgents),
				"attackSum":     fightActions.AttackSum,
				"shieldSum":     fightActions.ShieldSum,
				"numAgents":     len(agentMap),
				"maxAttack":     maxAttack,
			}, "Battle Summary")

			// NOTE: update the following function when you change AgentState
			damageCalculation(fightActions)

			channelsMap = addCommsChannels()

			if float64(len(agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {
				logging.Log(logging.Info, nil, fmt.Sprintf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap)))
				return
			}
			fightResultSlice = append(fightResultSlice, *decision.NewImmutableFightResult(fightActions, roundNum))
			roundNum++
		}

		// TODO: Loot Discussion Stage

		lootPool := generateLootPool(len(agentMap), globalState.CurrentLevel)
		lootTally := stages.AgentLootDecisions(*globalState, *lootPool, agentMap, channelsMap)
		lootActions := discussion.ResolveLootDiscussion(*globalState, agentMap, lootPool, agentMap[globalState.CurrentLeader], globalState.LeaderManifesto, lootTally)
		globalState = loot.HandleLootAllocation(*globalState, &lootActions, lootPool)

		channelsMap = addCommsChannels()

		hppool.UpdateHpPool(agentMap, globalState)

		// TODO: End of level Updates
		termLeft--
		globalState.MonsterHealth, globalState.MonsterAttack = gamemath.GetNextLevelMonsterValues(*gameConfig, globalState.CurrentLevel+1)
		*viewPtr = globalState.ToView()
		logging.Log(logging.Info, nil, fmt.Sprintf("------------------------------ Level %d Ended ----------------------------", globalState.CurrentLevel))

		immutableFightRounds := commons.NewImmutableList(fightResultSlice)
		votesResult := commons.MapToImmutable(votes)
		updateInternalStates(immutableFightRounds, &votesResult)
	}
	logging.Log(logging.Info, nil, fmt.Sprintf("Congratulations, The Peasants have escaped the pit with %d remaining.", len(agentMap)))
}
