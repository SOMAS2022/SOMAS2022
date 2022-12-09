import { useState, useEffect } from "react";
import { Box, Tab, Tabs } from "@mui/material";
import { Run } from "../../../common/types";
import ResultsOverview from "./results/Overview";
import Game from "./results/Game";
import Agents from "./results/Agents";

interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
}

const agentPropertyMap: Record<string, Record<string, Array<number>>> = {};
const teamPropertyMap: Record<string, Record<string, Array<Array<number>>>> = {};

function initAgentPropertyMap(run: Run) {
    run.Logs.Levels.forEach(lvl => {
        for (const id in lvl.AgentLogs) {
            const properties = lvl.AgentLogs[id].Properties;
            const recToMod = agentPropertyMap[id] ? agentPropertyMap[id] : {};
            for (const prop in properties) {
                const value = properties[prop];
                if (recToMod[prop]) {
                    recToMod[prop].push(value);
                } else {
                    recToMod[prop] = [value];
                }
            }
            agentPropertyMap[id] = recToMod;
        }
    });
}

function initTeamPropertyMap(run: Run) {
    run.Logs.Levels.forEach((lvl, idx) => {
        for (const id in lvl.AgentLogs) {
            const name = lvl.AgentLogs[id].Name;
            console.log(name);
            const properties = lvl.AgentLogs[id].Properties;
            const recToMod = teamPropertyMap[name] ? teamPropertyMap[name] : {};
            for (const prop in properties) {
                const value = properties[prop];
                if (recToMod[prop]) {
                    if (recToMod[prop][idx]) {
                        recToMod[prop][idx].push(value);
                    } else {
                        recToMod[prop][idx] = [value];
                    }
                } else {
                    recToMod[prop] = [];
                    recToMod[prop][idx] = [value];
                }
            }
            teamPropertyMap[name] = recToMod;
        }
    });
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
                    {children}
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
    run: Run
}

export default function Result ({run}: ResultProps) {
    const [tabValue, setTabValue] = useState(0);

    const handleChange = (_: React.SyntheticEvent, newValue: number) => {
        setTabValue(newValue);
    };

    useEffect(() => {
        initAgentPropertyMap(run);
        initTeamPropertyMap(run);
    }, []);

    return (
        <Box>
            <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
                <Tabs value={tabValue} onChange={handleChange} aria-label="basic tabs example">
                    <Tab label="Overview" {...a11yProps(0)} />
                    <Tab label="Game" {...a11yProps(1)} />
                    <Tab label="Agents" {...a11yProps(2)} />
                    <Tab label="Logs" {...a11yProps(3)} />
                </Tabs>
            </Box>
            <TabPanel value={tabValue} index={0}>
                <ResultsOverview run={run}/>
            </TabPanel>
            <TabPanel value={tabValue} index={1}>
                <Game run={run}/>
            </TabPanel>
            <TabPanel value={tabValue} index={2}>
                <Agents run={run} agentPropertyMap={agentPropertyMap} teamPropertyMap={teamPropertyMap}/>
            </TabPanel>
            <TabPanel value={tabValue} index={3}>
                <div><pre>{JSON.stringify(run, null, 2)}</pre></div>
            </TabPanel>
        </Box>
    );
}