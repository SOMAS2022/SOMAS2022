/* eslint-disable react/prop-types */
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
    scales: {
        yAxes: [
            {display: true, beginAtZero: true}, {display:false, beginAtZero:true}
        ]
    },
    plugins: {
        legend: {
            position: "top" as const,
        },
        title: {
            display: true,
            text: "Number of agents left per level",
        },
    },
};

interface ChartProps {
    agentsLeft: Array<number>
}

export default function Chart({ agentsLeft }: ChartProps) {
    const labels = agentsLeft.map((x, idx) => idx+1);
    // const labels = [1,2,3];
    const threshold = new Array(agentsLeft.length).fill(60, 0, agentsLeft.length); 
    const data = {
        labels,
        datasets: [
            {
                label: "Threshold",
                data: threshold,
                borderColor: "rgb(255, 99, 132)",
                backgroundColor: "rgba(255, 99, 132, 0.5)",
                yAxisID: "y"
            },
            {
                label: "Agents left",
                data: agentsLeft,
                borderColor: "rgb(53, 162, 235)",
                backgroundColor: "rgba(53, 162, 235, 0.5)",
                yAxisID: "y"
            },
        ],
    };
    return <Line options={options} data={data} />;
}