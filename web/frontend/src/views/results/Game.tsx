import { Box } from "@mui/material";
import { Run } from "../../../../common/types";
import { AgentDeathGraph } from "../../components/game/AgentDeathGraph";

interface ResultsOverviewProps {
    run: Run
}

export default function ResultsOverview({ run }: ResultsOverviewProps) {
    const agentsAlivePerRound: Array<number> = [];
    for (let i = 0; i < run.Logs.Levels.length; i++) {
        if (!run.Logs.Levels[i].FightStage.Occurred) {
            agentsAlivePerRound.push(run.Logs.Levels[i].LevelStats.NumberOfAgents);
        } else {
            const lastFightRound = run.Logs.Levels[i].FightStage.Rounds[run.Logs.Levels[i].FightStage.Rounds.length-1];
            agentsAlivePerRound.push(lastFightRound.AgentsRemaining);
        }
    }
    return (
        <Box>
            <AgentDeathGraph agents={agentsAlivePerRound} threshold={run.Config.PassThreshold*run.Logs.Levels[0].LevelStats.NumberOfAgents}/>
        </Box>
    );
}