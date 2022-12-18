#!/usr/bin/env python3

import os
import sys
import json
import math
import re
import matplotlib.pyplot as plt
import matplotlib.cm as cm
import numpy as np


def SocialNetwork(frame):
    network_str = frame.get("SocialNetwork")
    network = json.loads(network_str)
    result = []
    for _, item in network.items():
        s = item.get("StrategyScore")
        g = item.get("GoodwillScore")
        sp = item.get("Strategy")
        gp = item.get("Goodwill")
        colour = (3*sp+gp)
        result.append((g, s, sp, gp, colour))

    return result


def main(log_file):
    name = os.path.basename(log_file).split(".")[0]
    print(name)
    IDs = []
    all_logs = []
    with open(log_file, "r") as fp:
        while True:
            line = fp.readline()
            if not line:
                break
            data = json.loads(line)
            all_logs.append(data)
            IDs.append(data.get("ID"))

    ID = IDs[0]
    logs = list(filter(lambda log: log.get("ID") == ID, all_logs))
    n_levels = len(logs)

    # Print for which level
    LEVEL = n_levels

    for level in logs:
        if level.get("LEVEL") == LEVEL:
            network = SocialNetwork(level)
            res = list(zip(*network))
            plt.plot(res[0], res[1], "o")
            plt.axvline(0.2, color='r', linestyle='--')
            plt.axvline(0.8, color='g', linestyle='--')
            plt.axhline(0.2, color='r', linestyle='--')
            plt.axhline(0.8, color='g', linestyle='--')
            plt.xlim(0, 1)
            plt.ylim(0, 1)


                    
                                                
    plt.ylabel("Strategy Score")
    plt.xlabel("Goodwill Score")
    plt.title(f"Agent Categorisation on level {LEVEL}")
    plt.savefig(f"{name}_social_network.png")


if __name__ == "__main__":
    if len(sys.argv) < 1:
        exit(1)

    log_file = sys.argv[1]
    main(log_file)
