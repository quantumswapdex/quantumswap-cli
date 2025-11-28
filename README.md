# quantumswap-cli
CLI Tool for QuantumSwap. QuantumSwap is a DEX that runs on QuantumCoin (Q) blockchain.

## Prerequisites

Set following environment variables:

1) FROM_ADDRESS

## How to Swap Tokens 
Additionally, the `FEE` value used when creating the liquidity pool should be identified.

### Option A) Swapping with option of constant input tokens and minimum output tokens needed
```quantumswap-cli createpool TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE```

`FEE` : Use values 500 for 0.05%, 3000 for 0.3% or 10000 for 1% fee tier

### Option B) Swapping with option of constant output tokens and maximum input spend
```quantumswap-cli getpool TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE```

`FEE` : Use values 500 for 0.05%, 3000 for 0.3% or 10000 for 1% fee tier

## How to create Tokens, Check Balance, transfer etc.?

### Creating a new Token
```dputil createtoken FROM_ADDRESS TOKEN_NAME TOKEN_SYMBOL TOTAL_SUPPLY```

### Checking Token Balance
```dputil tokenbalance CONTRACT_ADDRESS ACCOUNT_ADDRESS```

### Transferring Tokens
```dputil transfertokens CONTRACT_ADDRESS FROM_ADDRESS TO_ADDRESS amount```

### Renouncing Token Ownership
```dputil renouncetokenownership CONTRACT_ADDRESS FROM_ADDRESS```

## How to create a Liquidity Pool and add Liquidity?

### Prerequisites
Set the following environment variables:

1) V3_CORE_FACTORY_CONTRACT_ADDRESS
2) NONFUNGIBLE_POSITION_MANAGER_CONTRACT_ADDRESS
3) SWAP_ROUTER_CONTRACT_ADDRESS
4) FROM_ADDRESS

### 1) Create a Liquidity Pool

```quantumswap-cli createpool TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE```

`FEE` : Use values 500 for 0.05%, 3000 for 0.3% or 10000 for 1% fee tier

### 2) Get the Pool Address
```quantumswap-cli getpool TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE```

`FEE` : Use values 500 for 0.05%, 3000 for 0.3% or 10000 for 1% fee tier

### 2) Initialize the Pool

```quantumswap-cli initializepool POOL_ADDRESS PRICE_IN_TOKEN_B_PER_TOKEN_A TOKEN_A_DECIMALS TOKEN_B_DECIMALS```

`TOKEN_A_DECIMALS` and `TOKEN_B_DECIMALS` should be between, 0 to 18

### 3) Approve TokenA
```quantumswap-cli approve TOKEN_ADDRESS APPROVAL_ADDRESS AMOUNT```

`TOKEN_ADDRESS`: Pass the `TOKEN_A_ADDRESS` 
`AMOUNT` You may give maximum or specific amount

### 4) Approve TokenB
```quantumswap-cli approve TOKEN_ADDRESS APPROVAL_ADDRESS AMOUNT```

`TOKEN_ADDRESS`: Pass the `TOKEN_A_ADDRESS`
`AMOUNT` You may give maximum or specific amount

### 5) Add Liquidity
```quantumswap-deploy addliquidity TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE TICK_LOWER TICK_UPPER AMOUNT_A AMOUNT_B AMOUNT_A_MIN AMOUNT_B_MIN```

`FEE` : Use values 500 for 0.05%, 3000 for 0.3% or 10000 for 1% fee tier


