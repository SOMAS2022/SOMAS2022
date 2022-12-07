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
            text: "Average agent ST per level",
        },
    },
};

const labels = [0];

const data = {
    labels: labels,
    datasets: [
        {
            label: "Average agent ST per level",
            data: [0],
            borderColor: "rgb(255, 99, 132)",
            backgroundColor: "rgba(255, 99, 132, 0.5)",
        },
    ],
};

interface AverageAgentSTGraphProps {
    points: Array<number>
}

export function AverageAgentSTGraph({ points }: AverageAgentSTGraphProps) {
    const xaxis = points.map((_, idx) => {
        return idx + 1;
    });
    data.labels = xaxis;
    data.datasets[0].data = points;

    return <Line key={uuid().toString()} options={options} data={data} />;
}
