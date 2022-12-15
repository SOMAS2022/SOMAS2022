import csv
import numpy as np
from sklearn.linear_model import LinearRegression
from collections import defaultdict
import random

path = "/mnt/c/Users/alexp/OneDrive - Imperial College London/Year 4/Y4_SOMAS/SOMAS2022/pkg/infra/teams/team1/post_log.csv"
final_weights_path = '/mnt/c/Users/alexp/OneDrive - Imperial College London/Year 4/Y4_SOMAS/SOMAS2022/pkg/infra/teams/team1/final_weights.csv'
lr = 0.001

data = [[],[],[]]
y_coop = [[],[],[]]
y_self = [[],[],[]]

# Episodes:
# actual reward y_coop, y_self
# next q-value given best next state
# Current q-value


### READ POST_LOG
# for given game, group together actions of the same player
# player_actions = defaultdict(list)
# with open(path, newline='') as f:
#     reader = csv.reader(f)

#     for row in reader:
#         game_id = row[0]
#         agent_id = row[1]
#         player_actions[game_id,agent_id].append(row)

# # Create episode structure
# episodes = []
# for play in player_actions.values():
#     for i in range(len(play)-1):
#         decision = play[i]
#         next_decision = play[i+1]
#         print(decision, next_decision)

with open(path, newline='') as f:
    reader = csv.reader(f)

    for row in reader:
        action = int(row[2])
        data[action].append([float(i) for i in
            [row[6],row[7],row[8],row[9],row[3],row[5],row[4]]
        ])
        y_coop[action].append(float(row[12]))
        y_self[action].append(float(row[11]))


## Q-learning approach

# Read final_weights.csv values
weights = []
with open(final_weights_path, newline='') as f:
    reader = csv.reader(f)
    for row in reader:
        weights.append([float(i) for i in row])
coopstrats = weights[:3]
selfstrats = weights[3:]


for action in range(3):
    # print("ACTION: ",action)
    X = np.array(data[action])
    y_coop_action = np.array(y_coop[action])
    y_self_action = np.array(y_self[action])

    

    # w_i <- w_i + lr (r + discount * max_a'(Q(s',a'))-Q(s,a)) * f_i(s,a)

    def update_strat(reward, prev_strat, n=5):
        # Select n plays from game to train on
        random_idx = random.sample(range(len(X)), n)
        for i in random_idx:
            episode = X[i,:]
            state = np.append([1],episode)
            r = reward[i]
            Q = np.dot(np.array(prev_strat),state)
            error = lr * (r - Q)
            for j in range(8):
                prev_strat[j] += min(1,error) * state[j] # too much scaling
        return prev_strat
    
    # Cooperative
    coopstrats[action] = update_strat(y_coop_action, coopstrats[action])
    # Selfish
    coopstrats[action] = update_strat(y_self_action, selfstrats[action])

# write new weights to final_weights.csv
f = open(final_weights_path, 'w')
writer = csv.writer(f)
for s in coopstrats:
    writer.writerow(s)
for s in selfstrats:
    writer.writerow(s)







