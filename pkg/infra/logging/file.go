package logging

import (
	"encoding/json"
	"fmt"
	"infra/game/commons"
	"os"
	"path"
)

// type LevelStage uint

// const (
// 	Election LevelStage = iota
// 	FightRounds
// 	Loot
// 	HPPool
// )

var fileLog GameLog
var runID string

type Outcome bool

const (
	Loss Outcome = false
	Win  Outcome = true
)

type Mode string

const (
	Default Mode = "default"
)

type VotingStrategy uint

const (
	SingleChoicePlurality VotingStrategy = 1
	BordaCount            VotingStrategy = 2
)

type GameLog struct {
	Warnings []LogField
	Errors   []LogField
	Outcome  Outcome
	Config   Config
	Levels   []LevelStages
}

type Config struct {
	Mode              Mode
	Levels            uint
	StartingHP        uint
	StartingAttack    uint
	StartingShield    uint
	BaseStamina       uint
	PassThreshold     float32
	VotingStrategy    VotingStrategy
	VotingPreferences uint
	AgentRandomQty    uint
	AgentTeam1Qty     uint
	AgentTeam2Qty     uint
	AgentTeam3Qty     uint
	AgentTeam4Qty     uint
	AgentTeam5Qty     uint
	AgentTeam6Qty     uint
}

type LevelStages struct {
	LevelStats    LevelStats
	ElectionStage ElectionStage
	VONCStage     VONCStage
	FightStage    FightStage
	LootStage     LootStage
	HPPoolStage   HPPoolStage
	AgentLogs     map[commons.ID]AgentLog
}

type AgentLog struct {
	Name       string
	ID         commons.ID
	Properties map[string]float32
}

type LevelStats struct {
	NumberOfAgents       uint
	SkippedThroughHpPool bool
	CurrentLevel         uint
	HPPool               uint
	MonsterHealth        uint
	MonsterAttack        uint
	LeaderBeforeElection commons.ID
	LeaderAfterElection  commons.ID
	AverageAgentHealth   uint
	AverageAgentStamina  uint
	AverageAgentShield   uint
	AverageAgentAttack   uint
}

type ElectionStage struct {
	Occurred  bool
	Winner    commons.ID
	Team      string
	Manifesto ManifestoLog
}

type ManifestoLog struct {
	FightImposition     bool
	LootImposition      bool
	TermLength          uint
	ThresholdPercentage uint
}

type VONCStage struct {
	Occurred  bool
	For       uint
	Against   uint
	Abstain   uint
	Threshold uint
}

type FightStage struct {
	Occurred bool
	Rounds   []FightLog
}

type FightLog struct {
	AttackingAgents []commons.ID
	ShieldingAgents []commons.ID
	CoweringAgents  []commons.ID
	AttackSum       uint
	ShieldSum       uint
	AgentsRemaining uint
}

type LootStage struct {
	Occurred bool
}

type HPPoolStage struct {
	Occurred         bool
	DonatedThisRound uint
	OldHPPool        uint
	NewHPPool        uint
}

func AgentLogToFile(fields LogField, msg string) {
	agentLog := AgentLog{}
	for k, v := range fields {
		if k == "agentName" {
			if name, ok := v.(string); ok {
				agentLog.Name = name
			}
		} else if k == "agentID" {
			if id, ok := v.(string); ok {
				agentLog.ID = id
			}
		} else {
			if num, ok := v.(float32); ok {
				agentLog.Properties[k] = num
			}
		}
	}
	fmt.Println(gameState.CurrentLeader)
	fmt.Println(fileLog.Levels)
	fileLog.Levels[gameState.CurrentLevel-1].AgentLogs[agentLog.ID] = agentLog
}

func LogToFile(lvl Level, fields LogField, msg string, level LevelStages) {
	switch lvl {
	case Error:
		fileLog.Errors = append(fileLog.Errors, CombineMessageToFields(fields, msg))
		Log(lvl, fields, msg)
		return
	case Warn:
		fileLog.Warnings = append(fileLog.Warnings, CombineMessageToFields(fields, msg))
		Log(lvl, fields, msg)
		return
	case Info:
		fileLog.Levels = append(fileLog.Levels, level)
		return
	}
}

func CombineMessageToFields(fields LogField, msg string) LogField {
	if fields == nil {
		fields = make(map[string]interface{})
		fields["msg"] = msg
	} else {
		fields["msg"] = msg
	}
	return fields
}

func OutputLog(outcome Outcome) {
	fileLog.Outcome = outcome
	jsonBuf, err := json.MarshalIndent(fileLog, "", "\t")
	if err != nil {
		log.Fatalf("Failed to Marshal gameStates: %v", err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	outputDir := path.Join(wd, "output/output.json")

	err = os.WriteFile(outputDir, jsonBuf, 0777)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
		return
	}
}
