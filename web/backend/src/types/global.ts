export enum team_names {
    team1 = "Team1",
    team2 = "Team2",
    team3 = "Team3",
    team4 = "Team4",
    team5 = "Team5",
    team6 = "Team6",
    randomAgent = "RandomAgent"
}

export interface team_score {
    team_name: team_names,
    wins: number
}

export enum simulation_status {
    Running = "Running",
    In_Queue = "In Queue",
    Finished = "Finished"
}

export enum run_result {
    Win = "Win",
    Loss = "Loss",
    Error = "Error"
}

export enum run_errors {
    "Time Exceeded" = "Time Exceeded"
}

export interface simulation_config {
    levels: number,
    runs: number,
    startingHP: number,
    startingAT: number,
    startingSH: number,
    baseSTAM: number,
    randomQty: number,
    team1Qty: number,
    team2Qty: number,
    team3Qty: number,
    team4Qty: number,
    team5Qty: number,
    team6Qty: number,
}

export interface simulation_result {
    name: string,
    id: string,
    time_queued: Date,
    time_completed: Date | null,
    sim_status: simulation_status,
    result: run_result | null,
    error: run_errors | null,
    winner: team_names | null,
    config: simulation_config
}

export interface battle_summary {
    level: string,
    currentLevel: number, 
    monsterHealth: number,
    monsterDamage: number,
    numberOfAgents: number,
    attackSum: number,
    shieldSum: number,
    time: Date
}

export interface game_summary {
    level: string,
    msg: string,
    time: Date
}