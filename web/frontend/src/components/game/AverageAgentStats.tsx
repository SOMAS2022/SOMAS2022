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
            text: "Average agent AT per level",
        },
    },
};

const labels = [0];

const data = {
    labels: labels,
    datasets: [
        {
            label: "Average agent HP per level",
            data: [0],
            borderColor: "rgb(255, 99, 132)",
            backgroundColor: "rgba(255, 99, 132, 0.5)",
        },
        {
            label: "Average agent AT per level",
            data: [0],
            borderColor: "rgb(43, 232, 167)",
            backgroundColor: "rgba(43, 232, 167, 0.5)",
        },
        {
            label: "Average agent SH per level",
            data: [0],
            borderColor: "rgb(149, 80, 79)",
            backgroundColor: "rgba(149, 80, 79, 0.5)",
        },
        {
            label: "Average agent ST per level",
            data: [0],
            borderColor: "rgb(206, 4, 249)",
            backgroundColor: "rgba(206, 4, 249, 0.5)",
        },
    ],
};

interface AverageAgentStatsProps {
    HP: Array<number>
    AT: Array<number>
    SH: Array<number>
    ST: Array<number>
}

export function AverageAgentStats({ HP, AT, SH, ST }: AverageAgentStatsProps) {
    const xaxis = HP.map((_, idx) => {
        return idx + 1;
    });
    data.labels = xaxis;
    data.datasets[0].data = HP;
    data.datasets[1].data = AT;
    data.datasets[2].data = SH;
    data.datasets[3].data = ST;

    return <Line key={uuid().toString()} style={{ minHeight: "50vh", maxHeight: "50vh" }} options={options} data={data} />;
}
