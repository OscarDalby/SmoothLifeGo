Example Inputs for Py SmoothLife functions:

# LogisticThreshold
Example inputs: (x=0.5, x0=0.5, alpha=0.1)
Py output:
logisticThreshold(x=-1, x0=0.5, alpha=0.1) => 8.75651076269652e-27
logisticThreshold(x=0.1, x0=0.5, alpha=0.1) => 1.12535162055095e-07
logisticThreshold(x=0.25, x0=0.5, alpha=0.1) => 4.5397868702434395e-05
logisticThreshold(x=0.5, x0=0.5, alpha=0.1) => 0.5
logisticThreshold(x=0.75, x0=0.5, alpha=0.1) => 0.9999546021312976
logisticThreshold(x=1, x0=0.5, alpha=0.1) => 0.9999999979388463
logisticThreshold(x=2, x0=0.5, alpha=0.1) => 1.0

# HardThreshold
Example inputs: x=0.5, x0=0.5
Py output:
hardThreshold(x=-1, x0=0.5) => False
hardThreshold(x=0.1, x0=0.5) => False
hardThreshold(x=0.25, x0=0.5) => False
hardThreshold(x=0.5, x0=0.5) => False
hardThreshold(x=0.75, x0=0.5) => True
hardThreshold(x=1, x0=0.5) => True
hardThreshold(x=2,x0=0.5) => True

# LinearisedThreshold
Example inputs: x=0.5, x0=0.5, alpha=0.1
Py output:
LinearisedThreshold(x=-1, x0=0.5, alpha=0.1) => 0.0
LinearisedThreshold(x=0.1, x0=0.5, alpha=0.1) => 0.0
LinearisedThreshold(x=0.25, x0=0.5, alpha=0.1) => 0.0
LinearisedThreshold(x=0.5, x0=0.5, alpha=0.1) => 0.5
LinearisedThreshold(x=0.75, x0=0.5, alpha=0.1) => 1.0
LinearisedThreshold(x=1, x0=0.5, alpha=0.1) => 1.0
LinearisedThreshold(x=2, x0=0.5, alpha=0.1) => 1.0

# LogisticInterval
Example inputs: x=0.5, a=0.3, b=0.7, alpha=0.1
Py output:
LogisticInterval(x=0.3, a=0.3, b=0.7, alpha=0.1) => 0.49999994373241896
LogisticInterval(x=0.4, a=0.3, b=0.7, alpha=0.1) => 0.9820077563737207
LogisticInterval(x=0.5, a=0.3, b=0.7, alpha=0.1) => 0.9993294121987771
LogisticInterval(x=0.6, a=0.3, b=0.7, alpha=0.1) => 0.9820077563737207
LogisticInterval(x=0.7, a=0.3, b=0.7, alpha=0.1) => 0.49999994373241896

# LinearisedInterval
Example inputs: x=0.5, a=0.3, b=0.7, alpha=0.1
LinearisedInterval(x=0.3, a=0.3, b=0.7, alpha=0.1) => 0.5
LinearisedInterval(x=0.4, a=0.3, b=0.7, alpha=0.1) => 1.0
LinearisedInterval(x=0.5, a=0.3, b=0.7, alpha=0.1) => 1.0
LinearisedInterval(x=0.6, a=0.3, b=0.7, alpha=0.1) => 1.0
LinearisedInterval(x=0.7, a=0.3, b=0.7, alpha=0.1) => 0.5

# Lerp
Example inputs: a=0, b=1, t=0.5
Lerp(a=0,b=1,t=0.5) => 0.5

# BasicRules.s
Example inputs: n=0.5, m=0.5, field=np.array([[0.1, 0.2], [0.3, 0.4]])

# ExtensiveRules.sigmoid_ab
Example inputs: x=0.5, a=0.3, b=0.7

# ExtensiveRules.sigmoid_mix
Example inputs: x=0.5, y=0.7, m=0.5

# ExtensiveRules.s
Example inputs: n=0.5, m=0.5, field=np.array([[0.1, 0.2], [0.3, 0.4]])

# SmoothTimestepRules.s
Example inputs: n=0.5, m=0.5, field=np.array([[0.1, 0.2], [0.3, 0.4]])

# antialiasedCircle
Example inputs: size=(100, 100), radius=50, roll=True, logres=2

# Multipliers
Example inputs: size=(100, 100), inner_radius=7.0, outer_radius=21.0

# SmoothLife
Example inputs: height=100, width=100

# SmoothLife.step
Example inputs: 

# SmoothLife.add_speckles
Example inputs: count=10, intensity=1.0