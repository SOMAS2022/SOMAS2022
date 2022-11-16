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
export interface simulation_config {
    random_agent_qty: number
}