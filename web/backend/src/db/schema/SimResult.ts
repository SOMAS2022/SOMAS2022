import mongoose from "mongoose";

const SimEntry = new mongoose.Schema({
    name: String,
    id: String,
    onGITCommit: String,
    time_queued: Date,
    time_taken: Number,
    sim_status: String,
    result: String,
    error: String,
    winner: String,
    config: {
        levels: Number,
        runs: Number,
        startingHP: Number,
        startingAT: Number,
        startingSH: Number,
        baseSTAM: Number,
        randomQty: Number,
        team1Qty: Number,
        team2Qty: Number,
        team3Qty: Number,
        team4Qty: Number,
        team5Qty: Number,
        team6Qty: Number,
    },
    logs: Array
});

const SimEntryModel = mongoose.model("SimEntry", SimEntry);

export default SimEntryModel;

