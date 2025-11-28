package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"quantumswap-cli/contracts/core"
	"quantumswap-cli/contracts/nonfungiblepositionmanager"
	"quantumswap-cli/contracts/swaprouter"
	"quantumswap-cli/contracts/v3pool"
	"strconv"
	"strings"
	"time"

	"github.com/quantumcoinproject/quantum-coin-go/accounts/abi/bind"
	"github.com/quantumcoinproject/quantum-coin-go/accounts/keystore"
	"github.com/quantumcoinproject/quantum-coin-go/common"
	"github.com/quantumcoinproject/quantum-coin-go/console/prompt"
	"github.com/quantumcoinproject/quantum-coin-go/core/types"
	"github.com/quantumcoinproject/quantum-coin-go/crypto/cryptobase"
	"github.com/quantumcoinproject/quantum-coin-go/crypto/signaturealgorithm"
	"github.com/quantumcoinproject/quantum-coin-go/ethclient"
	"github.com/quantumcoinproject/quantum-coin-go/params"
	"github.com/quantumcoinproject/quantum-coin-go/token"
)

const (
	ONE_MINUTE_SECONDS = 60
	ONE_HOUR_SECONDS   = ONE_MINUTE_SECONDS * 60
	ONE_DAY_SECONDS    = ONE_HOUR_SECONDS * 24
	ONE_MONTH_SECONDS  = ONE_DAY_SECONDS * 30
	ONE_YEAR_SECONDS   = ONE_DAY_SECONDS * 365
)

// 2592000
const MAX_INCENTIVE_START_LEAD_TIME = ONE_MONTH_SECONDS

// 1892160000
const MAX_INCENTIVE_DURATION = ONE_YEAR_SECONDS * 2

const GAS_LIMIT_ENV = "GAS_LIMIT"
const CHAIN_ID_ENV = "CHAIN_ID"
const DEFAULT_CHAIN_ID = 123123
const NATIVE_CURRENCY_LABEL = "Q"
const ONE_BP_FEE = 100
const ONE_BP_TICK_SPACING = 1

var NATIVE_CURRENCY_LABEL_BYTES = [32]byte(common.BytesToAddress([]byte(NATIVE_CURRENCY_LABEL)))

var fromAddress common.Address
var wqContractAddress common.Address
var v2CoreFactoryAddress common.Address
var v3CoreFactoryAddress common.Address
var multiCallContractAddress common.Address
var proxyAdminContractAddress common.Address
var tickLensContractAddress common.Address
var nftDescriptorLibraryAddress common.Address
var nftPositionDescriptorContractAddress common.Address
var transperentProxyAddress common.Address
var nonFungiblePositionManagerAddress common.Address
var v3MigratorContractAddress common.Address
var v3StakerContractAddress common.Address
var quoterv2ContractAddress common.Address
var v3SwapRouterContractAddress common.Address

func getChainId() (int64, error) {
	chainIdEnv := os.Getenv(CHAIN_ID_ENV)
	if len(chainIdEnv) > 0 {
		chainId, err := strconv.ParseUint(chainIdEnv, 10, 64)
		if err != nil {
			fmt.Println("Error parsing chain id, err")
			return int64(chainId), err
		}
		fmt.Println("Using CHAIN_ID passed using environment variable", chainId)
		return int64(chainId), nil
	} else {
		return DEFAULT_CHAIN_ID, nil
	}
}

func getGasLimit(defaultLimit uint64) (uint64, error) {
	gasLimitEnv := os.Getenv(GAS_LIMIT_ENV)
	if len(gasLimitEnv) > 0 {
		gasLimit, err := strconv.ParseUint(gasLimitEnv, 10, 64)
		if err != nil {
			fmt.Println("Error parsing gas limit", err)
			return gasLimit, err
		}
		fmt.Println("Using gas limit passed using environment variable", gasLimit)
		return gasLimit, nil
	} else {
		return defaultLimit, nil
	}
}

func ReadDataFile(filename string) ([]byte, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	fmt.Println("Successfully Opened ", filename)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue, nil
}

