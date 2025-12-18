package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/quantumcoinproject/quantum-coin-go/common"
	"github.com/quantumcoinproject/quantum-coin-go/console/prompt"
)

var rawURL string

func printHelp() {
	fmt.Println("--------")
	fmt.Println(" Usage")
	fmt.Println("--------")

	fmt.Println("(optional) quantumswap-cli createpair TOKEN_A_ADDRESS TOKEN_B_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("           V2_CORE_FACTORY_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-cli getpair TOKEN_A_ADDRESS TOKEN_B_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("           V2_CORE_FACTORY_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-deploy addliquidityv2 TOKEN_A_ADDRESS TOKEN_B_ADDRESS AMOUNT_A AMOUNT_B AMOUNT_A_MIN AMOUNT_B_MIN")
	fmt.Println(" !!!LIMITATION!!! Amount will be converted to wei based on 18 decimals internally. Other decimals not supported.")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("            SWAP_ROUTER_V2_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-deploy swapexacttokensFortokens TOKEN_IN_ADDRESS TOKEN_OUT_ADDRESS AMOUNT_IN AMOUNT_OUT_MIN")
	fmt.Println(" FEE should be 500 or 3000 or 10000 (For 0.3, 0.05%, 0.3%, or 1%)")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("            SWAP_ROUTER_V2_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-cli createpool TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE")
	fmt.Println(" FEE should be 500 or 3000 or 10000 (For 0.3, 0.05%, 0.3%, or 1%)")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("           V3_CORE_FACTORY_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-cli getpool TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE")
	fmt.Println(" FEE should be 500 or 3000 or 10000 (For 0.3, 0.05%, 0.3%, or 1%)")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("           V3_CORE_FACTORY_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-cli initializepool POOL_ADDRESS PRICE_IN_TOKEN_B_PER_TOKEN_A TOKEN_A_DECIMALS TOKEN_B_DECIMALS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")

	fmt.Println("(optional) quantumswap-deploy addliquidityv3 TOKEN_A_ADDRESS TOKEN_B_ADDRESS FEE TICK_LOWER TICK_UPPER AMOUNT_A AMOUNT_B AMOUNT_A_MIN AMOUNT_B_MIN")
	fmt.Println(" !!!LIMITATION!!! Amount will be converted to wei based on 18 decimals internally. Other decimals not supported.")
	fmt.Println(" FEE should be 500 or 3000 or 10000 (For 0.3, 0.05%, 0.3%, or 1%)")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("           NONFUNGIBLE_POSITION_MANAGER_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-deploy exactinputsingle TOKEN_IN_ADDRESS TOKEN_OUT_ADDRESS FEE AMOUNT_IN AMOUNT_OUT_MIN")
	fmt.Println(" FEE should be 500 or 3000 or 10000 (For 0.3, 0.05%, 0.3%, or 1%)")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("           SWAP_ROUTER_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-deploy exactoutputsingle TOKEN_IN_ADDRESS TOKEN_OUT_ADDRESS FEE AMOUNT_OUT AMOUNT_IN_MAX")
	fmt.Println(" FEE should be 500 or 3000 or 10000 (For 0.3, 0.05%, 0.3%, or 1%)")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           CHAIN_ID, DP_RAW_URL, DP_KEY_FILE_DIR or DP_KEY_FILE,GAS_LIMIT,FROM_ADDRESS")
	fmt.Println("      Set the following additional environment variables:")
	fmt.Println("           SWAP_ROUTER_CONTRACT_ADDRESS")

	fmt.Println("(optional) quantumswap-deploy ticktoprice TICK")

	fmt.Println("(optional) quantumswap-deploy pricetotick PRICE")
}

func main() {
	fmt.Println("===================")
	fmt.Println(" QuantumSwap CLI")
	fmt.Println("===================")

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	rawURL = os.Getenv("DP_RAW_URL")
	if len(rawURL) == 0 {
		runtimeOS := strings.ToLower(runtime.GOOS)
		if runtimeOS == "windows" {
			rawURL = "\\\\.\\pipe\\geth.ipc"
		} else {
			rawURL = "data/geth.ipc"
		}
	}

	if os.Args[1] == "createpair" {
		CreatePair()
	} else if os.Args[1] == "getpair" {
		GetPair()
	} else if os.Args[1] == "addliquidityv2" {
		AddLiquidityV2()
	} else if os.Args[1] == "swapexacttokensFortokens" {
		SwapExactTokensForTokens()
	} else if os.Args[1] == "createpool" {
		CreatePool()
	} else if os.Args[1] == "getpool" {
		GetPool()
	} else if os.Args[1] == "initializepool" {
		InitializePool()
	} else if os.Args[1] == "addliquidityv3" {
		AddLiquidityV3()
	} else if os.Args[1] == "exactinputsingle" {
		ExactInputSingle()
	} else if os.Args[1] == "exactoutputsingle" {
		ExactOutputSingle()
	} else if os.Args[1] == "ticktoprice" {
		TickToPrice()
	} else if os.Args[1] == "pricetotick" {
		PriceToTick()
	} else {
		printHelp()
	}
}

func CreatePair() {
	if len(os.Args) < 4 {
		printHelp()
		return
	}

	tokenAaddr := os.Args[2]
	if common.IsHexAddress(tokenAaddr) == false {
		fmt.Println("Invalid TOKEN_A_ADDRESS", tokenAaddr)
		return
	}
	tokenAaddress := common.HexToAddress(tokenAaddr)

	tokenBaddr := os.Args[3]
	if common.IsHexAddress(tokenBaddr) == false {
		fmt.Println("Invalid TOKEN_B_ADDRESS", tokenBaddr)
		return
	}
	tokenBaddress := common.HexToAddress(tokenBaddr)

	v2coreFactoryContractAddr := os.Getenv("V2_CORE_FACTORY_CONTRACT_ADDRESS")
	if common.IsHexAddress(v2coreFactoryContractAddr) == false {
		fmt.Println("Invalid V2_CORE_FACTORY_CONTRACT_ADDRESS", v2coreFactoryContractAddr)
		return
	}
	v2CoreFactoryAddress = common.HexToAddress(v2coreFactoryContractAddr)

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to CreatePair from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = createPair(tokenAaddress, tokenBaddress)
	if err != nil {
		fmt.Println("createPair error", err)
		return
	}
}

func GetPair() {
	if len(os.Args) < 4 {
		printHelp()
		return
	}

	tokenAaddr := os.Args[2]
	if common.IsHexAddress(tokenAaddr) == false {
		fmt.Println("Invalid TOKEN_A_ADDRESS", tokenAaddr)
		return
	}
	tokenAaddress := common.HexToAddress(tokenAaddr)

	tokenBaddr := os.Args[3]
	if common.IsHexAddress(tokenBaddr) == false {
		fmt.Println("Invalid TOKEN_B_ADDRESS", tokenBaddr)
		return
	}
	tokenBaddress := common.HexToAddress(tokenBaddr)

	v2coreFactoryContractAddr := os.Getenv("V2_CORE_FACTORY_CONTRACT_ADDRESS")
	if common.IsHexAddress(v2coreFactoryContractAddr) == false {
		fmt.Println("Invalid V2_CORE_FACTORY_CONTRACT_ADDRESS", v2coreFactoryContractAddr)
		return
	}
	v2CoreFactoryAddress = common.HexToAddress(v2coreFactoryContractAddr)

	_, err := getPair(tokenAaddress, tokenBaddress)
	if err != nil {
		fmt.Println("getPair error", err)
		return
	}
}

func AddLiquidityV2() {
	if len(os.Args) < 8 {
		printHelp()
		return
	}

	tokenAaddr := os.Args[2]
	if common.IsHexAddress(tokenAaddr) == false {
		fmt.Println("Invalid TOKEN_A_ADDRESS", tokenAaddr)
		return
	}
	tokenAaddress := common.HexToAddress(tokenAaddr)

	tokenBaddr := os.Args[3]
	if common.IsHexAddress(tokenBaddr) == false {
		fmt.Println("Invalid TOKEN_B_ADDRESS", tokenBaddr)
		return
	}
	tokenBaddress := common.HexToAddress(tokenBaddr)

	amountAval := os.Args[4]
	amountA, err := strconv.ParseUint(amountAval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_A", err)
		return
	}

	amountBval := os.Args[5]
	amountB, err := strconv.ParseUint(amountBval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_B", err)
		return
	}

	amountAminval := os.Args[6]
	amountAmin, err := strconv.ParseUint(amountAminval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_A_MIN", err)
		return
	}

	amountBminval := os.Args[7]
	amountBmin, err := strconv.ParseUint(amountBminval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_B_MIN", err)
		return
	}

	v2SwapRouterContractAddr := os.Getenv("SWAP_ROUTER_V2_CONTRACT_ADDRESS")
	if common.IsHexAddress(v2SwapRouterContractAddr) == false {
		fmt.Println("Invalid  SWAP_ROUTER_V2_CONTRACT_ADDRESS", v2SwapRouterContractAddr)
		return
	}
	v2SwapRouterContractAddress = common.HexToAddress(v2SwapRouterContractAddr)

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	fmt.Println("addLiquidityV2", "v2SwapRouterContractAddr", v2SwapRouterContractAddr, "tokenAaddress", tokenAaddress, "tokenBaddress", tokenBaddress,
		"amountA", amountA, "amountB", amountB, "amountAmin", amountAmin, "amountBmin", amountBmin)

	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to AddLiquidityV2 from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = addLiquidityV2(tokenAaddress, tokenBaddress, int64(amountA), int64(amountB), int64(amountAmin), int64(amountBmin))
	if err != nil {
		fmt.Println("addLiquidityV2 error", err)
		return
	}
}

func SwapExactTokensForTokens() {
	if len(os.Args) < 6 {
		printHelp()
		return
	}

	tokenInaddr := os.Args[2]
	if common.IsHexAddress(tokenInaddr) == false {
		fmt.Println("Invalid TOKEN_IN_ADDRESS", tokenInaddr)
		return
	}
	tokenInAddress := common.HexToAddress(tokenInaddr)

	tokenOutaddr := os.Args[3]
	if common.IsHexAddress(tokenOutaddr) == false {
		fmt.Println("Invalid TOKEN_OUT_ADDRESS", tokenOutaddr)
		return
	}
	tokenOutAddress := common.HexToAddress(tokenOutaddr)

	amountInVal := os.Args[4]
	amountIn, err := strconv.ParseUint(amountInVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_IN", err)
		return
	}

	amountOutMinVal := os.Args[5]
	amountOutMin, err := strconv.ParseUint(amountOutMinVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_OUT_MIN", err)
		return
	}

	v2SwapRouterContractAddr := os.Getenv("SWAP_ROUTER_V2_CONTRACT_ADDRESS")
	if common.IsHexAddress(v2SwapRouterContractAddr) == false {
		fmt.Println("Invalid  SWAP_ROUTER_V2_CONTRACT_ADDRESS", v2SwapRouterContractAddr)
		return
	}
	v2SwapRouterContractAddress = common.HexToAddress(v2SwapRouterContractAddr)

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	fmt.Println("SwapExactTokensForTokens", "v2SwapRouterContractAddr", v2SwapRouterContractAddr, "tokenInaddr", tokenInaddr, "tokenOutAddress", tokenOutAddress,
		"amountInVal", amountInVal, "amountOutMinVal", amountOutMinVal)

	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to SwapExactSingle from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = swapExactTokensForTokens(tokenInAddress, tokenOutAddress, int64(amountIn), int64(amountOutMin))
	if err != nil {
		fmt.Println("swapExactTokensForTokens error", err)
		return
	}
}

func CreatePool() {
	if len(os.Args) < 5 {
		printHelp()
		return
	}

	tokenAaddr := os.Args[2]
	if common.IsHexAddress(tokenAaddr) == false {
		fmt.Println("Invalid TOKEN_A_ADDRESS", tokenAaddr)
		return
	}
	tokenAaddress := common.HexToAddress(tokenAaddr)

	tokenBaddr := os.Args[3]
	if common.IsHexAddress(tokenBaddr) == false {
		fmt.Println("Invalid TOKEN_B_ADDRESS", tokenBaddr)
		return
	}
	tokenBaddress := common.HexToAddress(tokenBaddr)

	feeVal := os.Args[4]
	fee, err := strconv.ParseUint(feeVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing FEE", err)
		return
	}
	if fee != 500 && fee != 3000 && fee != 10000 {
		fmt.Println("Accepted values for FEE are 500, 3000, 10000")
		return
	}

	v3coreFactoryContractAddr := os.Getenv("V3_CORE_FACTORY_CONTRACT_ADDRESS")
	if common.IsHexAddress(v3coreFactoryContractAddr) == false {
		fmt.Println("Invalid V3_CORE_FACTORY_CONTRACT_ADDRESS", v3coreFactoryContractAddr)
		return
	}
	v3CoreFactoryAddress = common.HexToAddress(v3coreFactoryContractAddr)

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to CreatePool from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = createPool(tokenAaddress, tokenBaddress, int64(fee))
	if err != nil {
		fmt.Println("createPool error", err)
		return
	}
}

func GetPool() {
	if len(os.Args) < 5 {
		printHelp()
		return
	}

	tokenAaddr := os.Args[2]
	if common.IsHexAddress(tokenAaddr) == false {
		fmt.Println("Invalid TOKEN_A_ADDRESS", tokenAaddr)
		return
	}
	tokenAaddress := common.HexToAddress(tokenAaddr)

	tokenBaddr := os.Args[3]
	if common.IsHexAddress(tokenBaddr) == false {
		fmt.Println("Invalid TOKEN_B_ADDRESS", tokenBaddr)
		return
	}
	tokenBaddress := common.HexToAddress(tokenBaddr)

	feeVal := os.Args[4]
	fee, err := strconv.ParseUint(feeVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing FEE", err)
		return
	}
	if fee != 500 && fee != 3000 && fee != 10000 {
		fmt.Println("Accepted values for FEE are 500, 3000, 10000")
		return
	}

	v3coreFactoryContractAddr := os.Getenv("V3_CORE_FACTORY_CONTRACT_ADDRESS")
	if common.IsHexAddress(v3coreFactoryContractAddr) == false {
		fmt.Println("Invalid V3_CORE_FACTORY_CONTRACT_ADDRESS", v3coreFactoryContractAddr)
		return
	}
	v3CoreFactoryAddress = common.HexToAddress(v3coreFactoryContractAddr)

	_, err = getPool(tokenAaddress, tokenBaddress, int64(fee))
	if err != nil {
		fmt.Println("getPool error", err)
		return
	}
}

func InitializePool() {
	if len(os.Args) < 6 {
		printHelp()
		return
	}

	poolAddr := os.Args[2]
	if common.IsHexAddress(poolAddr) == false {
		fmt.Println("Invalid POOL_ADDRESS", poolAddr)
		return
	}
	poolAddress := common.HexToAddress(poolAddr)

	priceVal := os.Args[3]
	price, err := strconv.ParseUint(priceVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing PRICE_IN_TOKEN_B_PER_TOKEN_A", err)
		return
	}

	tokenAdecimalsVal := os.Args[4]
	tokenAdecimals, err := strconv.ParseUint(tokenAdecimalsVal, 10, 8)
	if err != nil {
		fmt.Println("Error parsing TOKEN_A_DECIMALS", err)
		return
	}
	if tokenAdecimals > 18 {
		fmt.Println("Invalid TOKEN_A_DECIMALS", err)
		return
	}

	tokenBdecimalsVal := os.Args[5]
	tokenBdecimals, err := strconv.ParseUint(tokenBdecimalsVal, 10, 8)
	if err != nil {
		fmt.Println("Error parsing TOKEN_B_DECIMALS", err)
		return
	}
	if tokenBdecimals > 18 {
		fmt.Println("Invalid TOKEN_A_DECIMALS", err)
		return
	}

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to InitializePool from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = initializePool(poolAddress, int64(price), uint8(tokenAdecimals), uint8(tokenBdecimals))
	if err != nil {
		fmt.Println("initializePool error", err)
		return
	}
}

func AddLiquidityV3() {
	if len(os.Args) < 11 {
		printHelp()
		return
	}

	tokenAaddr := os.Args[2]
	if common.IsHexAddress(tokenAaddr) == false {
		fmt.Println("Invalid TOKEN_A_ADDRESS", tokenAaddr)
		return
	}
	tokenAaddress := common.HexToAddress(tokenAaddr)

	tokenBaddr := os.Args[3]
	if common.IsHexAddress(tokenBaddr) == false {
		fmt.Println("Invalid TOKEN_B_ADDRESS", tokenBaddr)
		return
	}
	tokenBaddress := common.HexToAddress(tokenBaddr)

	feeVal := os.Args[4]
	fee, err := strconv.ParseUint(feeVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing FEE", err)
		return
	}
	if fee != 500 && fee != 3000 && fee != 10000 {
		fmt.Println("Accepted values for FEE are 500, 3000, 10000")
		return
	}

	tickLowerVal := os.Args[5]
	tickLower, err := strconv.ParseInt(tickLowerVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing TICK_LOWER", err)
		return
	}

	tickUpperVal := os.Args[6]
	tickUpper, err := strconv.ParseInt(tickUpperVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing TICK_UPPER", err)
		return
	}

	amountAval := os.Args[7]
	amountA, err := strconv.ParseUint(amountAval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_A", err)
		return
	}

	amountBval := os.Args[8]
	amountB, err := strconv.ParseUint(amountBval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_B", err)
		return
	}

	amountAminval := os.Args[9]
	amountAmin, err := strconv.ParseUint(amountAminval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_A_MIN", err)
		return
	}

	amountBminval := os.Args[10]
	amountBmin, err := strconv.ParseUint(amountBminval, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_B_MIN", err)
		return
	}

	nfPositionManagerAddr := os.Getenv("NONFUNGIBLE_POSITION_MANAGER_CONTRACT_ADDRESS")
	if common.IsHexAddress(nfPositionManagerAddr) == false {
		fmt.Println("Invalid NONFUNGIBLE_POSITION_MANAGER_CONTRACT_ADDRESS", nfPositionManagerAddr)
		return
	}
	nonFungiblePositionManagerAddress = common.HexToAddress(nfPositionManagerAddr)

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to AddLiquidityV3 from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = addLiquidityV3(tokenAaddress, tokenBaddress, int64(fee), tickLower, tickUpper, int64(amountA), int64(amountB), int64(amountAmin), int64(amountBmin))
	if err != nil {
		fmt.Println("addLiquidityV3 error", err)
		return
	}
}

func ExactInputSingle() {
	if len(os.Args) < 7 {
		printHelp()
		return
	}

	tokenInaddr := os.Args[2]
	if common.IsHexAddress(tokenInaddr) == false {
		fmt.Println("Invalid TOKEN_IN_ADDRESS", tokenInaddr)
		return
	}
	tokenInAddress := common.HexToAddress(tokenInaddr)

	tokenOutaddr := os.Args[3]
	if common.IsHexAddress(tokenOutaddr) == false {
		fmt.Println("Invalid TOKEN_OUT_ADDRESS", tokenOutaddr)
		return
	}
	tokenOutAddress := common.HexToAddress(tokenOutaddr)

	feeVal := os.Args[4]
	fee, err := strconv.ParseUint(feeVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing FEE", err)
		return
	}
	if fee != 500 && fee != 3000 && fee != 10000 {
		fmt.Println("Accepted values for FEE are 500, 3000, 10000")
		return
	}

	amountInVal := os.Args[5]
	amountIn, err := strconv.ParseUint(amountInVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_IN", err)
		return
	}

	amountOutMinVal := os.Args[6]
	amountOutMin, err := strconv.ParseUint(amountOutMinVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_OUT_MIN", err)
		return
	}

	v3SwapRouterContractAddr := os.Getenv("SWAP_ROUTER_CONTRACT_ADDRESS")
	if common.IsHexAddress(v3SwapRouterContractAddr) == false {
		fmt.Println("Invalid SWAP_ROUTER_CONTRACT_ADDRESS", v3SwapRouterContractAddr)
		return
	}
	v3SwapRouterContractAddress = common.HexToAddress(v3SwapRouterContractAddr)

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	fmt.Println("SwapExactSingle", "v3SwapRouterContractAddr", v3SwapRouterContractAddr, "tokenInaddr", tokenInaddr, "tokenOutAddress", tokenOutAddress, "fee", fee,
		"amountInVal", amountInVal, "amountOutMinVal", amountOutMinVal)
	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to SwapExactSingle from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = swapExactInputSingle(tokenInAddress, tokenOutAddress, int64(fee), int64(amountIn), int64(amountOutMin))
	if err != nil {
		fmt.Println("swapExactSingle error", err)
		return
	}
}

func ExactOutputSingle() {
	if len(os.Args) < 7 {
		printHelp()
		return
	}

	tokenInaddr := os.Args[2]
	if common.IsHexAddress(tokenInaddr) == false {
		fmt.Println("Invalid TOKEN_IN_ADDRESS", tokenInaddr)
		return
	}
	tokenInAddress := common.HexToAddress(tokenInaddr)

	tokenOutaddr := os.Args[3]
	if common.IsHexAddress(tokenOutaddr) == false {
		fmt.Println("Invalid TOKEN_OUT_ADDRESS", tokenOutaddr)
		return
	}
	tokenOutAddress := common.HexToAddress(tokenOutaddr)

	feeVal := os.Args[4]
	fee, err := strconv.ParseUint(feeVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing FEE", err)
		return
	}
	if fee != 500 && fee != 3000 && fee != 10000 {
		fmt.Println("Accepted values for FEE are 500, 3000, 10000")
		return
	}

	amountOutVal := os.Args[5]
	amountOut, err := strconv.ParseUint(amountOutVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_OUT", err)
		return
	}

	amountInMaxVal := os.Args[6]
	amountInMax, err := strconv.ParseUint(amountInMaxVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing AMOUNT_IN_MAX", err)
		return
	}

	v3SwapRouterContractAddr := os.Getenv("SWAP_ROUTER_CONTRACT_ADDRESS")
	if common.IsHexAddress(v3SwapRouterContractAddr) == false {
		fmt.Println("Invalid SWAP_ROUTER_CONTRACT_ADDRESS", v3SwapRouterContractAddr)
		return
	}
	v3SwapRouterContractAddress = common.HexToAddress(v3SwapRouterContractAddr)

	fromAddr := os.Getenv("FROM_ADDRESS")
	if common.IsHexAddress(fromAddr) == false {
		fmt.Println("Invalid FROM_ADDRESS", fromAddr)
		return
	}
	fromAddress = common.HexToAddress(fromAddr)

	fmt.Println("ExactOutputSingle", "v3SwapRouterContractAddr", v3SwapRouterContractAddr, "tokenInaddr", tokenInaddr, "tokenOutAddress", tokenOutAddress, "fee", fee,
		"amountOutVal", amountOutVal, "amountInMaxVal", amountInMaxVal)
	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you want to SwapExactSingle from %s?", fromAddress))
	if err != nil {
		fmt.Println("error", err)
		return
	}
	if ethConfirm != true {
		fmt.Println("confirmation not made")
		return
	}

	_, err = swapExactOutputSingle(tokenInAddress, tokenOutAddress, int64(fee), int64(amountOut), int64(amountInMax))
	if err != nil {
		fmt.Println("swapExactSingle error", err)
		return
	}
}

func PriceToTick() {
	if len(os.Args) < 3 {
		printHelp()
		return
	}

	priceVal := os.Args[2]
	price, err := ParseBigFloat(priceVal)
	if err != nil {
		fmt.Println("Error parsing PRICE", err)
		return
	}

	tick := getTickFromPrice(price)
	fmt.Println("Price", price, "Tick", tick)
}

func TickToPrice() {
	if len(os.Args) < 3 {
		printHelp()
		return
	}

	tickVal := os.Args[2]
	tick, err := strconv.ParseInt(tickVal, 10, 64)
	if err != nil {
		fmt.Println("Error parsing TICK", err)
		return
	}

	price := getPriceFromTick(int32(tick))
	fmt.Println("Tick", tick, "Price", price)
}
