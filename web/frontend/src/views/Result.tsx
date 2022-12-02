/* eslint-disable @typescript-eslint/no-explicit-any */
import { Box, Tab, Tabs, Typography } from "@mui/material";
import { useState } from "react";
import DeathChart from "../components/result/DeathChart";
import { simulation_result } from "../types/global";

interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
}

function TabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ p: 3 }}>
                    <Typography>{children}</Typography>
                </Box>
            )}
        </div>
    );
}

function a11yProps(index: number) {
    return {
        id: `simple-tab-${index}`,
        "aria-controls": `simple-tabpanel-${index}`,
    };
}

interface ResultProps {
    simRes: any
}

export default function Result ({simRes}: ResultProps) {
    console.log(simRes);
    const [value, setValue] = useState(0);

    const handleChange = (_: React.SyntheticEvent, newValue: number) => {
        setValue(newValue);
    };

    const agentsAlivePerLevel: Array<number> = new Array(100);
    agentsAlivePerLevel.fill(1000, 0, 100);
    console.log(agentsAlivePerLevel);
    console.log(simRes.logs);
    for (let i = 0; i < simRes.logs.length; i++) {
        console.log("yo");
        if (simRes.logs[i]?.currLevel) {
            if (simRes.logs[i].numAgents < agentsAlivePerLevel[simRes.logs[i].currLevel]) {
                agentsAlivePerLevel[simRes.logs[i].currLevel] = simRes.logs[i].numAgents;
            }
        }
        if(simRes.logs[i]?.msg && simRes.logs[i].msg.startsWith("Lost on level")) {
            const tokens: any = simRes.logs[i].msg.split(" ");
            const level: number = tokens[3];
            const agents: number = tokens[6];
            console.log(level, agents);
            agentsAlivePerLevel[level] = [agents];
        }
    }
    console.log(agentsAlivePerLevel.filter(x => x < 1000));

    return (
        <Box>        
            <Typography variant={"h5"}>
                {simRes.name}
            </Typography>
            <Typography variant="subtitle1">
                Result: {simRes.result ? simRes.result + (simRes.winner ? " - " + simRes.winner : "") : simRes.sim_status} 
            </Typography>
            <Typography variant="subtitle1">
                Time Taken: {simRes.time_taken}ms
            </Typography>
            <Typography variant="subtitle1">
                Status: {simRes.sim_status}
            </Typography>
            <Typography variant="subtitle1">
                Commit: {simRes.onGITCommit}
            </Typography>
            {/* <pre>{JSON.stringify(simRes, null, 2)}</pre>
            {/* {JSON.stringify(simRes)} */}

            <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
                <Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
                    <Tab label="Game Config" {...a11yProps(0)} />
                    <Tab label="Overview" {...a11yProps(1)} />
                    <Tab label="In detail" {...a11yProps(2)} />
                </Tabs>
            </Box>
            <TabPanel value={value} index={0}>
                <pre>{JSON.stringify(simRes.config, null, 2)}</pre>
            </TabPanel>
            <TabPanel value={value} index={1}>
                <DeathChart agentsLeft={agentsAlivePerLevel.filter(x => x < 1000)}/>
            </TabPanel>
            <TabPanel value={value} index={2}>
            </TabPanel>
        </Box>
    );
}