package contract

import (
	"fmt"
	"github.com/DSiSc/craft/rlp"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/util"
	"github.com/DSiSc/wallet/accounts"
	"github.com/DSiSc/wallet/utils"
	"github.com/DSiSc/wasm-cdt/common"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"time"
)

func DelpoyContract(contractFile, account, keystoreDir, endpoint string, networkId *big.Int, value *big.Int) error {
	nonce, err := common.GetAccountNonde(account, endpoint)
	if err != nil {
		return err
	}
	tx, err := buildDeployTx(contractFile, account, value, nonce)
	if err != nil {
		return err
	}
	password := common.GetPassPhrase("", false)
	fmt.Println("deploy contract...")
	signedTx, err := utils.SignTxByDir(account, tx, networkId, keystoreDir, &password)
	rawTx, err := rlp.EncodeToBytes(signedTx)
	txHash, err := common.SendRawTX(rawTx, endpoint)
	if err != nil {
		return err
	}

	timer := time.NewTimer(time.Second)
	for i := 0; i < 5; i++ {
		<-timer.C
		addr, err := common.GetContractAddress(txHash, endpoint)
		if err == nil {
			fmt.Printf("deploy finished, contract address: %s\n", addr)
			return nil
		}
		timer.Reset(2 * time.Second)
	}
	return fmt.Errorf("failed to deploy contract, please retry")
}

func buildDeployTx(contractFile, account string, value *big.Int, nonce uint64) (*types.Transaction, error) {
	file, err := os.Open(contractFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open contract file, as: %e", err)
	}
	contractBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read the content of contract file, as: %e", err)
	}

	from := util.HexToAddress(account)
	tx := &types.Transaction{
		Data: types.TxData{
			From:         &from,
			Amount:       value,
			GasLimit:     math.MaxUint64 / 2,
			Price:        big.NewInt(1),
			Payload:      contractBytes,
			AccountNonce: nonce,
		},
	}
	return tx, nil
}

func accountAddr(account accounts.Account) *types.Address {
	addr := &types.Address{}
	addrBytes := account.Address.Bytes()
	copy(addr[:], addrBytes)
	return addr
}
