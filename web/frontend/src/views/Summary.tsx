import {useState, useEffect} from "react";
import { Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import Leaderboard from "../components/summary/Leaderboard";
import { team_score } from "../types/global";

export default function Summary() {
    const [data, setData] = useState<team_score[] | null>(null);

    useEffect(() => {
        async function fetchLeaderboardData() {
            await fetch("http://localhost:9000/fetchLeaderboardData")
                .then(res => res.json())
                .then((res: team_score[]) => setData(res))
                .catch((err) => {console.log(err);});
        }
        fetchLeaderboardData();
    }, []);

    return (
        <Grid container spacing={4}>
            <Grid xs={4}>
                <Typography variant="h5" style={{textDecoration: "underline", textAlign: "center", marginBottom: "6px"}}>
                    Leaderboard
                </Typography>
                <Leaderboard data={data}/>
            </Grid>
            <Grid xs={4}>
            </Grid>
            <Grid xs={4}>
            </Grid>
        </Grid>
    );
}