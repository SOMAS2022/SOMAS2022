import { Box, TextField, Typography } from "@mui/material";
import LoadingButton from "@mui/lab/LoadingButton";
import Grid from "@mui/material/Unstable_Grid2/Grid2";
import { useState, FormEvent } from "react";
import { Run } from "../../../common/types";
import { v4 as uuidv4 } from "uuid";
import SubmitModal from "../components/schedule/SubmitModal";
import { FileUpload } from "@mui/icons-material";

export default function Schedule() {
    const [showSubmitModal, setShowSubmitModal] = useState(false);
    const [loading, setLoading] = useState(false); 
    const [formInput, setFormInput] = useState<Run>({
        Meta: {
            Name: new Date(Date.now()).toLocaleString(),
            Id: uuidv4(),
            OnGITCommit: "",
            TimeSubmitted: new Date(Date.now()),
            TimeTaken: 0
        },
        Config: {
            Mode: "default",
            Levels: 60,
            StartingHP: 1000,
            StartingAttack: 20,
            StartingShield: 20,
            BaseStamina: 2000,
            PassThreshold: 0.6,
            VotingStrategy: 1,
            VotingPreferences: 2,
            AgentRandomQty: 100,
            AgentTeam1Qty: 0,
            AgentTeam2Qty: 0,
            AgentTeam3Qty: 0,
            AgentTeam4Qty: 0,
            AgentTeam5Qty: 0,
            AgentTeam6Qty: 0
        },
        Logs: {
            Warnings: [],
            Errors: [],
            Levels: [],
            Outcome: false, 
        }
    });

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
        <Box 
            gap={12}
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
                <TextField id="nameField" label="Name" required variant="outlined" style={{ margin: "5px" }} defaultValue={formInput.Meta.Name} onChange={(e) => setFormInput({...formInput, Meta: {...formInput.Meta, Name: e.target.value}})}/>
                <TextField id="idField" label="ID" variant="outlined" disabled style={{ margin: "5px" }} defaultValue={formInput.Meta.Id} />
                <TextField id="levelsField" label="Levels" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.Config.Levels} onChange={(e) => setFormInput({ ...formInput, Config:{...formInput.Config, Levels: Number(e.target.value)}})} />
            </Grid>
            <Typography variant={"h6"} style={{ marginBottom: "6px" }}>
                Starting Values:
            </Typography>
            <Box sx={{ display: "flex", flexDirection: "row", padding: "5px" }}>
                <TextField id="startHPField" label="Starting HP" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.Config.StartingHP} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, StartingHP: Number(e.target.value) } })} />
                <TextField id="startATField" label="Starting Attack" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.Config.StartingAttack} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, StartingAttack: Number(e.target.value) } })} />
                <TextField id="startSHField" label="Starting Shield" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.Config.StartingShield} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, StartingShield: Number(e.target.value) } })} />
                <TextField id="startSTAMField" label="Base Stamina" variant="outlined" type="number" style={{ margin: "5px" }} defaultValue={formInput.Config.BaseStamina} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, BaseStamina: Number(e.target.value) } })} />
            </Box>
            <Typography variant={"h6"} style={{ marginBottom: "6px" }}>
                Agent Details:
            </Typography>
            <Box sx={{ display: "flex", flexDirection: "row", padding: "5px" }}>
                <TextField id="randomQty" label="Random Agent Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.Config.AgentRandomQty} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, AgentRandomQty: Number(e.target.value) } })} />
                <TextField id="team1Qty" label="Team 1 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.Config.AgentTeam1Qty} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, AgentTeam1Qty: Number(e.target.value) } })} />
                <TextField id="team2Qty" label="Team 2 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.Config.AgentTeam2Qty} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, AgentTeam2Qty: Number(e.target.value) } })} />
                <TextField id="team3Qty" label="Team 3 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.Config.AgentTeam3Qty} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, AgentTeam3Qty: Number(e.target.value) } })} />
                <TextField id="team4Qty" label="Team 4 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.Config.AgentTeam4Qty} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, AgentTeam4Qty: Number(e.target.value) } })} />
                <TextField id="team5Qty" label="Team 5 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.Config.AgentTeam4Qty} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, AgentTeam5Qty: Number(e.target.value) } })} />
                <TextField id="team6Qty" label="Team 6 Qty" variant="outlined" type={"number"} style={{ margin: "5px" }} defaultValue={formInput.Config.AgentTeam5Qty} onChange={(e) => setFormInput({ ...formInput, Config: { ...formInput.Config, AgentTeam6Qty: Number(e.target.value) } })} />
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
        </Box>
    );
}