package main

import (
	"fmt"
	"github.com/DSiSc/crypto-suite/crypto"
	"github.com/DSiSc/wallet/accounts"
	"github.com/DSiSc/wallet/accounts/keystore"
	"github.com/DSiSc/wallet/utils"
	common2 "github.com/DSiSc/wasm-cdt/common"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
)

var (
	DataDirFlagName  = "datadir"
	KeyStoreName     = "keystore"
	PasswordFlagName = "password"

	// General settings
	DataDirFlag = cli.StringFlag{
		Name:  "datadir, d",
		Usage: "Data directory for the databases and keystore",
		Value: filepath.Join(`/workspace/`, ".wallet"),
	}

	KeyStoreDirFlag = cli.StringFlag{
		Name:  "keystore, ks",
		Usage: "Directory for the keystore (default = inside the datadir)",
		Value: keystore.KeyStoreScheme,
	}

	PasswordFileFlag = cli.StringFlag{
		Name:  "password, p",
		Usage: "Password file to use for non-interactive password input",
		Value: "",
	}

	AccountCommand = cli.Command{
		Name:  "account",
		Usage: "Manage accounts",
		Description: `Manage accounts, list all existing accounts, import a private key into a new
account, create a new account or update an existing account.`,
		Subcommands: []cli.Command{
			{
				Name:   "list",
				Usage:  "Print summary of existing accounts",
				Action: ListAccounts,
				Flags: []cli.Flag{
					DataDirFlag,
					KeyStoreDirFlag,
				},
				Description: `Print a short summary of all accounts`,
			},
			{
				Name:   "new",
				Usage:  "Create a new account",
				Action: CreateAccount,
				Flags: []cli.Flag{
					DataDirFlag,
					KeyStoreDirFlag,
					PasswordFileFlag,
				},
				Description: `geth account new`,
			},
			{
				Name:   "import",
				Usage:  "Import a private key into a new account",
				Action: ImportAccount,
				Flags: []cli.Flag{
					DataDirFlag,
					KeyStoreDirFlag,
					PasswordFileFlag,
				},
				ArgsUsage:   "<keyFile>",
				Description: `Imports an unencrypted private key from <keyfile> and creates a new account`,
			},
		},
	}
)

//AccountList list all account in local keystore
func ListAccounts(ctx *cli.Context) error {
	keyStoreDir := GetKeystoreDir(ctx)

	accountList := GetAllAccount(keyStoreDir)
	for index, account := range accountList {
		fmt.Printf("Account #%d: {%x} %s\n", index, account.Address, &account.URL)
	}

	return nil
}

//AccountList list all account in local keystore
func GetAllAccount(keyStoreDir string) []accounts.Account {
	accountList := make([]accounts.Account, 0)
	manager, _, err := utils.MakeAccountManager(keyStoreDir)
	if err != nil {
		utils.Fatalf("Could not make account manager: %v", err)
	}

	for _, wallet := range manager.Wallets() {
		for _, account := range wallet.Accounts() {
			accountList = append(accountList, account)
		}
	}

	return accountList
}

// CreateAccount creates a new account into the keystore defined by the CLI flags.
func CreateAccount(ctx *cli.Context) error {
	keyStoreDir := GetKeystoreDir(ctx)
	password := common2.GetPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true)
	_, err := utils.NewAccount(keyStoreDir, password)
	return err
}

// Import an account from private key
func ImportAccount(ctx *cli.Context) error {
	keyfile := ctx.Args().First()
	if len(keyfile) == 0 {
		utils.Fatalf("keyfile must be given as argument")
	}
	key, err := crypto.LoadECDSA(keyfile)
	if err != nil {
		utils.Fatalf("Failed to load the private key: %v", err)
	}

	keyStoreDir := GetKeystoreDir(ctx)

	manager, _, _ := utils.MakeAccountManager(keyStoreDir)
	ks := manager.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	passphrase := common2.GetPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true)
	acct, err := ks.ImportECDSA(key, passphrase)
	if err != nil {
		utils.Fatalf("Could not create the account: %v", err)
	}
	fmt.Printf("Address: {%x}\n", acct.Address)
	return nil
}

func getDataDir(ctx *cli.Context) string {
	return ctx.String(DataDirFlagName)
}

func GetKeystoreDir(ctx *cli.Context) string {
	dataDir := getDataDir(ctx)
	keystoreName := ctx.String(KeyStoreName)
	return filepath.Join(dataDir, keystoreName)
}

func getPassword(ctx *cli.Context) string {
	return ctx.String(PasswordFlagName)
}
