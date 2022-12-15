import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
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
            text: "Agents Remaining at the end of each level",
        },
    },
};

const labels = [0];

const data = {
    labels: labels,
    datasets: [
        {
            label: "Agents Remaining",
            data: [0],
            borderColor: "rgb(255, 99, 132)",
            backgroundColor: "rgba(255, 99, 132, 0.5)",
        },
        {
            label: "Pass Threshold",
            data: [0],
            borderColor: "rgb(53, 162, 235)",
            backgroundColor: "rgba(53, 162, 235, 0.5)",
        },
    ],
};

interface AgentDeathGraphProps {
    agents: Array<number>
    threshold: number
}

export function AgentDeathGraph({agents, threshold}: AgentDeathGraphProps) {
    const xaxis = agents.map((_, idx) => {
        return idx+1;
    });
    const thresholdArr = agents.map(() => {
        return threshold;
    });
    data.labels = xaxis;
    data.datasets[0].data = agents;
    data.datasets[1].data= thresholdArr;
    
    return <Line key={uuid().toString()} style={{minHeight: "50vh", maxHeight: "50vh"}} options={options} data={data} />;
}
