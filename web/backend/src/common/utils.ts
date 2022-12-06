import {execSync} from "child_process";

export function getGITCommitHash(format: "long" | "short") {
    if (format == "long") {
        return execSync("git rev-parse HEAD").toString().trim();
    } else {
        return execSync("git rev-parse HEAD").toString().trim().substring(0, 7);
    }
}

// export function stdoutToLogJSON (stdout: string) {
//     const log: Array<battle_summary | game_summary> = [];
//     const array = stdout.split("\n");
//     for (let i = 0; i < array.length; i++) {
//         if (array[i][0] != "{") continue;
//         try {
//             const json = JSON.parse(array[i]) as battle_summary | game_summary;
//             json.time = new Date(json.time);
//             console.log(json);
//             log.push(json);
//         } catch {
//             console.log(`Invalid JSON: ${array[i]}`);
//             continue;
//         }
//     }
//     return log;
// }
