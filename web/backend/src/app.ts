import express from "express";
import cors from "cors";
import mongoose from "mongoose";
import { exec } from "child_process";
import { simulation_result, simulation_status, team_names, team_score } from "./types/global";
import { getGITCommitHash, stdoutToLogJSON } from "./common/utils";
import SimEntryModel from "./db/schema/SimResult";

const app = express();
const PORT = 9000;

app.use(cors());
app.use(express.json());
app.use(express.urlencoded());

app.get("/", (_, res) => {
    res.send("Hello SOMAS! " + getGITCommitHash("short"));
});

app.get("/test", (_, res) => {
    res.status(200).send(getGITCommitHash("short"));
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
    const startTime = new Date(Date.now());

    let hasError = false;
    let hasLost = false;

    exec("cd ../../ && make runWithJSON", { env: env_vars }, (error, stdout, stderr) => {
        if (error) {
            console.error(`error: ${error.message}`);
            hasError = true;
            return;
        }

        if (stderr) {
            console.error(`stderr: ${stderr}`);
            return;
        }
        
        const endTime = new Date(Date.now());
        const log = stdoutToLogJSON(stdout);
        if (log.at(-1)?.msg.startsWith("Lost")) {
            hasLost = true;
        }

        const newSimEntry = new SimEntryModel({
            name: config.name,
            id: config.id,
            onGITCommit: getGITCommitHash("short"),
            time_queued: config.time_queued,
            time_taken: endTime.getTime() - startTime.getTime(),
            sim_status: simulation_status.Finished,
            result: hasError? "Error" : hasLost? "Loss" : "Win",
            error: hasError ? "Crashed" : null,
            winner: config.winner,
            config: config.config,
            logs: log
        });
        newSimEntry.save((err: unknown) => {
            if (err) {
                console.log(err);
                res.status(500).send("Error saving to database");
            }
        }); 
        console.log(log);
        
    });
    res.status(202).send("Accepted");
});

app.get("/fetchSimResults", (_, res) => {
    SimEntryModel
        .find()
        .limit(20)
        .sort({ time_queued: -1 })
        .exec((err, docs) => {
            // console.log(err, docs);
            if (err) {
                return res.status(500).send("Error fetching from database");
            }
            return res.status(200).json(docs);
        });
});

app.listen(PORT, async () => {
    console.log(`[server]: Server is running at https://localhost:${PORT}`);
    await mongoose.connect("mongodb+srv://root:SOMAS2022@sim-logs.pk5ldje.mongodb.net/?retryWrites=true&w=majority");
    console.log("[server]: Connected to MongoDB");
    return;
});