import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
} from "chart.js";
import { Bar } from "react-chartjs-2";
import { v4 as uuid } from "uuid";

ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend
);

export const options = {
    responsive: true,
    scales: {
        y: {
            ticks: {
                stepSize: 1
            },
        }
    },
    plugins: {
        legend: {
            position: "top" as const,
        },
        title: {
            display: true,
            text: "Fight Rounds per Level",
        },
    },
};

const labels = [0];

const data = {
    labels: labels,
    datasets: [
        {
            label: "Fight Rounds",
            data: [0],
            borderColor: "rgb(255, 99, 132)",
            backgroundColor: "rgba(255, 99, 132, 0.5)",
        },
    ],
};

interface FightRoundGraph {
    fightRounds: Array<number>
}

export function FightRoundGraph({ fightRounds }: FightRoundGraph) {
    const xaxis = fightRounds.map((_, idx) => {
        return idx + 1;
    });
    data.labels = xaxis;
    data.datasets[0].data = fightRounds;

    return <Bar key={uuid().toString()} style={{ minHeight: "50vh", maxHeight: "50vh" }} options={options} data={data} />;
}
