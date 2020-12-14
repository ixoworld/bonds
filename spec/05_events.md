# Events

The bonds module emits the following events:

## EndBlocker

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| order\_cancel | bond | {token} |
| order\_cancel | order\_type | {orderType} |
| order\_cancel | address | {address} |
| order\_cancel | cancel\_reason | {cancelReason} |
| order\_fulfill | bond | {token} |
| order\_fulfill | order\_type | {orderType} |
| order\_fulfill | address | {address} |
| order\_fulfill | tokensMinted | {tokensMinted} |
| order\_fulfill | chargedPrices | {chargedPrices} |
| order\_fulfill | chargedFees | {chargedFees} |
| order\_fulfill | returnedToAddress | {returnedToAddress} |
| state\_change | bond | {token} |
| state\_change | old\_state | {oldState} |
| state\_change | new\_state | {newState} |

## Handlers

### MsgCreateBond

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| create\_bond | bond | {token} |
| create\_bond | name | {name} |
| create\_bond | description | {description} |
| create\_bond | function\_type | {functionType} |
| create\_bond | function\_parameters \[0\] | {functionParameters} |
| create\_bond | reserve\_tokens \[1\] | {reserveTokens} |
| create\_bond | tx\_fee\_percentage | {txFeePercentage} |
| create\_bond | exit\_fee\_percentage | {exitFeePercentage} |
| create\_bond | fee\_address | {feeAddress} |
| create\_bond | max\_supply | {maxSupply} |
| create\_bond | order\_quantity\_limits | {orderQuantityLimits} |
| create\_bond | sanity\_rate | {sanityRate} |
| create\_bond | sanity\_margin\_percentage | {sanityMarginPercentage} |
| create\_bond | allow\_sells | {allowSells} |
| create\_bond | signers \[2\] | {signers} |
| create\_bond | batch\_blocks | {batchBlocks} |
| create\_bond | state | {state} |
| message | module | bonds |
| message | action | create\_bond |
| message | sender | {senderAddress} |

* \[0\] Example formatting: `"{m:12,n:2,c:100}"`
* \[1\] Example formatting: `"[res,rez]"`
* \[2\] Example formatting: `"[ADDR1,ADDR2]"`

### MsgEditBond

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| edit\_bond | bond | {token} |
| edit\_bond | name | {name} |
| edit\_bond | description | {description} |
| edit\_bond | order\_quantity\_limits | {orderQuantityLimits} |
| edit\_bond | sanity\_rate | {sanityRate} |
| edit\_bond | sanity\_margin\_percentage | {sanityMarginPercentage} |
| message | module | bonds |
| message | action | edit\_bond |
| message | sender | {senderAddress} |

### MsgBuy

#### First Buy for Swapper Function Bond

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| init\_swapper | bond | {token} |
| init\_swapper | amount | {amount} |
| init\_swapper | charged\_prices | {chargedPrices} |
| message | module | bonds |
| message | action | buy |
| message | sender | {senderAddress} |

#### Otherwise

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| buy | bond | {token} |
| buy | amount | {amount} |
| buy | max\_prices | {maxPrices} |
| order\_cancel | bond | {token} |
| order\_cancel | order\_type | {orderType} |
| order\_cancel | address | {address} |
| order\_cancel | cancel\_reason | {cancelReason} |
| message | module | bonds |
| message | action | buy |
| message | sender | {senderAddress} |

### MsgSell

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| sell | bond | {token} |
| sell | amount | {amount} |
| message | module | bonds |
| message | action | buy |
| message | sender | {senderAddress} |

### MsgSwap

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| swap | bond | {token} |
| swap | amount | {amount} |
| swap | from\_token | {fromToken} |
| swap | to\_token | {toToken} |
| message | module | bonds |
| message | action | swap |
| message | sender | {senderAddress} |

### MsgMakeOutcomePayment

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| make\_outcome\_payment | bond | {token} |
| make\_outcome\_payment | address | {senderAddress} |
| message | module | bonds |
| message | action | make\_outcome\_payment |
| message | sender | {senderAddress} |

### MsgWithdrawShare

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| withdraw\_share | bond | {token} |
| withdraw\_share | address | {recipientAddress} |
| withdraw\_share | amount | {reserveOwed} |
| message | module | bonds |
| message | action | withdraw\_share |
| message | sender | {recipientAddress} |

