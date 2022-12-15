## WIP Qlearning

import csv
import numpy as np
from sklearn.linear_model import LinearRegression
from collections import defaultdict

path = "/mnt/c/Users/alexp/OneDrive - Imperial College London/Year 4/Y4_SOMAS/SOMAS2022/pkg/infra/teams/team1/post_log.csv"
final_weights_path = '/mnt/c/Users/alexp/OneDrive - Imperial College London/Year 4/Y4_SOMAS/SOMAS2022/pkg/infra/teams/team1/final_weights.csv'
lr = 0.1
discount = 0.9

data = [[],[],[]]
y_coop = [[],[],[]]
y_self = [[],[],[]]

# Episodes:
# actual reward y_coop, y_self
# next q-value given best next state
# Current q-value

episodes = []


### READ POST_LOG
with open(path, newline='') as f:
    reader = csv.reader(f)

    for row in reader:
        action = int(row[2])
        data[action].append([float(i) for i in
            [row[6],row[7],row[8],row[9],row[3],row[5],row[4]]
        ])
        y_coop[action].append(float(row[12]))
        y_self[action].append(float(row[11]))



## CALCULATE STRATEGIES
# coopstrats = []
# selfstrats = []
# for action in range(3):
#     # print("ACTION: ",action)
#     X = np.array(data[action])
#     y_coop_action = np.array(y_coop[action])
#     y_self_action = np.array(y_self[action])
#     # print(len(X),y_coop,y_self)

#     reg_coop = LinearRegression().fit(X, y_coop_action)
#     reg_self = LinearRegression().fit(X, y_self_action)
#     # print(reg_coop.intercept_)

#     weights_coop = np.append(reg_coop.intercept_, reg_coop.coef_)
#     coopstrats.append(weights_coop)
#     weights_self = np.append(reg_self.intercept_, reg_self.coef_)
#     selfstrats.append(weights_self)
#     print(weights_coop)
#     print(weights_self)
    # print(X,y)

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
    # Cooperative

    for i,episode in X:
        r = y_coop_action[i]
        # Q = max(coopstrats[action] * 
        # next_state = 
        Q_next = max(coopstrats[action] * next_state)
        error = lr * (r + discount * Q_next - Q)
        for j in range(8):
            coopstrats[action][j] += error * episode[j]



    # reg_coop = LinearRegression().fit(X, y_coop_action)
    # reg_self = LinearRegression().fit(X, y_self_action)
    # # print(reg_coop.intercept_)

    # weights_coop = np.append(reg_coop.intercept_, reg_coop.coef_)
    # coopstrats.append(weights_coop)
    # weights_self = np.append(reg_self.intercept_, reg_self.coef_)
    # selfstrats.append(weights_self)
    # print(weights_coop)
    # print(weights_self)
    # # print(X,y)



# write new weights to final_weights.csv
f = open(final_weights_path, 'w')
writer = csv.writer(f)
for s in coopstrats:
    writer.writerow(s)
for s in selfstrats:
    writer.writerow(s)







