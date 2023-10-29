import numpy as np
import matplotlib.pyplot as plt

# Define the range for a and b values
a_values = np.linspace(-2, 1, 4)
b_values = np.linspace(-1, 2, 4)

# Set common plotting parameters
y_max = 5
x_max = 5

# Create a 4x4 grid of subplots
fig, axes = plt.subplots(4, 4, figsize=(15, 15), sharex='col', sharey='row')

# Set x and y labels for rows and columns
for i, a in enumerate(a_values):
    for j, b in enumerate(b_values):
        if i == 0:
            axes[i, j].set_title(f'b={b:.0f}')
        if j == 0:
            axes[i, j].set_ylabel(f'a={a:.0f}')

# Iterate over a and b values
for i, a in enumerate(a_values):
    for j, b in enumerate(b_values):
        y, x = np.ogrid[-y_max:y_max:100j, -x_max:x_max:100j]
        curve = pow(y, 2) - pow(x, 3) - x * a - b
        axes[i, j].contour(x.ravel(), y.ravel(), curve, [0], colors='royalblue')
        axes[i, j].grid()

# Adjust the layout of subplots
plt.tight_layout()

# Save the figure as a PNG file
plt.savefig("docs/images/elliptic_curves_grid.png")
plt.savefig("docs/images/elliptic_curves_grid.svg")

# Display the plot
plt.show()
