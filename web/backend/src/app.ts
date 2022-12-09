import express from "express";
import cors from "cors";
import mongoose from "mongoose";
import { exec } from "child_process";
import { GameLog, Run, TeamNames, TeamScore } from "../../common/types";
import { getGITCommitHash  } from "./common/utils";
import RunModel from "./db/schema/Run";
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
    const config = req.body as Run;

    const env_vars = process.env;
    env_vars["MODE"] = config.Config.Mode;
    env_vars["LEVELS"] = config.Config.Levels.toString();
    env_vars["STARTING_HP"] = config.Config.StartingHP.toString();
    env_vars["STARTING_ATTACK"] = config.Config.StartingAttack.toString();
    env_vars["STARTING_SHIELD"] = config.Config.StartingShield.toString();
    env_vars["BASE_STAMINA"] = config.Config.BaseStamina.toString();
    env_vars["THRESHOLD_PCT"] = config.Config.PassThreshold.toString();
    env_vars["VOTING_STRATEGY"] = config.Config.VotingStrategy.toString();
    env_vars["VOTING_PREFERENCES"] = config.Config.VotingPreferences.toString();
    env_vars["AGENT_RANDOM_QUANTITY"] = config.Config.AgentRandomQty.toString();
    env_vars["AGENT_TEAM1_QUANTITY"] = config.Config.AgentTeam1Qty.toString();
    env_vars["AGENT_TEAM2_QUANTITY"] = config.Config.AgentTeam2Qty.toString();
    env_vars["AGENT_TEAM3_QUANTITY"] = config.Config.AgentTeam3Qty.toString();
    env_vars["AGENT_TEAM4_QUANTITY"] = config.Config.AgentTeam4Qty.toString();
    env_vars["AGENT_TEAM5_QUANTITY"] = config.Config.AgentTeam5Qty.toString();
    env_vars["AGENT_TEAM6_QUANTITY"] = config.Config.AgentTeam6Qty.toString();
        
    const startTime = new Date(Date.now());
    exec(`cd ../../ && make ID=${config.Meta.Id} runWithID`, { env: env_vars }, (error, stdout, stderr) => {
        if (error) {
            console.error(`error: ${error.message}`);
            return;
        }

        if (stderr) {
            console.error(`stderr: ${stderr}`);
            return;
        }
        
        const endTime = new Date(Date.now());
        config.Meta.TimeTaken = endTime.getTime()-startTime.getTime();
        console.log(endTime.getTime() - startTime.getTime());
        config.Meta.OnGITCommit = getGITCommitHash("short");
        fs.writeFile(`../../logs/${config.Meta.Id}.log`, stdout, function (err) {
            if (err) {
                return console.log(err);
            }
            console.log("The file was saved!");
        }); 

        fs.readFile(`../../logs/${config.Meta.Id}.json`, "utf8", function (err, data) {
            if (err) throw err;
            const log = JSON.parse(data) as GameLog;
            
            config.Logs = log;
            config; 

            const newRun = new RunModel(config);
            newRun.save((err) => {
                if (err) {
                    console.log(err);
                    return;
                }
            });
        });
    });
    res.status(202).send("Accepted");
});

app.get("/fetchRuns", (_, res) => {
    RunModel
        .find()
        .limit(20)
        .sort({ "Meta.TimeSubmitted": -1 })
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
    await mongoose.connect("");
    console.log("[server]: Connected to MongoDB");
    return;
});