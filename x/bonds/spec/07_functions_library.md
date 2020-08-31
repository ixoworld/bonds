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

<img alt="power function price" src="./img/power1.png" height="20"/>

Integral (used as reserve function):

<img alt="power function reserve" src="./img/power2.png" height="40"/>

### Logistic Function (sigmoid)

Function (used as pricing function):

<img alt="sigmoid function price" src="./img/sigmoid1.png" height="80"/>

Integral (used as reserve function):

<img alt="sigmoid function reserve" src="./img/sigmoid2.png" height="55"/>

### Augmented Bonding Curves (augmented)

Initial reserve:

<img alt="initial reserve" src="./img/augmented1.png" height="20"/>

Initial supply:

<img alt="initial supply" src="./img/augmented2.png" height="20"/>

Constant power function invariant:

<img alt="constant power" src="./img/augmented3.png" height="40"/>

Invariant function:

<img alt="invariant function" src="./img/augmented4.png" height="55"/>

Pricing function:

<img alt="pricing function" src="./img/augmented5.png" height="55"/>

Reserve function:

<img alt="reserve function" src="./img/augmented6.png" height="50"/>

Ref: https://medium.com/giveth/deep-dive-augmented-bonding-curves-3f1f7c1fa751

### Constant Product Function (swapper)

Reserve function:

<img alt="swapper function" src="./img/swapper.png" height="20"/>
