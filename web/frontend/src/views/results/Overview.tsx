import { Alert, Box, Typography } from "@mui/material";
import { Run } from "../../../../common/types";

interface ResultsOverviewProps {
    run: Run
}

export default function ResultsOverview({ run }: ResultsOverviewProps) {
    return (
        <Box>
            {
                run.Logs.Warnings.map(w => {
                    return <Alert key={w.msg} severity="warning">{w.msg}</Alert>;
                })
            }
            {
                run.Logs.Errors?.map(w => {
                    return <Alert key={w.msg} severity="error">{w.msg}</Alert>;
                })
            }
            <Typography mt={2}>
                Name: {run.Meta.Name}
            </Typography>
            <Typography>
                GIT Commit: {run.Meta.OnGITCommit}
            </Typography>
            <Typography>
                Time Submitted: {run.Meta.TimeSubmitted.toLocaleString()}
            </Typography>
            <Typography>
                Time Taken: {run.Meta.TimeTaken}ms
            </Typography>

            <Typography mt={2}>
                Outcome: {run.Logs.Outcome ? "Win" : "Loss"}
            </Typography>
        </Box>
    );
}