export enum TeamNames {
    team1 = "Team1",
    team2 = "Team2",
    team3 = "Team3",
    team4 = "Team4",
    team5 = "Team5",
    team6 = "Team6",
    randomAgent = "RandomAgent"
}

export interface TeamScore {
    team_name: TeamNames,
    wins: number
}

export interface Run {
    Meta: Meta
    Config: Config
    Logs: GameLog
}

export interface Meta {
    Name: string,
    Id: string,
    OnGITCommit: string,
    TimeSubmitted: Date,
    TimeTaken: number,
}

export interface GameLog {
    Warnings: Array<Warnings>
    Errors: Array<Errors>
    Levels: Array<LevelStages>
    Outcome: boolean
}

export interface Warnings {
    msg: string
    [key: string]: string
}

export interface Errors {
    msg: string,
    [key: string]: string
}

export enum VotingStrategy {
    SingleChoicePlurality = 1,
    BordaCount
}

export interface Config {
    Mode: "default"
    Levels: number
    StartingHP: number
    StartingAttack: number
    StartingShield: number
    BaseStamina: number
    PassThreshold: number
    VotingStrategy: number
    VotingPreferences: VotingStrategy
    AgentRandomQty: number
    AgentTeam1Qty: number
    AgentTeam2Qty: number
    AgentTeam3Qty: number
    AgentTeam4Qty: number
    AgentTeam5Qty: number
    AgentTeam6Qty: number
}

export interface LevelStages {
    LevelStats: LevelStats
    ElectionStage: ElectionStage
    VONCStage: VONCStage
    FightStage: FightStage
    LootStage: LootStage
    HPPoolStage: HPPoolStage
    AgentLogs: Record<string, AgentLog>
}

export interface AgentLog {
    Name: string,
    ID: string,
    Properties: Record<string, number>
}

export interface PropertyLog {
    PropertyName: string
    PropertyValues: Record<string, number>
}

export interface LevelStats {
    NumberOfAgents: number
    SkippedThroughHpPool: boolean
    CurrentLevel: number
    HPPool: number
    MonsterHealth: number
    MonsterAttack: number
    LeaderBeforeElection: string
    LeaderAfterElection: string
    AverageAgentHealth: number
    AverageAgentStamina: number
    AverageAgentShield: number
    AverageAgentAttack: number
}

export interface ElectionStage {
    Occurred: boolean
    Winner: string
    Team: string
    Manifesto: ManifestoLog
}

export interface ManifestoLog {
    FightImposition: boolean
    LootImposition: boolean
    TermLength: number
    ThresholdPercentage: number
}

export interface VONCStage {
    Occurred: boolean
    For: number
    Against: number
    Abstain: number
    Threshold: number
}

export interface FightStage {
    Occurred: boolean
    Rounds: Array<FightLog>
}

export interface FightLog {
    AttackingAgents: Array<string>
    ShieldingAgents: Array<string>
    CoweringAgents: Array<string>
    AttackSum: number
    ShieldSum: number
    AgentsRemaining: number
}

export interface LootStage {
    Occurred: boolean
}

export interface HPPoolStage {
    Occurred: boolean
    DonatedThisRound: number
    OldHPPool: number
    NewHPPool: number
}