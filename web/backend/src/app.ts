import express from "express";
import cors from "cors";
import { exec } from "child_process";
import { battle_summary, game_summary, simulation_result, team_names, team_score } from "./types/global";

const app = express();
const PORT = 9000;

app.use(cors());
app.use(express.json());
app.use(express.urlencoded());

app.get("/", (_, res) => {
    res.send("Hello SOMAS!");
});

app.get("/test", (_, res) => {
    res.status(200).send("Connected");
});

app.get("/fetchLeaderboardData", (_, res) => {
    const data: team_score[] = [
        { team_name: team_names.team1, wins: 50 },
        { team_name: team_names.team2, wins: 100 },
        { team_name: team_names.team3, wins: 60 },
        { team_name: team_names.team4, wins: 30 },
        { team_name: team_names.team5, wins: 10 },
        { team_name: team_names.team6, wins: 2 },
        { team_name: team_names.randomAgent, wins: 786 },
    ];
    res.status(200).json(data);
});

app.post("/sendToQueue", (req, res) => {
    const config = req.body as simulation_result;
    const env_vars = process.env;
    env_vars["LEVELS"] = config.config.levels.toString();
    exec("cd ../../ && make runWithJSON", { env: env_vars }, (error, stdout, stderr) => {
        if (error) {
            console.error(`error: ${error.message}`);
            return;
        }

        if (stderr) {
            console.error(`stderr: ${stderr}`);
            return;
        }
        const log: Array<battle_summary | game_summary> = [];

        const stdoutToLogJSON = (stdout: string) => {
            const array = stdout.split("\n");
            for (let i = 0; i < array.length; i++) {
                if (array[i][0] != "{") continue;
                try {
                    const json = JSON.parse(array[i]) as battle_summary | game_summary;
                    json.time = new Date(json.time);
                    console.log(json);
                    log.push(json);
                } catch {
                    console.log(`Invalid JSON: ${array[i]}`);
                    continue;
                }
            }
        };
        
        stdoutToLogJSON(stdout);
        console.log(log);
        
    });
    res.status(202).send("Accepted");
});

app.listen(PORT, () => {
    console.log(`[server]: Server is running at https://localhost:${PORT}`);
    return;
});