import {useState, useEffect} from "react";
import { Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import Leaderboard from "../components/summary/Leaderboard";
import { TeamScore } from "../../../common/types";

export default function Summary() {
    const [data, setData] = useState<TeamScore[] | null>(null);

    useEffect(() => {
        async function fetchLeaderboardData() {
            await fetch("http://localhost:9000/fetchLeaderboardData")
                .then(res => res.json())
                .then((res: TeamScore[]) => setData(res))
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