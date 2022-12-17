import { useState, useContext, useEffect } from "react";
import { styled, useTheme } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
import CssBaseline from "@mui/material/CssBaseline";
import MuiAppBar, { AppBarProps as MuiAppBarProps } from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import List from "@mui/material/List";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import MenuIcon from "@mui/icons-material/Menu";
import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Brightness4Icon from "@mui/icons-material/Brightness4";
import Brightness7Icon from "@mui/icons-material/Brightness7";

import { ColourModeContext } from "../App";
import { Book, NoteAdd, SignalCellular4Bar, WifiOff } from "@mui/icons-material";
import { Tooltip } from "@mui/material";
import Schedule from "./Schedule";
import Result from "./Result";
import { Run } from "../../../common/types";

const drawerWidth = 240;

const Main = styled("main", { shouldForwardProp: (prop) => prop !== "open" })<{
    open?: boolean;
}>(({ theme, open }) => ({
    flexGrow: 1,
    padding: theme.spacing(3),
    transition: theme.transitions.create("margin", {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
    }),
    marginLeft: `-${drawerWidth}px`,
    ...(open && {
        transition: theme.transitions.create("margin", {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen,
        }),
        marginLeft: 0,
    }),
}));

interface AppBarProps extends MuiAppBarProps {
    open?: boolean;
}

const AppBar = styled(MuiAppBar, {
    shouldForwardProp: (prop) => prop !== "open",
})<AppBarProps>(({ theme, open }) => ({
    transition: theme.transitions.create(["margin", "width"], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
    }),
    ...(open && {
        width: `calc(100% - ${drawerWidth}px)`,
        marginLeft: `${drawerWidth}px`,
        transition: theme.transitions.create(["margin", "width"], {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen,
        }),
    }),
}));

const DrawerHeader = styled("div")(({ theme }) => ({
    display: "flex",
    alignItems: "center",
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: "flex-end",
}));

interface LayoutProps {
    connected: boolean,
    latestGitCommit: string
}

export default function Layout({ connected, latestGitCommit }: LayoutProps) {
    const theme = useTheme();
    const colourMode = useContext(ColourModeContext);
    const [open, setOpen] = useState<boolean>(true);
    const [view, setView] = useState<number>(0);
    const [selectedRun, setSelectedRun] = useState<Run | null>(null);
    const [runs, setRuns] = useState<Run[]>([]);
    
    useEffect(() => {
        async function fetchResults() {
            const res = await fetch("http://localhost:9000/fetchRuns");
            const data = await res.json();
            console.log(data);
            setRuns(data);
        }

        fetchResults();
    }, []);
        
    const handleDrawerOpen = () => {
        setOpen(true);
    };

    const handleDrawerClose = () => {
        setOpen(false);
    };

    return (
        <Box sx={{ display: "flex" }}>
            <CssBaseline />
            <AppBar position="fixed" open={open}>
                <Toolbar>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        onClick={handleDrawerOpen}
                        edge="start"
                        sx={{ mr: 2, ...(open && { display: "none" }) }}
                    >
                        <MenuIcon />
                    </IconButton>
                    <Box sx={{ width: "100%" }}>
                        <Typography sx={{ml: 1, float:"left"}} variant="h6" noWrap component="div">
                        SOMAS2022 - Escape from the dark Pit(t)
                        </Typography>
                        <Tooltip title={theme.palette.mode === "dark" ? "Switch to light" : "Switch to dark"}>
                            <IconButton sx={{ ml: 1, float: "right" }} onClick={colourMode.toggleColourMode} color="inherit">
                                {theme.palette.mode === "dark" ? <Brightness7Icon /> : <Brightness4Icon />}
                            </IconButton>
                        </Tooltip>
                        <Tooltip title={connected ? "Connected to server" : "Not connected to server"}>
                            <IconButton sx={connected ? { ml: 1, float: "right", color: theme.palette.success.main } : { ml: 1, float: "right", color: theme.palette.error.main }} >
                                {connected ? <SignalCellular4Bar /> : <WifiOff />}
                            </IconButton>
                        </Tooltip>
                        <Tooltip title={"Git Commit"}>
                            <Typography variant="subtitle1" sx={{ ml: 1, float: "right", marginTop: "8px" }}>
                                Commit: {latestGitCommit}
                            </Typography>
                        </Tooltip>
                    </Box>
                </Toolbar>
            </AppBar>
            <Drawer
                sx={{
                    width: drawerWidth,
                    flexShrink: 0,
                    "& .MuiDrawer-paper": {
                        width: drawerWidth,
                        boxSizing: "border-box",
                    },
                }}
                variant="persistent"
                anchor="left"
                open={open}
            >
                <DrawerHeader>
                    <Typography style={{ width: "100%", textAlign: "right", marginRight: "5%" }}>
                        Runs
                    </Typography>
                    <IconButton onClick={handleDrawerClose}>
                        {theme.direction === "ltr" ? <ChevronLeftIcon /> : <ChevronRightIcon />}
                    </IconButton>
                </DrawerHeader>
                <Divider />
                <List>
                    <ListItem key={"create"} disablePadding onClick={() => setView(0)}>
                        <ListItemButton>
                            <ListItemIcon>
                                <NoteAdd />
                            </ListItemIcon>
                            <ListItemText primary={"Schedule a new run"} />
                        </ListItemButton>
                    </ListItem>
                </List>
                <Divider />
                <List>
                    {
                        runs.map((run) => {
                            const iconColour = theme.palette.success.main;

                            return (
                                <ListItem key={run.Meta.Id} disablePadding onClick={() => { setView(1); setSelectedRun(run);}}>
                                    <ListItemButton>
                                        <ListItemIcon>
                                            <Book style={{color: iconColour}} />
                                        </ListItemIcon>
                                        <ListItemText primary={run.Meta.Name} />
                                    </ListItemButton>
                                </ListItem>
                            );
                        })
                    }
                </List>
            </Drawer>
            <Main open={open}>
                <DrawerHeader />
                {
                    view === 0 ?
                        <Schedule />
                        :
                        selectedRun == null ?
                            <p>Loading...</p> :
                            <Result run={selectedRun} />
                }
            </Main>
        </Box>
    );
}