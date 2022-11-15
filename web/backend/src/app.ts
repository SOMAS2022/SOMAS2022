import express from "express";
const app = express();
const PORT = 9000;

app.get("/", (_, res) => {
    res.send("Hello SOMAS!");
});

app.listen(PORT, () => {
    console.log(`[server]: Server is running at https://localhost:${PORT}`);
    return;
});