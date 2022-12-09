import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    ChartData,
} from "chart.js";
import { Line } from "react-chartjs-2";
import { v4 as uuid } from "uuid";

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

export const options = {
    responsive: true,
    plugins: {
        legend: {
            position: "top" as const,
        },
        title: {
            display: true,
            text: "Agent Logs",
        },
    },
};

const labels = [0];

const data: ChartData< "line", (number | null)[], number>= {
    labels: labels,
    datasets: [
        {
            label: "yo",
            data: [0],
            borderColor: "rgb(1,2,3)",
            backgroundColor: "rgb(1,2,3)"
        }
    ],
};

interface AgentPropertyGraphProps {
    properties: Record<string, Array<number>>
}

function random_rgba() {
    const o = Math.round, r = Math.random, s = 255;
    return "rgba(" + o(r() * s) + "," + o(r() * s) + "," + o(r() * s) + "," + r().toFixed(1) + ")";
}

export function AgentPropertyGraph({ properties }: AgentPropertyGraphProps) {
    console.log(properties);
    let xaxis: number[] = [];
    const dataset = [];
    for (const prop in properties) {
        xaxis = Array(properties[prop].length).fill(0);
        const colour = random_rgba();
        dataset.push({
            label: prop,
            data: properties[prop],
            borderColor: colour,
            backgroundColor: colour
        });
    }
    xaxis = xaxis.map((_, i) => {return i+1;});
    data.labels = xaxis;
    data.datasets = dataset;
    console.log(dataset);
    console.log(xaxis);

    return <Line key={uuid().toString()} style={{ minHeight: "50vh", maxHeight: "50vh" }} options={options} data={data} />;
}
