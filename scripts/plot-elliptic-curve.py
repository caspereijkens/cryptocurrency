import numpy as np
import matplotlib.pyplot as plt

a = 0
b = 7

# Set common plotting parameters
y_max = 15
x_max = 9

# Create a 1x2 grid of subplots
fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(16, 6))


y, x = np.ogrid[-y_max:y_max:100j, -x_max:x_max:100j]

cubic_curve = y - pow(x, 3) - x * a - b
cubic_curve_neg = -y - pow(x, 3) - x * a - b
# Plot both equations on ax1
ax1.contour(x.ravel(), y.ravel(), cubic_curve, [0], colors='royalblue')
ax1.contour(x.ravel(), y.ravel(), cubic_curve_neg, [0], colors='red')

ax1.set_xlabel('x')
ax1.set_ylabel('y')
ax1.axhline()
ax1.axvline()
ax1.grid()

# Create the equation title dynamically based on a and b
eq_title = f'y = x^3 {f"- {a}x" if a != 0 else ""} {f"- {b}" if b != 0 else ""}'

# Set the title with the equation
ax1.set_title(f'${eq_title}$')

elliptic_curve = pow(y, 2) - pow(x, 3) - x * a - b
# Plot both equations on ax1
ax2.contour(x.ravel(), y.ravel(), elliptic_curve, [0], colors='royalblue')

ax2.set_xlabel('x')
ax2.set_ylabel('y')
ax2.axhline()
ax2.axvline()
ax2.grid()

# Create the equation title dynamically based on a and b
eq_title = f'y^2 = x^3 {f"- {a}x" if a != 0 else ""} {f"- {b}" if b != 0 else ""}'

# Set the title with the equation
ax2.set_title(f'${eq_title}$')


# Adjust the layout of subplots
plt.tight_layout()


# Save the figure as a PNG file
plt.savefig("docs/images/both_curves.png")
plt.savefig("docs/images/both_curves.svg")

# Display the plot
plt.show()
