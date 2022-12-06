import { Box } from "@mui/material";
import { Run } from "../../../../common/types";
import { AgentDeathGraph } from "../../components/game/AgentDeathGraph";
import { HPPoolGraph } from "../../components/game/HPPoolGraph";

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

    const HPPool = run.Logs.Levels.map(round => {
        return round.HPPoolStage.OldHPPool;
    });
    const MonsterResilience = run.Logs.Levels.map(round => {
        return round.LevelStats.MonsterAttack;
    });
    return (
        <Box>
            <AgentDeathGraph agents={agentsAlivePerRound} threshold={run.Config.PassThreshold*run.Logs.Levels[0].LevelStats.NumberOfAgents}/>
            <HPPoolGraph pool={HPPool} monsterHP={MonsterResilience}/>
        </Box>
    );
}