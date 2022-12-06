package logging

import (
	"encoding/json"
	"infra/game/commons"
	"os"
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
}

type LevelStats struct {
	NumberOfAgents       uint
	SkippedThroughHpPool bool
	CurrentLevel         uint
	HPPool               uint
	MonsterHealth        uint
	MonsterAttack        uint
	CurrentLeader        commons.ID
}

type ElectionStage struct {
	Occured   bool
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
	Occured   bool
	For       uint
	Against   uint
	Abstain   uint
	Threshold uint
}

type FightStage struct {
	Occured bool
	Rounds  []FightLog
}

type FightLog struct {
	AttackingAgents []commons.ID
	ShieldingAgents []commons.ID
	CoweringAgents  []commons.ID
	AttackSum       uint
	ShieldSum       uint
}

type LootStage struct {
	Occured bool
}

type HPPoolStage struct {
	Occured          bool
	DonatedThisRound uint
	OldHPPool        uint
	NewHPPool        uint
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
		// fmt.Print(level)
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
	logJSON, _ := json.Marshal(fileLog)
	os.WriteFile("logs/"+runID+".json", logJSON, 0644)
}
