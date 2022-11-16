import express from "express";
import cors from "cors";
import { exec } from "child_process";
import { team_names, team_score } from "./types/global";

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
    console.log(req.body);
    const env_vars = process.env;
    env_vars["LEVELS"] = "400";
    exec("cd ../../ && make run", { env: env_vars }, (error, stdout, stderr) => {
        if (error) {
            console.error(`error: ${error.message}`);
            return;
        }

        if (stderr) {
            console.error(`stderr: ${stderr}`);
            return;
        }

        console.log(`stdout:\n${stdout}`);
    });
    res.status(202).send("Accepted");
});

app.listen(PORT, () => {
    console.log(`[server]: Server is running at https://localhost:${PORT}`);
    return;
});