import numpy as np
import matplotlib.pyplot as plt

a = -2
b = 2

# Set common plotting parameters
x_max = 3
y_max = 5

# Create a 1x2 grid of subplots
fig, ax1 = plt.subplots(1, 1, figsize=(8, 6))

y, x = np.ogrid[-y_max:y_max:100j, -x_max:x_max:100j]

# Set the title with the equation
ax1.set_title('Point addition')

elliptic_curve = pow(y, 2) - pow(x, 3) - x * a - b
# Plot both equations on ax1
ax1.contour(x.ravel(), y.ravel(), elliptic_curve, [0], colors='royalblue')
ax1.set_xlabel('x')
ax1.set_ylabel('y')
ax1.set_ylim(y_max, -y_max)
ax1.axhline()
ax1.axvline()
ax1.grid()


# Define two points P and Q
P = (-1, np.sqrt((-1) ** 3 + a * (-1) + b))  # Point P
Q = (2, np.sqrt(2 ** 3 + a * 2 + b))  # Point Q

# Plot the points P and Q
ax1.plot(P[0], P[1], 'ro', label='P')
ax1.plot(Q[0], Q[1], 'go', label='Q')

# Draw a line connecting P and Q
ax1.plot([P[0], Q[0]], [P[1], Q[1]], 'k--', label='P + Q')

ax1.legend()

# Adjust the layout of subplots
plt.tight_layout()


# Save the figure as a PNG file
# plt.savefig("docs/images/cubic-to-elliptic.png")
# plt.savefig("docs/images/cubic-to-elliptic.svg")

# Display the plot
plt.show()
