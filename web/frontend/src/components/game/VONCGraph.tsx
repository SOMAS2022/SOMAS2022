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
import { Bar } from "react-chartjs-2";
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
            text: "Vote of No Confidence Result",
        },
    },
    scales: {
        x: {
            stacked: true,
        },
        y: {
            stacked: true,
        },
    },
};

const labels = [0];

const data = {
    labels: labels,
    datasets: [
        {
            label: "For",
            data: [0],
            backgroundColor: "rgba(75, 192, 192, 0.5)",
        },
        {
            label: "Against",
            data: [0],
            borderColor: "rgb(53, 162, 235)",
            backgroundColor: "rgba(255, 99, 132, 0.5)",
        },
        {
            label: "Abstain",
            data: [0],
            borderColor: "rgb(53, 162, 235)",
            backgroundColor: "rgba(53, 162, 235, 0.5)",
        },
    ],
};

interface VONCGraphProps {
    forVotes: Array<number>
    againstVotes: Array<number>
    abstainVotes: Array<number>
}

export function VONCGraph({ forVotes, againstVotes, abstainVotes }: VONCGraphProps) {
    const xaxis = forVotes.map((_, idx) => {
        return idx + 1;
    });
    data.labels = xaxis;
    data.datasets[0].data = forVotes;
    data.datasets[1].data = againstVotes;
    data.datasets[2].data = abstainVotes;

    return <Bar key={uuid().toString()} options={options} data={data} style={{ minHeight: "50vh", maxHeight: "50vh" }} />;
}
