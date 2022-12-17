import { Box } from "@mui/material";
import { Run } from "../../../../common/types";
import { AgentDeathGraph } from "../../components/game/AgentDeathGraph";
import { AgentDeathRoundGraph } from "../../components/game/AgentDeathRoundGraph";
import { AverageAgentStats } from "../../components/game/AverageAgentStats";
import { FightRoundGraph } from "../../components/game/FightRoundGraph";
import { HPPoolGraph } from "../../components/game/HPPoolGraph";
import { VONCGraph } from "../../components/game/VONCGraph";

interface ResultsOverviewProps {
    run: Run
}

export default function Game({ run }: ResultsOverviewProps) {
    const agentsAlivePerLevel: Array<number> = [];
    for (let i = 0; i < run.Logs.Levels.length; i++) {
        if (!run.Logs.Levels[i].FightStage.Occurred) {
            agentsAlivePerLevel.push(run.Logs.Levels[i].LevelStats.NumberOfAgents);
        } else {
            const lastFightRound = run.Logs.Levels[i].FightStage.Rounds[run.Logs.Levels[i].FightStage.Rounds.length-1];
            agentsAlivePerLevel.push(lastFightRound.AgentsRemaining);
        }
    }

    const agentsAlivePerRound: Array<number> = [];
    for (let i = 0; i < run.Logs.Levels.length; i++) {
        if (!run.Logs.Levels[i].FightStage.Occurred) {
            agentsAlivePerRound.push(run.Logs.Levels[i].LevelStats.NumberOfAgents);
            continue;
        }
        for(let j = 0; j < run.Logs.Levels[i].FightStage.Rounds.length; j++) {
            agentsAlivePerRound.push(run.Logs.Levels[i].FightStage.Rounds[j].AgentsRemaining);
        }
    }

    const AverageHP = run.Logs.Levels.map(level => {
        return level.LevelStats.AverageAgentHealth;
    });

    const AverageAT = run.Logs.Levels.map(level => {
        return level.LevelStats.AverageAgentAttack;
    });

    const AverageSH = run.Logs.Levels.map(level => {
        return level.LevelStats.AverageAgentShield;
    });

    const AverageST = run.Logs.Levels.map(level => {
        return level.LevelStats.AverageAgentStamina;
    });

    const HPPool = run.Logs.Levels.map(level => {
        return level.HPPoolStage.NewHPPool;
    });
    const MonsterResilience = run.Logs.Levels.map(level => {
        return level.LevelStats.MonsterAttack;
    });

    const FightRoundsPerLevel = run.Logs.Levels.map(level => {
        return level.FightStage.Rounds?.length;
    });

    const VONCForVotes = run.Logs.Levels.map((level, idx) => {
        if (level.VONCStage.Occurred) {
            return level.VONCStage.For;
        } else {
            console.log(idx);
            return 0;
        }
    });
    const VONCAgainstVotes = run.Logs.Levels.map(level => {
        if (level.VONCStage.Occurred) {
            return level.VONCStage.Against;
        } else {
            return 0;
        }
    });
    const VONCAbstainVotes = run.Logs.Levels.map(level => {
        if (level.VONCStage.Occurred) {
            return level.VONCStage.Abstain;
        } else {
            return 0;
        }
    });
    return (
        <Box>
            <AgentDeathGraph agents={agentsAlivePerLevel} threshold={run.Config.PassThreshold * run.Logs.Levels[0].LevelStats.NumberOfAgents} />
            <AgentDeathRoundGraph agents={agentsAlivePerRound} threshold={run.Config.PassThreshold * run.Logs.Levels[0].LevelStats.NumberOfAgents} />
            <AverageAgentStats HP={AverageHP} AT={AverageAT} SH={AverageSH} ST={AverageST}/>
            <HPPoolGraph pool={HPPool} monsterHP={MonsterResilience}/>
            <FightRoundGraph fightRounds={FightRoundsPerLevel}/>
            <VONCGraph forVotes={VONCForVotes} againstVotes={VONCAgainstVotes} abstainVotes={VONCAbstainVotes}/>
        </Box>
    );
}