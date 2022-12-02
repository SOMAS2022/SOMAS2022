import { Box, Select, Stack, TextField, Typography, InputLabel, MenuItem, SelectChangeEvent } from "@mui/material";
import LoadingButton from "@mui/lab/LoadingButton";
import Grid from "@mui/material/Unstable_Grid2/Grid2";
import { useState, FormEvent } from "react";
import { simulation_result, simulation_status } from "../types/global";
import { v4 as uuidv4 } from "uuid";
import SubmitModal from "../components/schedule/SubmitModal";
import { FileUpload } from "@mui/icons-material";

export default function Schedule() {
    const [showSubmitModal, setShowSubmitModal] = useState(false);
    const [loading, setLoading] = useState(false); 
    const [formInput, setFormInput] = useState<simulation_result>({
        name: new Date(Date.now()).toLocaleString(),
        id: uuidv4(),
        onGITCommit: "",
        time_queued: new Date(Date.now()),
        time_taken: 0,
        sim_status: simulation_status.In_Queue,
        result: null,
        error: null,
        winner: null,
        config: {
            levels: 60,
            runs: 1,
            startingHP: 1000,
            startingAT: 20,
            startingSH: 20,
            baseSTAM: 2000,
            randomQty: 100,
            team1Qty: 0,
            team2Qty: 0,
            team3Qty: 0,
            team4Qty: 0,
            team5Qty: 0,
            team6Qty: 0,
        },
        logs: [],
    });
    const [age, setAge] = useState(1);

    const handleChange = (event: SelectChangeEvent) => {
        setAge(Number(event.target.value.toString()));
    };

    const toggleModal = () => {
        setLoading(false);
        setShowSubmitModal(!showSubmitModal);
        window.location.reload();
    };

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        console.log("submitted", formInput);
        setShowSubmitModal(true);

        const res =  await fetch("http://localhost:9000/sendToQueue", {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify(formInput)
        }).then(res => {
            setShowSubmitModal(true);
            return res.text();
        }).catch(err => console.log(err));

        console.log(res);
    };

    return (
        <Stack 
            spacing={2}
            justifyContent="space-between"
            component="form"
            autoComplete="off"
            onSubmit={handleSubmit}
        >
            <SubmitModal status={202} text={"Accepted"} initialOpen={showSubmitModal} toggle={toggleModal}/>
            <Typography variant={"h4"} style={{marginBottom: "6px"}}>
                Configure a new run:
            </Typography>
            <Typography variant={"h6"} style={{ marginBottom: "6px" }}>
                Simulation Details:
            </Typography>
            <Grid sx={{display: "flex", flexDirection: "row", padding: "5px"}}>
                <TextField id="nameField" label="Name" required variant="outlined" style={{ margin: "5px" }} defaultValue={formInput.name} onChange={(e) => setFormInput({...formInput, name: e.target.value})}/>
                <TextField id="idField" label="ID" variant="outlined" disabled style={{ margin: "5px" }} defaultValue={formInput.id} />
                <TextField id="levelsField" label="Levels" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.config.levels} onChange={(e) => setFormInput({ ...formInput, config: {...formInput.config, levels: Number(e.target.value)} })} />
                <TextField id="runsField" label="Runs" variant="outlined" type="number" disabled style={{ margin: "5px" }} defaultValue={formInput.config.runs} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, runs: Number(e.target.value) } })} />
            </Grid>
            <Typography variant={"h6"} style={{ marginBottom: "6px" }}>
                Starting Values:
            </Typography>
            <Box sx={{ display: "flex", flexDirection: "row", padding: "5px" }}>
                <TextField id="startHPField" label="Starting HP" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.config.startingHP} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, startingHP: Number(e.target.value) } })} />
                <TextField id="startATField" label="Starting Attack" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.config.startingAT} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, startingAT: Number(e.target.value) } })} />
                <TextField id="startSHField" label="Starting Shield" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.config.startingSH} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, startingSH: Number(e.target.value) } })} />
                <TextField id="startSTAMField" label="Base Stamina" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.config.baseSTAM} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, baseSTAM: Number(e.target.value) } })} />
                <TextField id="thresholdPCT" label="Threshold Percentage" variant="outlined" type="number" style={{ margin: "5px" }} disabled defaultValue={0.6} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, baseSTAM: Number(e.target.value) } })} />
            </Box>
            <Typography variant={"h6"} style={{ marginBottom: "6px" }}>
                Governance
            </Typography>
            <Box sx={{}}>
                <InputLabel id="demo-simple-select-label">Voting Strategy</InputLabel>
                <Select
                    labelId="demo-simple-select-label"
                    id="demo-simple-select"
                    value={age.toString()}
                    label="Age"
                    onChange={handleChange}
                    style={{width: "100%"}}
                >
                    <MenuItem value={1}>Single Choice Plurality</MenuItem>
                    <MenuItem value={2}>Borda Count</MenuItem>
                </Select>
            </Box>
            <Typography variant={"h6"} style={{ marginBottom: "6px" }}>
                Agent Details:
            </Typography>
            <Box sx={{ display: "flex", flexDirection: "row", padding: "5px"}}>
                <TextField id="randomQty" label="Random Agent Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.config.randomQty} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, randomQty: Number(e.target.value) } })} />
                <TextField id="team1Qty" label="Team 1 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.config.team1Qty} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, team1Qty: Number(e.target.value) } })} />
                <TextField id="team2Qty" label="Team 2 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.config.team2Qty} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, team2Qty: Number(e.target.value) } })} />
                <TextField id="team3Qty" label="Team 3 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.config.team3Qty} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, team3Qty: Number(e.target.value) } })} />
                <TextField id="team4Qty" label="Team 4 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.config.team4Qty} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, team4Qty: Number(e.target.value) } })} />
                <TextField id="team5Qty" label="Team 5 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.config.team5Qty} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, team5Qty: Number(e.target.value) } })} />
                <TextField id="team6Qty" label="Team 6 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.config.team6Qty} onChange={(e) => setFormInput({ ...formInput, config: { ...formInput.config, team6Qty: Number(e.target.value) } })} />
            </Box>
            <LoadingButton
                color="success"
                onClick={() => setLoading(true)}
                loading={loading}
                loadingPosition="start"
                startIcon={<FileUpload />}
                variant="contained"
                type="submit"
            >
                Add To Queue
            </LoadingButton>
        </Stack>
    );
}