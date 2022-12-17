import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Filler,
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
    Filler,
    Legend
);

const options = {
    responsive: true,
    plugins: {
        legend: {
            position: "top" as const,
        },
        title: {
            display: true,
            text: "Histogram",
        },
    },
    scales: {
        x: {
            min: 0,
            max: 100,
            ticks: {
                stepSize: 1
            }
        }
    },
};

const labels: number[] = [];

const data = {
    labels,
    datasets: [
        {
            fill: true,
            label: "Histogram",
            data: [0],
            borderColor: "rgb(53, 162, 235)",
            backgroundColor: "rgba(53, 162, 235, 0.5)",
            xAxisID: "x"
        },
    ],
};

interface HistogramGraph {
    property: string
    values: Array<Array<number>>
    idx: number
    binSize: number
}

export function HistogramGraph({property, values, idx, binSize}: HistogramGraph) {
    data.datasets[0].label = property;
    if (values.length == 0) {
        return <Line key={uuid().toString()} options={options} data={data} style={{ minHeight: "50vh", maxHeight: "50vh" }} />;
    } else {
        const flat = values.flat();
        const min = Math.min(...flat);
        const max = Math.max(...flat);

        const freqs: Record<number, number> = {};
        const binnedFreqs: Record<number, number> = {};
        for(let i = 0; i < values[idx].length; i++) {
            if (!freqs[values[idx][i]]) {
                freqs[values[idx][i]] = 1;
                binnedFreqs[values[idx][i]] = 1;
            } else {
                freqs[values[idx][i]] += 1;
                binnedFreqs[values[idx][i]] += 1;
            }
        }
        console.log(values[idx], freqs);

        const keys = Object.keys(freqs) as unknown as number[];
        console.log(keys, Math.min(...keys));

        if (binSize > 1) {
            console.log(keys);
            for (let i = 0; i < keys.length; i += binSize) {
                let sum = binnedFreqs[keys[i]];
                for(let j = i+1; j < i+binSize; j++) {
                    sum += binnedFreqs[keys[j]];
                }
                console.log(i, sum);
                for (let j = i; j < i+binSize; j++) {
                    binnedFreqs[keys[j]] = sum;
                }
            }
            console.log(freqs, binnedFreqs);
        }

        for (let i = min; i < Math.min(...keys) - binSize; i++) {
            // console.log(i);
            binnedFreqs[i] = 0;
        }
        for (let i = Math.max(...keys) + binSize + 1; i <= max; i++) {
            // console.log(i);
            binnedFreqs[i] = 0;
        }

        data.labels = Object.keys(binnedFreqs) as unknown as number[];
        data.datasets[0].data = Object.values(binnedFreqs);
        return <Line key={uuid().toString()} options={options} data={data} style={{ minHeight: "50vh", maxHeight: "50vh" }} />;
    }
}
