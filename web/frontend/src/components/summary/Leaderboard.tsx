import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";
import { Typography } from "@mui/material";
import { team_score } from "../../types/global";

function compareScores(a: team_score, b: team_score) {
    if (a.wins < b.wins) {
        return 1;
    }
    if (a.wins > b.wins) {
        return -1;
    }
    return 0;
}

interface LeaderboardProps {
    data: team_score[] | null
}

export default function Leaderboard({data}:LeaderboardProps) {
    if (data === null) {
        return <Typography paragraph style={{textAlign: "center"}}>Loading...</Typography>;
    }
    return (
        <TableContainer component={Paper}>
            <Table aria-label="simple table">
                <TableHead>
                    <TableRow>
                        <TableCell>Team</TableCell>
                        <TableCell align="right">Wins</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {data.sort(compareScores).map((row) => (
                        <TableRow
                            key={row.team_name}
                            sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                        >
                            <TableCell component="th" scope="row">
                                {row.team_name}
                            </TableCell>
                            <TableCell align="right">{row.wins}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}