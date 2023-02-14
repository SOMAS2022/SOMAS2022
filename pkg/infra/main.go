package main

import (
	"flag"
	"fmt"
	cmdline "infra/cmdLine"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	gamemath "infra/game/math"
	"infra/game/message"
	"infra/game/stage/discussion"
	"infra/game/stage/fight"
	"infra/game/stage/hppool"
	"infra/game/stage/loot"
	"infra/game/stage/trade"
	"infra/game/stages"
	"infra/logging"
	"infra/teams/team3"
	"math"
	"time"

	"github.com/benbjohnson/immutable"
)

var InitAgentMap = map[commons.ID]func() agent.Strategy{
	// "RANDOM": example.NewRandomAgent,
	// "TEAM1":  team1.NewSocialAgent,
	// "TEAM2":  team2.NewAgent2,
	"TEAM3NEUT": team3.NewAgentThreeNeutral,
	"TEAM3PAS":  team3.NewAgentThreePassive,
	"TEAM3AGR":  team3.NewAgentThreeAggressive,
	// "TEAM5":  team5.NewAgent5,
	// "TEAM6":  team6.NewTeam6Agent,
	// "TEAM4":  team4.NewAgentFour,
}

func main() {
	// define flags
	time := time.Now()
	useJSONFormatter := flag.Bool("j", false, "Whether to output logs in JSON")
	debug := flag.Bool("d", false, "Whether to run in debug mode. If false, only logs with level info or above will be shown")
	id := flag.String("i", time.String(), "Provide an ID for a given run")
	fixedSanction := flag.Int("fSanc", 1, "Provide fixed sanction length")
	dynamicSanction := flag.Bool("dSanc", false, "Toggle dynamic sanctioning")
	graduatedSanction := flag.Bool("gSanc", false, "Toggle graduated sanctioning")
	verbose := flag.Bool("verbose", true, "Toggle logger")

	flag.Parse()

	if *dynamicSanction && *graduatedSanction {
		fmt.Println("Cannot have both dynamic and graduated sanctions.")
		return
	}

	cmdline.CmdLineInits.FixedSanctionDuration = *fixedSanction
	cmdline.CmdLineInits.DynamicSanctions = *dynamicSanction
	cmdline.CmdLineInits.GraduatedSanctions = *graduatedSanction

	logging.InitLogger(*verbose, *useJSONFormatter, *debug, *id, globalState)
	initGame()
	startGameLoop()
}

