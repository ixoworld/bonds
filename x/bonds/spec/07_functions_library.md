# Functions Library
The Bonds Module is deployed with a built-in set of libraries for commonly-used algorithmic pricing and reserve functions. It also includes algorithmic application logic and features, such as *Augmented Bonding*. Additional functions can be added to the Library through SDK updates. This requires a formal process of governance to approve updates, to assure the integrity of these functions.

## Function Types
The following function types will be included in the standard Bonds SDK Module:
* Power (exponential)
* Logistic (sigmoidal)
* Constant Product (swapper)
Algorithmic Applications include:
* Alpha Bonds (Risk-adjusted bonding)
* Innovation Bonds (offers bond shareholders contingent rights to future IP rights and/or revenues)
* Impact Bonds (offers bond shareholders contingent rights to success-based outcomes payments and/or rewards)

### Exponential Function (power)

Function (used as pricing function):

<img alt="drawing" src="./img/power1.png"/>

Integral (used as reserve function):

<img alt="drawing" src="./img/power2.png"/>

### Logistic Function (sigmoid)

Function (used as pricing function):

<img alt="drawing" src="./img/sigmoid1.png"/>

Integral (used as reserve function):

<img alt="drawing" src="./img/sigmoid2.png"/>

### Augmented Bonding Curves (augmented)

Initial reserve and initial supply:

<img alt="drawing" src="./img/augmented1.png"/>
<img alt="drawing" src="./img/augmented2.png"/>

Constant power function invariant and invariant function:

<img alt="drawing" src="./img/augmented3.png"/>
<img alt="drawing" src="./img/augmented4.png"/>

Pricing function:

<img alt="drawing" src="./img/augmented5.png"/>

Reserve function:

<img alt="drawing" src="./img/augmented6.png"/>

Ref: https://medium.com/giveth/deep-dive-augmented-bonding-curves-3f1f7c1fa751

### Constant Product Function (swapper)

Reserve function:

<img alt="drawing" src="./img/swapper.png"/>