package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"os/exec"
	"time"
)

func init() {
	confInit()
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no wallet parameter")
		return
	}

	wallet := os.Args[1]

	client, err := ethclient.Dial(env().Infura)
	if err != nil {
		log.Panic(err)
	}

	d := time.Now().Add(5000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	unlocked, err := keystore.DecryptKey([]byte(env().Key), env().Password)
	if err != nil {
		log.Panic("Bad Password")
	}

	addWallet(client, ctx, unlocked, wallet)
}

func addWallet(client *ethclient.Client, ctx context.Context, unlocked *keystore.Key, wal string) {
	file_abi, err := os.Open(env().Abi)
	if err != nil {
		log.Panic(err)
	}
	defer file_abi.Close()

	coin_abi, err := abi.JSON(file_abi)
	if err != nil {
		log.Panic(err)
	}

	bytesData, err := coin_abi.Pack("addToWhitelist", common.HexToAddress(wal))
	if err != nil {
		log.Panic(err)
	}

	nonce, err := client.NonceAt(ctx, unlocked.Address, nil)

	if err != nil {
		log.Panic(err)
	}

	tx := types.NewTransaction(nonce,
		common.HexToAddress(env().CrowdsaleAddress),
		nil,
		env().Gas,
		big.NewInt(env().GasPrice),
		bytesData,
	)

	signTx, err := types.SignTx(tx, types.HomesteadSigner{}, unlocked.PrivateKey)
	if err != nil {
		log.Panic(err)
	}

	err = client.SendTransaction(ctx, signTx)
	if err != nil {
		log.Panic(err)
	}

	select {
	case <-time.After(5000 * time.Millisecond):
		fmt.Println("time out")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	default:
		runcmd(wal)
	}
}

func runcmd(wal string) {
	cmd := fmt.Sprintf("%s %s", env().Cmd, wal)
	_, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
