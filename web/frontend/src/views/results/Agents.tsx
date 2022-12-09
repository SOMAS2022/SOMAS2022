import { useState } from "react"; 
import { Box, Select, InputLabel, MenuItem, SelectChangeEvent, Slider, Typography } from "@mui/material";
import { Run } from "../../../../common/types";
import { AgentPropertyGraph } from "../../components/game/AgentPropertyGraph";
import { HistogramGraph } from "../../components/game/HistogramGraph";

interface AgentsProps {
    run: Run
    agentPropertyMap: Record<string, Record<string, Array<number>>>
    teamPropertyMap: Record<string, Record<string, Array<Array<number>>>>
}

export default function Agents({ run, agentPropertyMap, teamPropertyMap }: AgentsProps) {
    const [bins, setBins] = useState<number | string | Array<number | string>>(1);
    const [level, setLevel] = useState<number | string | Array<number | string>>(1);
    const [agentID, setAgentID] = useState<string>("");
    const [teamSelect, setTeamSelect] = useState<string>("");
    const [selectedAgent, setSelectedAgent] = useState<string>("");
    const [propertySelect, setPropertySelect] = useState<string>("");

    const agents: Array<string> = [];
    // console.log(teamPropertyMap["RANDOM"]["bravery"][0], frequencies(teamPropertyMap["RANDOM"]["bravery"][20], 1));

    // console.log(teamSelect in teamPropertyMap && propertySelect in teamPropertyMap[teamSelect] ? teamPropertyMap[teamSelect][propertySelect].flat() : 0);
    // console.log(teamSelect in teamPropertyMap && propertySelect in teamPropertyMap[teamSelect] ? Math.min(...(teamPropertyMap[teamSelect][propertySelect].flat())) : 0);
    // console.log(teamSelect in teamPropertyMap && propertySelect in teamPropertyMap[teamSelect] ? Math.max(...(teamPropertyMap[teamSelect][propertySelect].flat())) : 0);

    for (const a in run.Logs.Levels[0].AgentLogs) {
        const name = run.Logs.Levels[0].AgentLogs[a].Name;
        agents.push(`${a}, ${name}`);
    }

    // console.log(agentPropertyMap);
    console.log(teamPropertyMap);

    const handleBinChange = (e: Event, newValue: number | number[]) => {
        e.preventDefault();
        setBins(newValue);
    };

    const handleLevelChange = (e: Event, newValue: number | number[]) => {
        e.preventDefault();
        setLevel(newValue);
    };
    
    const handleTeamSelect = (event: SelectChangeEvent) => {
        setTeamSelect(event.target.value as string);
    };

    const handlePropertySelect = (event: SelectChangeEvent) => [
        setPropertySelect(event.target.value as string)
    ];

    const handleIDChange = (event: SelectChangeEvent) => {
        const idValue = event.target.value as string;
        const id = idValue.split(",")[0];
        setAgentID(idValue);
        setSelectedAgent(id);
    };

    const binMarks = [
        {
            value: 1,
            label: 1,
        },
        {
            value: 2,
            label: 2,
        },
    ];
    const levelMarks = [
        {
            value: 1,
            label: 1,
        },
        {
            value: run.Logs.Levels.length-1,
            label: run.Logs.Levels.length-1,
        },
    ];

    return (
        <Box sx={{ minWidth: 120 }}>
            <InputLabel id="Teams">Teams</InputLabel>
            <Select
                labelId="TeamSelect"
                id="TeamSelect"
                value={teamSelect}
                label="Agents"
                onChange={handleTeamSelect}
            >
                {Object.entries(teamPropertyMap).map(([key]) => {
                    return <MenuItem key={key} value={key}>{key}</MenuItem>;
                })}
            </Select>
            { teamSelect === "" 
                ? "" 
                : <><InputLabel id="Property">Property</InputLabel>
                    <Select
                        labelId="PropertySelect"
                        id="PropertySelect"
                        value={propertySelect}
                        label="Property"
                        onChange={handlePropertySelect}
                    >
                        {Object.entries(teamPropertyMap[teamSelect]).map(([prop]) => {
                            return <MenuItem key={prop} value={prop}>{prop}</MenuItem>;
                        })}
                    </Select></>
            } 
            
            <Typography>
                Level
            </Typography>
            <Slider id="level" defaultValue={1} value={typeof level === "number" ? level : 0} onChange={handleLevelChange} aria-label="Default" valueLabelDisplay="on" style={{ marginTop: "32px" }} marks={levelMarks} min={levelMarks[0].value} max={levelMarks[1].value} /> 

            <Typography>
                Bin Size
            </Typography>
            <Slider id="bin" defaultValue={1} value={typeof bins === "number" ? bins : 0} onChange={handleBinChange} aria-label="Default" valueLabelDisplay="on" style={{ marginTop: "32px" }} marks={binMarks} min={binMarks[0].value} max={binMarks[1].value} />

            <HistogramGraph property={propertySelect} values={teamSelect in teamPropertyMap && propertySelect in teamPropertyMap[teamSelect] ? teamPropertyMap[teamSelect][propertySelect] : []} idx={typeof level === "number" ? level-1 : 0} binSize={typeof bins === "number" ? bins : 0}/>
            
            <InputLabel id="Agents" style={{marginTop:24}}>Agents</InputLabel>
            <Select
                labelId="Agents"
                id="Agents"
                value={agentID}
                label="Agents"
                onChange={handleIDChange}
            >
                {agents.map(a => {
                    return <MenuItem key={a} value={a}>{a}</MenuItem>;
                })}
            </Select>
            {agentID === "" ? "" : <AgentPropertyGraph properties={agentPropertyMap[selectedAgent]}/>}
        </Box>
    );
}