package main

import (
	"fmt"
	"github.com/DSiSc/wasm-cdt/contract"
	"gopkg.in/urfave/cli.v1"
	"math/big"
)

var (
	DeployCommand = cli.Command{
		Action: DeployContract,
		Name:   "deploy",
		Usage:  "deploy wasm contract",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "wasm contract `FILE`",
			},
			cli.StringFlag{
				Name:  "account, ac",
				Usage: "account used to deploy contract",
			},
			cli.StringFlag{
				Name:  "endpoint, ep",
				Usage: "endpoint used to deploy contract",
			},
			cli.StringFlag{
				Name:  "network, net",
				Usage: "network id(main net:1, local test: >1, default: 2)",
				Value: "2",
			},
			cli.StringFlag{
				Name:  "value, v",
				Usage: "value to transfer to contract",
				Value: "0",
			},
			DataDirFlag,
			KeyStoreDirFlag,
		}}
)

func DeployContract(ctx *cli.Context) error {
	sourceFile := ctx.String("file")
	account := ctx.String("account")
	endpoint := ctx.String("endpoint")
	network := ctx.String("network")
	value := ctx.String("value")

	networkId, ok := new(big.Int).SetString(network, 0)
	if !ok {
		return fmt.Errorf("invalid network id %s", network)
	}

	amount, ok := new(big.Int).SetString(value, 0)
	if !ok {
		return fmt.Errorf("invalid transfer value")
	}

	return contract.DelpoyContract(sourceFile, account, GetKeystoreDir(ctx), endpoint, networkId, amount)
}
