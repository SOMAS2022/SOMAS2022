import { useState } from "react";
import { Box, Tab, Tabs } from "@mui/material";
import { Run } from "../../../common/types";
import ResultsOverview from "./results/Overview";
import Game from "./results/Game";

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

    return (
        <Box>
            <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
                <Tabs value={tabValue} onChange={handleChange} aria-label="basic tabs example">
                    <Tab label="Overview" {...a11yProps(0)} />
                    <Tab label="Game" {...a11yProps(1)} />
                    {/* <Tab label="Agents" {...a11yProps(2)} /> */}
                </Tabs>
            </Box>
            <TabPanel value={tabValue} index={0}>
                <ResultsOverview run={run}/>
            </TabPanel>
            <TabPanel value={tabValue} index={1}>
                <Game run={run}/>
            </TabPanel>
            <TabPanel value={tabValue} index={2}>
                {/* <Agents run={run}/> */}
            </TabPanel>
        </Box>
    );
}