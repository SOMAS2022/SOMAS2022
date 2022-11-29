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
	"infra/game/stage/fight"
	"infra/game/stages"
	"infra/logging"
	"math"
	"math/rand"

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
		if termLeft == 0 {
			termLeft = runElection()
		} else {
			termLeft = runConfidenceVote(termLeft)
		}

		// TODO: Fight Discussion Stage

		// Battle Rounds
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		for globalState.MonsterHealth != 0 {
			decisionMapView := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
			for u, action := range decisionMap {
				decisionMapView.Set(u, action)
			}
			fightRoundResult := decision.FightResult{Choices: stages.AgentFightDecisions(*globalState, *agentMap, *decisionMapView.Map(), channelsMap)}
			stateAfterFight := fight.HandleFightRound(*globalState, gameConfig.StartingHealthPoints, &fightRoundResult)
			globalState = &stateAfterFight
			*viewPtr = globalState.ToView()

			logging.Log(logging.Info, logging.LogField{
				"currLevel":     globalState.CurrentLevel,
				"monsterHealth": globalState.MonsterHealth,
				"monsterDamage": globalState.MonsterAttack,
				"numCoward":     len(fightRoundResult.CoweringAgents),
				"attackSum":     fightRoundResult.AttackSum,
				"shieldSum":     fightRoundResult.ShieldSum,
				"numAgents":     len(*agentMap),
			}, "Battle Summary")

			damageCalculation(fightRoundResult)

			channelsMap = addCommsChannels()

			if float64(len(*agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {
				logging.Log(logging.Info, nil, fmt.Sprintf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(*agentMap)))
				return
			}
		}

		// TODO: Loot Discussion Stage
		weaponLoot, shieldLoot := make([]uint, len(*agentMap)), make([]uint, len(*agentMap))

		for i := range weaponLoot {
			weaponLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
			shieldLoot[i] = globalState.CurrentLevel * uint(rand.Intn(3))
		}

		newGlobalState := stages.AgentLootDecisions(*globalState, *agentMap, weaponLoot, shieldLoot)
		globalState = &newGlobalState

		// TODO: End of level Updates
		termLeft--
		globalState.MonsterHealth, globalState.MonsterAttack = gamemath.GetNextLevelMonsterValues(*gameConfig, globalState.CurrentLevel+1)
		*viewPtr = globalState.ToView()
		logging.Log(logging.Info, nil, fmt.Sprintf("------------------------------ Level %d Ended ----------------------------", globalState.CurrentLevel))
	}
	logging.Log(logging.Info, nil, fmt.Sprintf("Congratulations, The Peasants have escaped the pit with %d remaining.", len(*agentMap)))
}
