import matplotlib
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

matplotlib.use("TkAgg")

iterations = []
mean_rank_diffs = []
average_rank_diffs = []

results_file = open("initial_guess_impact_results", "r")
for row in results_file:
    cells = row.split()
    iterations.append(int(cells[0]))
    mean_rank_diffs.append(float(cells[1]))
    average_rank_diffs.append(float(cells[2]))
results_file.close()

plt.figure(figsize=(20, 10))
plt.xticks(np.linspace(-1000, 6_000, 27))
sns.distplot(iterations, bins=np.linspace(-1000, 10_000, 110), label="Iterations")
plt.legend()
plt.show(block=False)

plt.figure(figsize=(20, 10))
plt.xticks(np.linspace(-2, 6, 27))
sns.distplot(mean_rank_diffs, label="Mean diff")
plt.legend()
plt.show(block=False)

plt.figure(figsize=(20, 10))
plt.xticks(np.linspace(-2, 6, 27))
sns.distplot(average_rank_diffs, label="Average diff")
plt.legend()
plt.show()
