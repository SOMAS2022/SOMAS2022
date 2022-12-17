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
            text: "HP Pool",
        },
    },
};

const labels = [0];

const data = {
    labels: labels,
    datasets: [
        {
            label: "HP Pool at end of level",
            data: [0],
            borderColor: "rgb(255, 99, 132)",
            backgroundColor: "rgba(255, 99, 132, 0.5)",
        },
        {
            label: "Monster Resilience",
            data: [0],
            borderColor: "rgb(53, 162, 235)",
            backgroundColor: "rgba(53, 162, 235, 0.5)",
        },
    ],
};

interface HPPoolGraphProps {
    pool: Array<number>
    monsterHP: Array<number>
}

export function HPPoolGraph({ pool, monsterHP }: HPPoolGraphProps) {
    const xaxis = pool.map((_, idx) => {
        return idx + 1;
    });
    data.labels = xaxis;
    data.datasets[0].data = pool;
    data.datasets[1].data = monsterHP;

    return <Line key={uuid().toString()} options={options} data={data} style={{ minHeight: "50vh", maxHeight: "50vh" }} />;
}
