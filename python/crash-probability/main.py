# Cell 1
import numpy as np
import matplotlib.pyplot as plt
import mplcursors

# Parameters
a = 0.006
n = 5
dt = 0.1

# Initialize y, v, d, f
t = [0.]
y = [0.]
y1 = [0.05]
y2 = [0.10]
v = [1.]
d = [0.]
f = [0.]
f1 = [0.05]
f2 = [0.10]
odds = [100]

i = 1
while odds[i-1] < 50000:
    t.append(t[i-1] + dt)
    d.append(int((t[i-1]+n-0.1)/n) * a)
    v.append(v[i-1] + d[i-1])
    odds.append((v[i-1] + 0.005)*100)
    f.append((1 - (1 - y[0]) / v[i-1] - y[i-1]) / (1 - y[i-1]))
    y.append(f[i-1] * (1-y[i-1]) + y[i-1])
    i += 1

# Create figure for plotting
fig, ax1 = plt.subplots()

# Plot f(t)
line1, = ax1.plot(t, f, label='f(t)')
ax1.set(xlabel='t', ylabel='f(t)', title='Plot of f(t) and d(t)')
ax1.grid()

# Create second y-axis
ax2 = ax1.twinx()
line2, = ax2.plot(t, d, color='red', label='d(t)')
ax2.set_ylabel('d(t)')

# Add interactivity
cursor1 = mplcursors.cursor(line1, hover=True)
@cursor1.connect("add")
def on_add(sel):
    sel.annotation.set_text(f'f(t): {sel.target[1]:.10f}')

cursor2 = mplcursors.cursor(line2, hover=True)
@cursor2.connect("add")
def on_add(sel):
    sel.annotation.set_text(f'd(t): {sel.target[1]:.10f}')



import pandas as pd

# Create a DataFrame
df = pd.DataFrame({
    'odds': odds,
    't': t,
    'y(0) == 0': f,  # Replace y with the actual list
})

# Write the DataFrame to an Excel file
df.to_excel('output.xlsx', index=False)

# Show the plot
plt.show()