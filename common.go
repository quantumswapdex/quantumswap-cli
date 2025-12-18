package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/quantumcoinproject/quantum-coin-go/accounts/keystore"
	"github.com/quantumcoinproject/quantum-coin-go/common"
	"github.com/quantumcoinproject/quantum-coin-go/console/prompt"
	"github.com/quantumcoinproject/quantum-coin-go/crypto/signaturealgorithm"
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
var v2SwapRouterContractAddress common.Address

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

// ParseBigFloat parse string value to big.Float
func ParseBigFloat(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}
