package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"quantumswap-cli/contracts/core"
	"quantumswap-cli/contracts/nonfungiblepositionmanager"
	"quantumswap-cli/contracts/swaprouter"
	"quantumswap-cli/contracts/v3pool"
	"time"

	"github.com/quantumcoinproject/quantum-coin-go/accounts/abi/bind"
	"github.com/quantumcoinproject/quantum-coin-go/common"
	"github.com/quantumcoinproject/quantum-coin-go/core/types"
	"github.com/quantumcoinproject/quantum-coin-go/crypto/cryptobase"
	"github.com/quantumcoinproject/quantum-coin-go/ethclient"
	"github.com/quantumcoinproject/quantum-coin-go/params"
)

// getPriceFromTick converts a tick to a price using the formula: 1.0001^tick
// tick is a signed 24-bit integer (int24 in Solidity, int32 in Go)
// Uses big.Float for high precision calculations
func getPriceFromTick(tick int32) *big.Float {
	// Base: 1.0001
	base := big.NewFloat(1.0001)

	// If tick is 0, return 1.0
	if tick == 0 {
		return big.NewFloat(1.0)
	}

	// Use efficient exponentiation by squaring
	result := big.NewFloat(1.0)
	absTick := int32(tick)
	if tick < 0 {
		absTick = -tick
	}

	// Binary exponentiation (exponentiation by squaring)
	currentPower := new(big.Float).Set(base)
	for absTick > 0 {
		if absTick&1 == 1 {
			result.Mul(result, currentPower)
		}
		currentPower.Mul(currentPower, currentPower)
		absTick >>= 1
	}

	// If tick was negative, take the reciprocal
	if tick < 0 {
		result.Quo(big.NewFloat(1.0), result)
	}

	return result
}

// getPriceFromTickFloat64 is a simpler version using float64 (less precise but faster)
func getPriceFromTickFloat64(tick int32) float64 {
	return math.Pow(1.0001, float64(tick))
}

// getTickFromPrice converts a price to a tick using the formula: log(price) / log(1.0001)
// price is a uint256 in Solidity, represented as *big.Float in Go
func getTickFromPrice(price *big.Float) int32 {
	// Calculate log(price) / log(1.0001)
	// Using natural logarithm
	priceFloat, _ := price.Float64()
	logBase := math.Log(1.0001)

	if priceFloat <= 0 {
		return 0 // Invalid price
	}

	tickFloat := math.Log(priceFloat) / logBase

	// Round to nearest integer (int24)
	return int32(math.Round(tickFloat))
}

// getTickFromPriceFloat64 is a simpler version using float64
func getTickFromPriceFloat64(price float64) int32 {
	if price <= 0 {
		return 0
	}

	logBase := math.Log(1.0001)
	tickFloat := math.Log(price) / logBase

	return int32(math.Round(tickFloat))
}

// calculateSqrtPriceX96 calculates the square root price scaled by 2^96
// price: the price value
// decimals0: number of decimals for token0
// decimals1: number of decimals for token1
// Returns: the square root price as a uint160 (represented as *big.Int)
func calculateSqrtPriceX96(price *big.Int, decimals0, decimals1 uint8) *big.Int {
	// Adjust for decimal differences between tokens
	decimalsDiff := new(big.Int).Sub(
		big.NewInt(int64(decimals0)),
		big.NewInt(int64(decimals1)),
	)

	// Calculate 10^(decimals0 - decimals1)
	ten := big.NewInt(10)
	decimalsMultiplier := new(big.Int).Exp(ten, decimalsDiff, nil)

	// adjustedPrice = price * 10^(decimals0 - decimals1)
	adjustedPrice := new(big.Int).Mul(price, decimalsMultiplier)

	// Calculate square root of (adjustedPrice * 10^18)
	tenTo18 := new(big.Int).Exp(ten, big.NewInt(18), nil)
	priceForSqrt := new(big.Int).Mul(adjustedPrice, tenTo18)

	sqrtPrice := sqrt(priceForSqrt)

	// Scale by 2^96
	twoTo96 := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	numerator := new(big.Int).Mul(sqrtPrice, twoTo96)

	// Divide by 10^9
	tenTo9 := new(big.Int).Exp(ten, big.NewInt(9), nil)
	result := new(big.Int).Div(numerator, tenTo9)

	return result
}

