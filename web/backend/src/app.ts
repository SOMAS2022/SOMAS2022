import express from "express";
import cors from "cors";
import mongoose from "mongoose";
import { exec } from "child_process";
import { GameLog, RunDetails, TeamNames, TeamScore } from "../../common/types";
import { getGITCommitHash  } from "./common/utils";
import SimEntryModel from "./db/schema/Run";
import fs from "fs";

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
    const data: TeamScore[] = [
        { team_name: TeamNames.team1, wins: 50 },
        { team_name: TeamNames.team2, wins: 100 },
        { team_name: TeamNames.team3, wins: 60 },
        { team_name: TeamNames.team4, wins: 30 },
        { team_name: TeamNames.team5, wins: 10 },
        { team_name: TeamNames.team6, wins: 2 },
        { team_name: TeamNames.randomAgent, wins: 786 },
    ];
    res.status(200).json(data);
});

app.post("/sendToQueue", (req, res) => {
    const config = req.body as RunDetails;

    const env_vars = process.env;
    env_vars["MODE"] = config.Logs.Config.Mode;
    env_vars["LEVELS"] = config.Logs.Config.Levels.toString();
    env_vars["STARTING_HP"] = config.Logs.Config.Levels.toString();
    env_vars["STARTING_ATTACK"] = config.Logs.Config.Levels.toString();
    env_vars["STARTING_SHIELD"] = config.Logs.Config.Levels.toString();
    env_vars["BASE_STAMINA"] = config.Logs.Config.Levels.toString();
    env_vars["THRESHOLD_PCT"] = config.Logs.Config.Levels.toString();
    env_vars["VOTING_STRATEGY"] = config.Logs.Config.Levels.toString();
    env_vars["VOTING_PREFERENCES"] = config.Logs.Config.Levels.toString();
    env_vars["AGENT_RANDOM_QUANTITY"] = config.Logs.Config.Levels.toString();
    env_vars["AGENT_TEAM1_QUANTITY"] = config.Logs.Config.Levels.toString();
    env_vars["AGENT_TEAM2_QUANTITY"] = config.Logs.Config.Levels.toString();
    env_vars["AGENT_TEAM3_QUANTITY"] = config.Logs.Config.Levels.toString();
    env_vars["AGENT_TEAM4_QUANTITY"] = config.Logs.Config.Levels.toString();
    env_vars["AGENT_TEAM5_QUANTITY"] = config.Logs.Config.Levels.toString();
    env_vars["AGENT_TEAM6_QUANTITY"] = config.Logs.Config.Levels.toString();
        
    const startTime = new Date(Date.now());
    exec(`cd ../../ && make ID=${config.Meta.id} runWithID`, { env: env_vars }, (error, stdout, stderr) => {
        if (error) {
            console.error(`error: ${error.message}`);
            return;
        }

        if (stderr) {
            console.error(`stderr: ${stderr}`);
            return;
        }
        
        const endTime = new Date(Date.now());
        console.log(endTime.getTime()-startTime.getTime());
        fs.writeFile(`../../logs/${config.Meta.id}.log`, stdout, function (err) {
            if (err) {
                return console.log(err);
            }
            console.log("The file was saved!");
        }); 

        fs.readFile(`../../logs/${config.Meta.id}.json`, "utf8", function (err, data) {
            if (err) throw err;
            const log = JSON.parse(data) as GameLog;
            config.Logs = log;
            
        });

        // const newSimEntry = new SimEntryModel({
        //     name: config.name,
        //     id: config.id,
        //     onGITCommit: getGITCommitHash("short"),
        //     time_queued: config.time_queued,
        //     time_taken: endTime.getTime() - startTime.getTime(),
        //     sim_status: simulation_status.Finished,
        //     result: hasError? "Error" : hasLost? "Loss" : "Win",
        //     error: hasError ? "Crashed" : null,
        //     winner: config.winner,
        //     config: config.config,
        //     logs: log
        // });
        // newSimEntry.save((err: unknown) => {
        //     if (err) {
        //         console.log(err);
        //         res.status(500).send("Error saving to database");
        //     }
        // }); 
        // console.log(log);
        
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