func findKeyFile(keyAddress string) (string, error) {
	keyfile := os.Getenv("DP_KEY_FILE")
	if len(keyfile) > 0 {
		return keyfile, nil
	}

	keyfileDir := os.Getenv("DP_KEY_FILE_DIR")
	if len(keyfileDir) == 0 {
		return "", errors.New("Both DP_KEY_FILE and DP_KEY_FILE_DIR environment variables not set")
	}

	files, err := ioutil.ReadDir(keyfileDir)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	addr := strings.ToLower(strings.Replace(keyAddress, "0x", "", 1))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(strings.ToLower(file.Name()), addr) {
			return filepath.Join(keyfileDir, file.Name()), nil
		}
	}

	return "", errors.New("could not find key file")
}

func GetKeyFromFile(keyFile string, accPwd string) (*signaturealgorithm.PrivateKey, error) {
	secretKey, err := ReadDataFile(keyFile)
	if err != nil {
		return nil, err
	}

	password := accPwd
	key, err := keystore.DecryptKey(secretKey, password)
	if err != nil {
		return nil, err
	}

	return key.PrivateKey, nil
}

func GetKey(address string) (*signaturealgorithm.PrivateKey, error) {
	keyFile, err := findKeyFile(address)
	if err != nil {
		return nil, err
	}

	fmt.Println("keyFile", keyFile)
	secretKey, err := ReadDataFile(keyFile)
	if err != nil {
		return nil, err
	}
	password := os.Getenv("DP_ACC_PWD")
	if len(password) == 0 {
		password, err = prompt.Stdin.PromptPassword(fmt.Sprintf("Enter the wallet password : "))
		if err != nil {
			return nil, err
		}
		if len(password) == 0 {
			return nil, errors.New("password is not correct")
		}
	}
	key, err := keystore.DecryptKey(secretKey, password)
	if err != nil {
		return nil, err
	}

	if key.Address.IsEqualTo(common.HexToAddress(address)) == false {
		return nil, errors.New("address mismatch in key file")
	}

	return key.PrivateKey, nil
}

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

	fmt.Println("Your request to create a pool has been added to the queue for processing. Please check your account after 10 minutes.")
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

	fmt.Println("poolAddress", poolAddress)
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

	fmt.Println("Your request to initialize pool has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash(), "sqrtPriceX96", sqrtPriceX96)
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func approve(tokenAddress common.Address, approveAddress common.Address, amount int64) (*types.Transaction, error) {
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

	contract, err := token.NewToken(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.Approve(txnOpts, approveAddress, params.EtherToWei(big.NewInt(amount)))
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to approve spend has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func addLiquidity(tokenAaddress common.Address, tokenBaddress common.Address, fee int64, tickLower int64, tickUpper int64,
	amountA int64, amountB int64, amountAmin int64, amountBmin int64) (*types.Transaction, error) {
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

	var mintParams nonfungiblepositionmanager.INonfungiblePositionManagerMintParams
	mintParams.Fee = big.NewInt(fee)
	mintParams.TickLower = big.NewInt(tickLower)
	mintParams.TickUpper = big.NewInt(tickUpper)
	mintParams.Recipient = fromAddress //todo: correct check?
	mintParams.Deadline = big.NewInt(999999999999999999)

	if bytes.Compare(tokenAaddress.Bytes(), tokenBaddress.Bytes()) < 0 {
		mintParams.Token0 = tokenAaddress
		mintParams.Token1 = tokenBaddress
		mintParams.Amount0Desired = params.EtherToWei(big.NewInt(amountA))
		mintParams.Amount1Desired = params.EtherToWei(big.NewInt(amountB))
		mintParams.Amount0Min = params.EtherToWei(big.NewInt(amountAmin))
		mintParams.Amount1Min = params.EtherToWei(big.NewInt(amountBmin))
	} else {
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

	fmt.Println("Your request to add liquidity (mint) has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func swapExactSingle(tokenInAddress common.Address, tokenOutAddress common.Address, fee int64, amountIn int64, amountOutMinimum int64) (*types.Transaction, error) {
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
	swapParams.AmountIn = big.NewInt(amountIn)
	swapParams.AmountOutMinimum = big.NewInt(amountOutMinimum)
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

	fmt.Println("Your request to swapExactSingle has been added to the queue for processing. Please check your account after 10 minutes.")
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
	swapParams.AmountOut = big.NewInt(amountOut)
	swapParams.AmountInMaximum = big.NewInt(amountInMaximum)

	contract, err := swaprouter.NewSwaprouter(v3SwapRouterContractAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.ExactOutputSingle(txnOpts, swapParams)
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to swapExactOutputSingle has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

// ParseBigFloat parse string value to big.Float
func ParseBigFloat(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}