func startGameLoop() {
	var decisionMap map[commons.ID]decision.FightAction
	var channelsMap map[commons.ID]chan message.TaggedMessage
	var termLeft uint
	channelsMap = addCommsChannels()
	*viewPtr = globalState.ToView()

	w, csvFile := initCsvLogging()
	for globalState.CurrentLevel = 1; globalState.CurrentLevel < (gameConfig.NumLevels + 1); globalState.CurrentLevel++ {
		levelLog := logging.LevelStages{}
		// Election Stage
		_, alive := agentMap[globalState.CurrentLeader]
		var votes map[decision.Intent]uint
		leaderBeforeElection := globalState.CurrentLeader
		if termLeft == 0 || !alive {
			termLeft = runElection()
			// fmt.Println(globalState.LeaderManifesto)
			levelLog.ElectionStage = logging.ElectionStage{
				Occurred: true,
				Winner:   globalState.CurrentLeader,
				Team:     agentMap[globalState.CurrentLeader].BaseAgent.Name(),
				Manifesto: logging.ManifestoLog{
					FightImposition:     globalState.LeaderManifesto.FightDecisionPower(),
					LootImposition:      globalState.LeaderManifesto.LootDecisionPower(),
					TermLength:          globalState.LeaderManifesto.TermLength(),
					ThresholdPercentage: globalState.LeaderManifesto.OverthrowThreshold(),
				},
			}
			logging.Log(logging.Info, logging.LogField{
				"Fight Imp": globalState.LeaderManifesto.FightDecisionPower(),
				"Loot Imp":  globalState.LeaderManifesto.LootDecisionPower(),
				"Term":      globalState.LeaderManifesto.TermLength(),
				"Threshold": globalState.LeaderManifesto.OverthrowThreshold(),
				"Winner":    globalState.CurrentLeader,
				"Team":      agentMap[globalState.CurrentLeader].BaseAgent.Name(),
			}, "Re-Election Vote")
		} else {
			levelLog.VONCStage = logging.VONCStage{Occurred: true, Threshold: globalState.LeaderManifesto.OverthrowThreshold()}
			termLeft, votes = runConfidenceVote(termLeft)
			levelLog.VONCStage.For = votes[decision.Positive]
			levelLog.VONCStage.Against = votes[decision.Negative]
			levelLog.VONCStage.Abstain = votes[decision.Abstain]
		}

		levelLog.LevelStats.SkippedThroughHpPool = checkHpPool()

		// allow agents to change the weapon and the shield in use
		globalState = loot.UpdateItems(*globalState, agentMap)
		*viewPtr = globalState.ToView()

		// Battle Rounds
		// TODO: Ambiguity in specification - do agents have a upper limit of rounds to try and slay the monster?
		fightResultSlice := make([]decision.ImmutableFightResult, 0)
		roundNum := uint(0)
		for globalState.MonsterHealth != 0 {
			levelLog.FightStage.Occurred = true
			// find out the maximum attack from alive agents
			maxAttack := uint(0)
			for _, agentState := range globalState.AgentState {
				if agentState.Hp > 0 {
					maxAttack += agentState.TotalAttack()
				}
			}
			// calculate average stats
			avgHP, avgAT, avgSH, avgST := uint(0), uint(0), uint(0), uint(0)
			for _, a := range agentMap {
				state := a.AgentState()
				avgHP += a.AgentState().Hp
				avgAT += state.TotalAttack()
				avgSH += state.TotalDefense()
				avgST += a.AgentState().Stamina
				// avgPers += a.personality
			}
			agents := uint(len(agentMap))
			avgHP, avgAT, avgSH, avgST = avgHP/agents, avgAT/agents, avgSH/agents, avgST/agents

			levelLog.LevelStats = logging.LevelStats{
				NumberOfAgents:       uint(len(agentMap)),
				CurrentLevel:         globalState.CurrentLevel,
				LeaderBeforeElection: leaderBeforeElection,
				LeaderAfterElection:  globalState.CurrentLeader,
				HPPool:               globalState.HpPool,
				MonsterHealth:        globalState.MonsterHealth,
				MonsterAttack:        globalState.MonsterAttack,
				AverageAgentHealth:   avgHP,
				AverageAgentAttack:   avgAT,
				AverageAgentShield:   avgSH,
				AverageAgentStamina:  avgST,
			}
			// end calc average stats

			decisionMapView := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
			for u, action := range decisionMap {
				decisionMapView.Set(u, action)
			}
			fightTally := stages.AgentFightDecisions(*globalState, agentMap, *decisionMapView.Map(), channelsMap)
			fightActions := discussion.ResolveFightDiscussion(*globalState, agentMap, agentMap[globalState.CurrentLeader], globalState.LeaderManifesto, fightTally)
			// Printing voted rules
			// prop := fightTally.GetMax()
			// rules := prop.Rules()
			// iterator := rules.Iterator()
			// for !iterator.Done() {
			// 	rule, _ := iterator.Next()
			// 	condition := rule.Condition()
			// 	associatedAction := ""
			// 	switch rule.Action() {
			// 	case 0:
			// 		associatedAction = "Defend"
			// 	case 1:
			// 		associatedAction = "Cower"
			// 	default:
			// 		associatedAction = "Attack"
			// 	}
			// 	fmt.Println(associatedAction, condition)
			// }
			globalState = fight.HandleFightRound(*globalState, gameConfig.Stamina, gameConfig.StartingHealthPoints, &fightActions)
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
				"AvHP":          avgHP,
				"AvST":          avgST,
				"AvAttack":      avgAT,
				"AvDefence":     avgSH,
			}, "Battle Summary")
			// NOTE: update the following function when you change AgentState
			damageCalculation(fightActions)
			levelLog.FightStage.Rounds = append(levelLog.FightStage.Rounds, logging.FightLog{
				AttackingAgents: fightActions.AttackingAgents,
				CoweringAgents:  fightActions.CoweringAgents,
				ShieldingAgents: fightActions.ShieldingAgents,
				AttackSum:       fightActions.AttackSum,
				ShieldSum:       fightActions.ShieldSum,
				AgentsRemaining: uint(len(agentMap)),
			})

			channelsMap = addCommsChannels()

			// log last round if we lose
			if float64(len(agentMap)) < math.Ceil(float64(gameConfig.ThresholdPercentage)*float64(gameConfig.InitialNumAgents)) {

				logLevel(levelLog, agentMap, w)

				logging.Log(logging.Info, nil, fmt.Sprintf("Lost on level %d  with %d remaining", globalState.CurrentLevel, len(agentMap)))
				logging.LogToFile(logging.Info, nil, "", levelLog)
				logging.OutputLog(logging.Loss)

				csvFile.Close()
				fmt.Println(csvFile.Name())
				return
			}
			fightResultSlice = append(fightResultSlice, *decision.NewImmutableFightResult(fightActions, roundNum))
			roundNum++
		}

		// TODO: Loot Discussion Stage

		lootPool := generateLootPool(uint(len(agentMap)))
		prunedAgentMap := stages.AgentPruneMapping(agentMap, globalState)
		// fmt.Println("PRUNED: ", prunedAgentMap)
		lootTally := stages.AgentLootDecisions(*globalState, *lootPool, agentMap, channelsMap)
		lootActions := discussion.ResolveLootDiscussion(*globalState, prunedAgentMap, lootPool, agentMap[globalState.CurrentLeader], globalState.LeaderManifesto, lootTally)
		globalState = loot.HandleLootAllocation(*globalState, lootActions, lootPool, prunedAgentMap)

		// fmt.Println(lootTally.ProposalTally())

		trade.HandleTrade(*globalState, agentMap, 5, 3)

		channelsMap = addCommsChannels()

		stages.HandleTrustStage(agentMap, channelsMap)

		channelsMap = addCommsChannels()

		levelLog.HPPoolStage = logging.HPPoolStage{Occurred: true, OldHPPool: globalState.HpPool}
		hppool.UpdateHpPool(agentMap, globalState)
		levelLog.HPPoolStage.NewHPPool = globalState.HpPool
		levelLog.HPPoolStage.DonatedThisRound = levelLog.HPPoolStage.NewHPPool - levelLog.HPPoolStage.OldHPPool

		// TODO: End of level Updates
		termLeft--
		globalState.MonsterHealth, globalState.MonsterAttack = gamemath.GetNextLevelMonsterValues(*gameConfig, globalState.CurrentLevel+1)
		*viewPtr = globalState.ToView()
		logging.Log(logging.Info, nil, fmt.Sprintf("------------------------------ Level %d Ended ----------------------------", globalState.CurrentLevel))

		immutableFightRounds := commons.NewImmutableList(fightResultSlice)
		votesResult := commons.MapToImmutable(votes)
		levelLog.AgentLogs = stages.UpdateInternalStates(agentMap, globalState, immutableFightRounds, &votesResult)

		logging.LogToFile(logging.Info, nil, "", levelLog)

		logLevel(levelLog, agentMap, w)
	}
	logging.Log(logging.Info, nil, fmt.Sprintf("Congratulations, The Peasants have escaped the pit with %d remaining.", len(agentMap)))
	logging.OutputLog(logging.Win)
	csvFile.Close()
	fmt.Println(csvFile.Name())
}
