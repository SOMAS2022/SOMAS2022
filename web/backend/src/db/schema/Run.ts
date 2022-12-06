import mongoose from "mongoose";

const Run = new mongoose.Schema({
    Meta: Object,
    Logs: Object
});

const RunModel = mongoose.model("Run", Run);

export default RunModel;

