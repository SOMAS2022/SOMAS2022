import matplotlib.pyplot as plt
import matplotlib.cbook as cbook

import numpy as np
import pandas as pd


# file name here
filename = "19-58-gameLog.csv"
data = pd.read_csv('logCSV/'+filename)
data.set_index("level")
data.reset_index(drop=True, inplace=True)

plt.plot(data["level"],data["total agents alive"])
plt.title('Number of Agents per Level')
plt.ylabel('Number of Agents')
plt.xlabel('Level')
plt.legend(loc="upper right")
# plt.show()

plt.plot(data["average health"], label = 'Average Agent Health')
plt.plot(data["average stamina"], label = 'Average Agent Stamina')
# dosen't show any var
# plt.plot(data["average attack"], label = 'Average Agent Attack')
# plt.plot(data["average defense"], label = 'Average Agent Defense')
# plt.yscale("log")
plt.title('Evolution of agent stats per level')
plt.ylabel('Stats Value')
plt.xlabel('Level')
plt.legend()
# plt.show()

plt.plot(data["count selfless"], label = 'Proportion of Selfless Agents')
plt.plot(data["count selfish"], label = 'Proportion of Selfish agents')
plt.plot(data["count collective"], label = 'Proportion of Collective Agents')
# plt.plot(history.history['val_loss'], label='Base Model Validation loss')
# plt.yscale("log")
plt.title('Agent Personality Distribution per Level')
plt.ylabel('Agents')
plt.xlabel('Level')
plt.legend()
# plt.show()

# think adding amount of agents is interesting bc people become increasingly selfish as less people remain
# NB: sanction code not finished in this dataset
# plt.plot(data["total agents alive"], label = 'Amount of Agents')
plt.plot(data["average personality"], label = 'Average Agent Personality')
plt.plot(data["average sanctioned"], label = 'Average Agent Sanctioned')
plt.title('Agent Personality and Sanction per Level')
plt.ylabel('Stat Amount')
plt.xlabel('Level')
plt.legend()
plt.show()