import { Box, Typography } from "@mui/material";
import { simulation_result } from "../types/global";

interface ResultProps {
    simRes: simulation_result
}

export default function Result ({simRes}: ResultProps) {
    return (
        <Box>        
            <Typography variant={"h5"}>
                {simRes.name}
            </Typography>
            <Typography variant="subtitle1">
                ID: {simRes.id}
            </Typography>
            <Typography variant="subtitle1">
                Queued: {simRes.time_queued?.toLocaleString()} - Completed: {simRes.time_completed?.toLocaleString()}
            </Typography>
            <Typography variant="subtitle1">
                Status: {simRes.sim_status}
            </Typography>
            <Typography variant="subtitle1">
                Result: {simRes.result ? simRes.result + " - " + simRes.winner : simRes.sim_status} 
            </Typography>
        </Box>
    );
}