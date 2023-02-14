import os
import json

OUTPUT_FILE_LOCATION = "./output/output.json"
RUN_COMMAND = "go run ./pkg/infra "
NUM_ITERATIONS = 1


def parseJSON(data):
    level_data = data["Levels"]
    return len(level_data)


def fixedLength():
    duration_comp = []
    for duration in range(0, 2):
        avg_level_reached = 0
        for _ in range(NUM_ITERATIONS):
            os.system(RUN_COMMAND + f"-fSanc={duration}")

            with open(OUTPUT_FILE_LOCATION) as OUTPUT_JSON:
                DATA = json.load(OUTPUT_JSON)
                lvl = parseJSON(DATA)
                avg_level_reached += lvl
        avg_level_reached /= NUM_ITERATIONS
        duration_comp.append(avg_level_reached)

    for duration, score in enumerate(duration_comp):
        print(f"duration:{duration}, score:{score}")

def dynamic():
    return

def graduated():
    return

if __name__ == "__main__":
    fixedLength()
    # dynamic()
    # graduated()
