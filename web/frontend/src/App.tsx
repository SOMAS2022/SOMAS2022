import {useState, useMemo, createContext, useEffect} from "react";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import Layout from "./views/Layout";
import { useMediaQuery } from "@mui/material";

export const ColourModeContext = createContext({ toggleColourMode: () => {return;} });

export default function App() {
    const prefersDarkMode = useMediaQuery("(prefers-color-scheme: light)");
    const [mode, setMode] = useState<"light" | "dark">(prefersDarkMode ? "dark" : "light");
    const [conn, setConn] = useState<boolean>(false);
    const [latestGitCommit, setLatestGitCommit] = useState<string>("");
    const colorMode = useMemo(
        () => ({
            toggleColourMode: () => {
                setMode((prevMode) => (prevMode === "light" ? "dark" : "light"));
            },
        }),
        [],
    );

    useEffect(() => {
        async function establishServerConnection() {
            await fetch("http://localhost:9000/test")
                .then(res => {
                    console.info(res); 
                    setConn(true); 
                    return res.text();
                }).then(txt => setLatestGitCommit(txt)).catch((err) => {
                    console.error(err); 
                    setConn(false);
                });
        }
        establishServerConnection();
    }, []);

    const theme = useMemo(
        () =>
            createTheme({
                palette: {
                    mode,
                    primary: {
                        main: "#40467f"
                    },
                    success: {
                        main: "#b4bf58"
                    },
                    warning: {
                        main: "#884d85"
                    },
                    error: {
                        main: "#d94b38"
                    },
                },
            }),
        [mode],
    );

    return (
        <ColourModeContext.Provider value={colorMode}>
            <ThemeProvider theme={theme}>
                <Layout connected={conn} latestGitCommit={latestGitCommit}/>
            </ThemeProvider>
        </ColourModeContext.Provider>
    );
}
