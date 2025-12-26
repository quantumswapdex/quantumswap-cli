# quantumswap-cli
CLI Tool for QuantumSwap. QuantumSwap is a DEX that runs on QuantumCoin (Q) blockchain. You also need `dputil` CLI tool from https://github.com/quantumcoinproject/quantum-coin-go/releases

## Prerequisites

Set following environment variables:

1) `DP_RAW_URL`, `DP_KEY_FILE`
2) Create two addresses with enough Q for gas, one for the `token creator` and one for the `token swapper`. 
3) For each subsequent command below, set the `DP_KEY_FILE`, `DP_KEY_FILE_DIR`, `FROM_ADDRESS` to the corresponding values based on whether `token creator` or `token swapper` is used.
4) `WQ_CONTRACT_ADDRESS` (Wrapped Q)
5) `V2_CORE_FACTORY_CONTRACT_ADDRESS`
6) `SWAP_ROUTER_V2_CONTRACT_ADDRESS`
7) To check the transaction status at any step, use `dputil txn TXN_HASH` and ensure the receipt status is `0x1`.

## Example for creating tokens, token pairs, adding liquidity

The below instructions are for Windows. For Linux, Mac, use the `export` command instead of the `set` command and change to `$` prefix instead of `% %` for setting environment variables.

`set FROM_ADDRESS=%TOKEN_CREATOR_ADDRESS%`

### Create new tokens

Run the following two commands and note down the contract address. If your goal is to add liquidity between Token and Q, set `TOKEN_B_ADDRESS` value to `WQ_CONTRACT_ADDRESS` (Wrapped Q), instead of creating new token.

`dputil createtoken %FROM_ADDRESS% "Quantum Shiba" "qshib" 1000000000`

`dputil createtoken %FROM_ADDRESS% "Doge Protocol" "dogep" 1000000000`

#### Set token environment variables

Set the environment variables to the contract addresses from above step. 

`set TOKEN_A_ADDRESS=0x7B385c8525D707c8444a95674d467B830e8a3041c5d6458d0CF8E1c4FfefdEfB`
`set TOKEN_B_ADDRESS=0xCeF0799Ccd42A95AaC1f6b9db30F89255F7b116d429Ef27021f31cBB8B01B143`

#### Check token balance

Ensure balance matches the values above when the token was created.

`dputil tokenbalance %TOKEN_A_ADDRESS% %FROM_ADDRESS%`

`dputil tokenbalance %TOKEN_B_ADDRESS% %FROM_ADDRESS%`

### Create a Liquidity Pair

`quantumswap-cli createpair %TOKEN_A_ADDRESS% %TOKEN_B_ADDRESS%`

Note down the pair address by running the following command. 
`quantumswap-cli getpair %TOKEN_A_ADDRESS% %TOKEN_B_ADDRESS%`

Set the PAIR_ADDRESS environment variable to the output value from above command.
`set PAIR_ADDRESS=0x68e8Ac81Dd2Ef3F7dCF5c40ff9A9a4ff09484e064C07b7fdfD78839697f59074`

### Add Liquidity

#### Approve the tokens for adding liquidity

Approval should be given to the `SWAP_ROUTER_V2_CONTRACT_ADDRESS` contract address.

`dputil tokenapprove %TOKEN_A_ADDRESS% %FROM_ADDRESS% %SWAP_ROUTER_V2_CONTRACT_ADDRESS% 1000000`

`dputil tokenapprove %TOKEN_B_ADDRESS% %FROM_ADDRESS% %SWAP_ROUTER_V2_CONTRACT_ADDRESS% 1000000`

#### Check token allowance for the swap router contract

After giving approval, validate the approval has been given. The output values should match the number of tokens approved.

`dputil tokenallowance %TOKEN_A_ADDRESS% %FROM_ADDRESS% %SWAP_ROUTER_V2_CONTRACT_ADDRESS%`

`dputil tokenallowance %TOKEN_B_ADDRESS% %FROM_ADDRESS% %SWAP_ROUTER_V2_CONTRACT_ADDRESS%`

#### Adding liquidity

First set example token ratio environment variables. Adjust values as desired. Below is example for 1:1.

`set AMOUNT_A=10000`
`set AMOUNT_B=10000`
`set AMOUNT_A_MIN=1`
`set AMOUNT_B_MIN=1`

Add the liquidity.
`quantumswap-cli addliquidityv2 %TOKEN_A_ADDRESS% %TOKEN_B_ADDRESS% %AMOUNT_A% %AMOUNT_B% %AMOUNT_A_MIN% %AMOUNT_B_MIN%`

Check whether liquidity has been added, by checking balance in `PAIR_ADDRESS` and whether the token balance has decreased from `token creator`

`dputil tokenbalance %TOKEN_A_ADDRESS% %PAIR_ADDRESS%`
`dputil tokenbalance %TOKEN_B_ADDRESS% %PAIR_ADDRESS%`

`dputil tokenbalance %TOKEN_A_ADDRESS% %FROM_ADDRESS%`
`dputil tokenbalance %TOKEN_B_ADDRESS% %FROM_ADDRESS%`

### Demonstration of Swapping

First, send token to `TOKEN_SWAPPER_ADDRESS` from the `TOKEN_CREATOR_ADDRESS` (or any address that has the tokens), to demonstrate swapping of tokens.

`dputil transfertokens %TOKEN_A_ADDRESS% %TOKEN_CREATOR_ADDRESS% %TOKEN_SWAPPER_ADDRESS%`

Check balance of the swapper.
`dputil tokenbalance %TOKEN_A_ADDRESS% %TOKEN_SWAPPER_ADDRESS%`

#### Approve the tokens for swapping

Switch the `DP_KEY_FILE`, `DP_KEY_FILE_DIR` environment variables to that of the `TOKEN_SWAPPER_ADDRESS`.

Approval should be given to the `SWAP_ROUTER_V2_CONTRACT_ADDRESS` contract address.

`dputil tokenapprove %TOKEN_A_ADDRESS% %TOKEN_SWAPPER_ADDRESS% %SWAP_ROUTER_V2_CONTRACT_ADDRESS% 1000000`

#### Check TokenA allowance for the swap router contract

After giving approval, validate the approval has been given. The output values should match the number of tokens approved.

`dputil tokenallowance %TOKEN_A_ADDRESS% %TOKEN_SWAPPER_ADDRESS% %SWAP_ROUTER_V2_CONTRACT_ADDRESS%`

#### Swap the tokens

Adjust the following values as desired.
`set AMOUNT_IN=100`
`set AMOUNT_OUT_MIN=1`
`set FROM_ADDRESS=%TOKEN_SWAPPER_ADDRESS%`

`quantumswap-cli swapexacttokensFortokens %TOKEN_A_ADDRESS% %TOKEN_B_ADDRESS% %AMOUNT_IN% %AMOUNT_OUT_MIN%`

Now check balance of both tokens for `TOKEN_SWAPPER_ADDRESS` and `PAIR_ADDRESS`. TokenA should have decreased for the swapper, TokenB should have increased, while its vice versa for the `PAIR_ADDRESS`

`dputil tokenbalance %TOKEN_A_ADDRESS% %TOKEN_SWAPPER_ADDRESS%`
`dputil tokenbalance %TOKEN_B_ADDRESS% %TOKEN_SWAPPER_ADDRESS%`

`dputil tokenbalance %TOKEN_A_ADDRESS% %PAIR_ADDRESS%`
`dputil tokenbalance %TOKEN_B_ADDRESS% %PAIR_ADDRESS%`