// sqrt calculates the square root of x using the Babylonian method
func sqrt(x *big.Int) *big.Int {
	if x.Sign() == 0 {
		return big.NewInt(0)
	}

	// z = (x + 1) / 2
	one := big.NewInt(1)
	z := new(big.Int).Add(x, one)
	z.Div(z, big.NewInt(2))

	y := new(big.Int).Set(x)

	// Babylonian method: while z < y
	for z.Cmp(y) < 0 {
		y.Set(z)
		// z = (x / z + z) / 2
		temp := new(big.Int).Div(x, z)
		temp.Add(temp, z)
		z.Div(temp, big.NewInt(2))
	}

	return y
}

func createPool(tokenAaddress common.Address, tokenBaddress common.Address, fee int64) (*types.Transaction, error) {
	key, err := GetKey(fromAddress.Hex())
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	fromAddress, err = cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	chainId, err := getChainId()
	if err != nil {
		return nil, err
	}
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainId))

	if err != nil {
		return nil, err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit, err = getGasLimit(uint64(6000000))
	if err != nil {
		return nil, err
	}

	txnOpts.Value = big.NewInt(0)

	contract, err := core.NewCore(v3CoreFactoryAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.CreatePool(txnOpts, tokenAaddress, tokenBaddress, big.NewInt(fee))
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to create a v3 pool has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func getPool(tokenAaddress common.Address, tokenBaddress common.Address, fee int64) (*common.Address, error) {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	contract, err := core.NewCore(v3CoreFactoryAddress, client)
	if err != nil {
		return nil, err
	}

	poolAddress, err := contract.GetPool(nil, tokenAaddress, tokenBaddress, big.NewInt(fee))
	if err != nil {
		return nil, err
	}

	fmt.Println("v3 poolAddress", poolAddress)
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return &poolAddress, nil
}

func initializePool(poolAddress common.Address, price int64, tokenAdecimals uint8, tokenBdecimals uint8) (*types.Transaction, error) {
	key, err := GetKey(fromAddress.Hex())
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	fromAddress, err = cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	chainId, err := getChainId()
	if err != nil {
		return nil, err
	}
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainId))

	if err != nil {
		return nil, err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit, err = getGasLimit(uint64(6000000))
	if err != nil {
		return nil, err
	}

	txnOpts.Value = big.NewInt(0)

	contract, err := v3pool.NewV3pool(poolAddress, client)
	if err != nil {
		return nil, err
	}

	sqrtPriceX96 := calculateSqrtPriceX96(big.NewInt(price), tokenAdecimals, tokenBdecimals)

	var tx *types.Transaction
	tx, err = contract.Initialize(txnOpts, sqrtPriceX96)
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to initialize pool v3 has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash(), "sqrtPriceX96", sqrtPriceX96)
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func addLiquidityV3(tokenAaddress common.Address, tokenBaddress common.Address, fee int64, tickLower int64, tickUpper int64,
	amountA int64, amountB int64, amountAmin int64, amountBmin int64) (*types.Transaction, error) {
	key, err := GetKey(fromAddress.Hex())
	if err != nil {
		return nil, err
	}

	fmt.Println("addLiquidityV3", "nonFungiblePositionManagerAddress", nonFungiblePositionManagerAddress, "tokenAaddress", tokenAaddress, "tokenBaddress", tokenBaddress, "fee", fee,
		"tickLower", tickLower, "tickUpper", tickUpper, "amountA", amountA, "amountB", amountB, "amountAmin", amountAmin, "amountBmin", amountBmin)

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	fromAddress, err = cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	chainId, err := getChainId()
	if err != nil {
		return nil, err
	}
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainId))

	if err != nil {
		return nil, err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit, err = getGasLimit(uint64(6000000))
	if err != nil {
		return nil, err
	}

	txnOpts.Value = big.NewInt(0)

	var mintParams nonfungiblepositionmanager.INonfungiblePositionManagerMintParams
	mintParams.Fee = big.NewInt(fee)
	mintParams.TickLower = big.NewInt(tickLower)
	mintParams.TickUpper = big.NewInt(tickUpper)
	mintParams.Recipient = fromAddress //todo: correct check?
	//mintParams.Recipient = nonFungiblePositionManagerAddress //todo: correct check?
	mintParams.Deadline = big.NewInt(9999999999)

	if bytes.Compare(tokenAaddress.Bytes(), tokenBaddress.Bytes()) < 0 {
		fmt.Println("option A")
		mintParams.Token0 = tokenAaddress
		mintParams.Token1 = tokenBaddress
		mintParams.Amount0Desired = params.EtherToWei(big.NewInt(amountA))
		mintParams.Amount1Desired = params.EtherToWei(big.NewInt(amountB))
		mintParams.Amount0Min = params.EtherToWei(big.NewInt(amountAmin))
		mintParams.Amount1Min = params.EtherToWei(big.NewInt(amountBmin))
	} else {
		fmt.Println("option B")
		mintParams.Token0 = tokenBaddress
		mintParams.Token1 = tokenAaddress
		mintParams.Amount0Desired = params.EtherToWei(big.NewInt(amountB))
		mintParams.Amount1Desired = params.EtherToWei(big.NewInt(amountA))
		mintParams.Amount0Min = params.EtherToWei(big.NewInt(amountBmin))
		mintParams.Amount1Min = params.EtherToWei(big.NewInt(amountAmin))
	}

	contract, err := nonfungiblepositionmanager.NewNonfungiblepositionmanager(nonFungiblePositionManagerAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.Mint(txnOpts, mintParams)
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to add liquidity v3 (mint) has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func swapExactInputSingle(tokenInAddress common.Address, tokenOutAddress common.Address, fee int64, amountIn int64, amountOutMinimum int64) (*types.Transaction, error) {
	key, err := GetKey(fromAddress.Hex())
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	fromAddress, err = cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	chainId, err := getChainId()
	if err != nil {
		return nil, err
	}
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainId))

	if err != nil {
		return nil, err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit, err = getGasLimit(uint64(6000000))
	if err != nil {
		return nil, err
	}

	txnOpts.Value = big.NewInt(0)

	var swapParams swaprouter.IV3SwapRouterExactInputSingleParams
	swapParams.TokenIn = tokenInAddress
	swapParams.TokenOut = tokenOutAddress
	swapParams.Fee = big.NewInt(fee)
	swapParams.Recipient = fromAddress //todo: check if correct
	swapParams.AmountIn = params.EtherToWei(big.NewInt(amountIn))
	swapParams.AmountOutMinimum = params.EtherToWei(big.NewInt(amountOutMinimum))
	swapParams.SqrtPriceLimitX96 = big.NewInt(0)

	contract, err := swaprouter.NewSwaprouter(v3SwapRouterContractAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.ExactInputSingle(txnOpts, swapParams)
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to swapExactSingle v3 has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func swapExactOutputSingle(tokenInAddress common.Address, tokenOutAddress common.Address, fee int64, amountOut int64, amountInMaximum int64) (*types.Transaction, error) {
	key, err := GetKey(fromAddress.Hex())
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	fromAddress, err = cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	chainId, err := getChainId()
	if err != nil {
		return nil, err
	}
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(chainId))

	if err != nil {
		return nil, err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit, err = getGasLimit(uint64(6000000))
	if err != nil {
		return nil, err
	}

	txnOpts.Value = big.NewInt(0)

	var swapParams swaprouter.IV3SwapRouterExactOutputSingleParams
	swapParams.TokenIn = tokenInAddress
	swapParams.TokenOut = tokenOutAddress
	swapParams.Fee = big.NewInt(fee)
	swapParams.Recipient = fromAddress //todo: check if correct
	swapParams.AmountOut = params.EtherToWei(big.NewInt(amountOut))
	swapParams.AmountInMaximum = params.EtherToWei(big.NewInt(amountInMaximum))
	swapParams.SqrtPriceLimitX96 = big.NewInt(0)

	contract, err := swaprouter.NewSwaprouter(v3SwapRouterContractAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.ExactOutputSingle(txnOpts, swapParams)
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to swapExactOutputSingle v3 has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}
