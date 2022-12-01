package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/example"
	gamemath "infra/game/math"
	"infra/game/message"
	"infra/game/stage/discussion"
	"infra/game/stage/fight"
	"infra/game/stage/loot"
	"infra/game/stages"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

var InitAgentMap = map[commons.ID]func() agent.Strategy{
	"RANDOM": example.NewRandomAgent,
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
		if termLeft == 0 || !alive {
			termLeft = runElection()
		} else {
			termLeft = runConfidenceVote(termLeft)
		}

		// allow agents to change the weapon and the shield in use
		updatedGlobalState := loot.UpdateItems(*globalState, agentMap)
		globalState = &updatedGlobalState

		// TODO: Fight Discussion Stage

		// Battle Rounds
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth != 0 {
			// find out the maximum attack from alive agents
			var maxAttack uint = 0
			for _, agentState := range globalState.AgentState {
				if agentState.Hp > 0 {
					maxAttack = maxAttack + agentState.TotalAttack(*globalState)
				}
			}

			decisionMapView := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
			for u, action := range decisionMap {
				decisionMapView.Set(u, action)
			}
			tally := stages.AgentFightDecisions(*globalState, agentMap, *decisionMapView.Map(), channelsMap)
			fightActions := discussion.ResolveFightDiscussion(agentMap, agentMap[globalState.CurrentLeader], globalState.LeaderManifesto, tally)
			stateAfterFight := fight.HandleFightRound(*globalState, gameConfig.StartingHealthPoints, &fightActions)
			globalState = &stateAfterFight
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
		}

		// TODO: Loot Discussion Stage
		weaponLoot, shieldLoot := make([]uint, len(agentMap)), make([]uint, len(agentMap))

		for i := range weaponLoot {
			weaponLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
			shieldLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
		}

		newGlobalState := stages.AgentLootDecisions(*globalState, agentMap, weaponLoot, shieldLoot)
		globalState = &newGlobalState

		// TODO: End of level Updates
		termLeft--
		globalState.MonsterHealth, globalState.MonsterAttack = gamemath.GetNextLevelMonsterValues(*gameConfig, globalState.CurrentLevel+1)
		*viewPtr = globalState.ToView()
		logging.Log(logging.Info, nil, fmt.Sprintf("------------------------------ Level %d Ended ----------------------------", globalState.CurrentLevel))
	}
	logging.Log(logging.Info, nil, fmt.Sprintf("Congratulations, The Peasants have escaped the pit with %d remaining.", len(agentMap)))
}
