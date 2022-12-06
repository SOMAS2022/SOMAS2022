export enum TeamNames {
    team1 = "Team1",
    team2 = "Team2",
    team3 = "Team3",
    team4 = "Team4",
    team5 = "Team5",
    team6 = "Team6",
    randomAgent = "RandomAgent"
}

export interface GameLog {
    warnings: Array<Map<string, string|number>>
    errors: Array<Map<string, string|number>>
    Levels: Array<LevelStages>
}

export interface LevelStages {
    LevelStats:    LevelStats
	ElectionStage: ElectionStage
	VONCStage:     VONCStage
	FightStage:    FightStage
	LootStage:     LootStage
	HPPoolStage:   HPPoolStage
}

export interface LevelStats {
	NumberOfAgents:       number
	SkippedThroughHpPool: boolean
	CurrentLevel:         number
	HPPool:               number
	MonsterHealth:        number
	MonsterAttack:        number
	CurrentLeader:        string
}

export interface ElectionStage {
    Occured:   boolean
	Winner:    string
	Team:      string
	Manifesto: ManifestoLog
}

export interface ManifestoLog {
	FightImposition:     boolean
	LootImposition:      boolean
	TermLength:          number
	ThresholdPercentage: number
}

export interface VONCStage {
	Occured:   boolean
	For:       number
	Against:   number
	Abstain:   number
	Threshold: number
}

export interface FightStage {
	Occured: boolean
    Rounds: Array<FightLog>
} 

export interface FightLog {
    AttackingAgents: Array<string>
    ShieldingAgents: Array<string>
    CoweringAgents: Array<string>
	AttackSum:       number
	ShieldSum:       number
}

export interface LootStage {
	Occured: boolean
}

export interface HPPoolStage {
	Occured :         boolean
	DonatedThisRound: number
	OldHPPool:        number
	NewHPPool:        number
}