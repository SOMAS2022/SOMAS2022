import { Box, Typography } from "@mui/material";
import { simulation_result } from "../types/global";

interface ResultProps {
    simRes: simulation_result
}

export default function Result ({simRes}: ResultProps) {
    console.log(simRes);
    return (
        <Box>        
            <Typography variant={"h5"}>
                {simRes.name}
            </Typography>
            <Typography variant="subtitle1">
                ID: {simRes.id}
            </Typography>
            <Typography variant="subtitle1">
                Queued: {simRes.time_queued?.toLocaleString()} - Time Taken: {simRes.time_taken}
            </Typography>
            <Typography variant="subtitle1">
                Status: {simRes.sim_status}
            </Typography>
            <Typography variant="subtitle1">
                Result: {simRes.result ? simRes.result + " - " + simRes.winner : simRes.sim_status} 
            </Typography>
            <Typography variant="subtitle1">
                {simRes.onGITCommit}
            </Typography>
            <pre>{JSON.stringify(simRes, null, 2)}</pre>
            {/* {JSON.stringify(simRes)} */}
        </Box>
    );
}