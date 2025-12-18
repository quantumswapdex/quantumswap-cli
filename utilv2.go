package main

import (
	"context"
	"fmt"
	"math/big"
	"quantumswap-cli/contracts/corev2"
	"quantumswap-cli/contracts/v2swaprouter"
	"time"

	"github.com/quantumcoinproject/quantum-coin-go/accounts/abi/bind"
	"github.com/quantumcoinproject/quantum-coin-go/common"
	"github.com/quantumcoinproject/quantum-coin-go/core/types"
	"github.com/quantumcoinproject/quantum-coin-go/crypto/cryptobase"
	"github.com/quantumcoinproject/quantum-coin-go/ethclient"
	"github.com/quantumcoinproject/quantum-coin-go/params"
)

func createPair(tokenAaddress common.Address, tokenBaddress common.Address) (*types.Transaction, error) {
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

	contract, err := corev2.NewCorev2(v2CoreFactoryAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.CreatePair(txnOpts, tokenAaddress, tokenBaddress)
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to create a v2 pair has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func getPair(tokenAaddress common.Address, tokenBaddress common.Address) (*common.Address, error) {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	contract, err := corev2.NewCorev2(v2CoreFactoryAddress, client)
	if err != nil {
		return nil, err
	}

	pairAddress, err := contract.GetPair(nil, tokenAaddress, tokenBaddress)
	if err != nil {
		return nil, err
	}

	fmt.Println("v2 pairAddress", pairAddress)
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return &pairAddress, nil
}

func addLiquidityV2(tokenAaddress common.Address, tokenBaddress common.Address,
	amountA int64, amountB int64, amountAmin int64, amountBmin int64) (*types.Transaction, error) {
	key, err := GetKey(fromAddress.Hex())
	if err != nil {
		return nil, err
	}

	fmt.Println("addLiquidityV2", "v2SwapRouterContractAddress", v2SwapRouterContractAddress, "tokenAaddress", tokenAaddress, "tokenBaddress", tokenBaddress,
		"amountA", amountA, "amountB", amountB, "amountAmin", amountAmin, "amountBmin", amountBmin)

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

	contract, err := v2swaprouter.NewV2swaprouter(v2SwapRouterContractAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.AddLiquidity(txnOpts, tokenAaddress, tokenBaddress, params.EtherToWei(big.NewInt(amountA)), params.EtherToWei(big.NewInt(amountB)),
		params.EtherToWei(big.NewInt(amountAmin)), params.EtherToWei(big.NewInt(amountBmin)), fromAddress, big.NewInt(9999999999))
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to add liquidity v2 (mint) has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}

func swapExactTokensForTokens(tokenInAddress common.Address, tokenOutAddress common.Address, amountIn int64, amountOutMinimum int64) (*types.Transaction, error) {
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

	contract, err := v2swaprouter.NewV2swaprouter(v2SwapRouterContractAddress, client)
	if err != nil {
		return nil, err
	}

	var tx *types.Transaction
	tx, err = contract.SwapExactTokensForTokens(txnOpts, params.EtherToWei(big.NewInt(amountIn)),
		params.EtherToWei(big.NewInt(amountOutMinimum)), []common.Address{tokenInAddress, tokenOutAddress}, fromAddress, big.NewInt(9999999999))
	if err != nil {
		return nil, err
	}

	fmt.Println("Your request to swapExactTokensForTokens v2 has been added to the queue for processing. Please check your account after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return tx, nil
}
