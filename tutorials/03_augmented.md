# Creating an augmented bonding curve

Throughout this tutorial, some knowledge around [Augmented Bonding Curves](https://medium.com/giveth/deep-dive-augmented-bonding-curves-3f1f7c1fa751) will be assumed.

## Contents

- [Bond Configuration](#bond-configuration)

## Bond Configuration

### Curve Function

In this tutorial, an augmented function bond will be created. The augmented function implemented by the Bonds module is shown below, where `y` represents the price per bond token, in reserve tokens, for a specific supply `x` of bond tokens:

- Initial reserve:
  <img alt="drawing" src="./img/augmented1.png" height="20"/>
- Initial supply:
  <img alt="drawing" src="./img/augmented2.png" height="20"/>
- Constant power function invariant:
  <img alt="drawing" src="./img/augmented3.png" height="40"/>
- Invariant function:
  <img alt="drawing" src="./img/augmented4.png" height="55"/>
- Pricing function:
  <img alt="drawing" src="./img/augmented5.png" height="55"/>
- Reserve function:
  <img alt="drawing" src="./img/augmented6.png" height="50"/>

Ref: https://medium.com/giveth/deep-dive-augmented-bonding-curves-3f1f7c1fa751
