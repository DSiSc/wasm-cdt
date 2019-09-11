package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// GetAccountNonde get account nonce
func GetAccountNonde(account, endpoint string) (uint64, error) {
	reqTemplateStr := `{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x%s","%s"],"id":1}`
	resp, err := doPost(fmt.Sprintf(reqTemplateStr, account, "pending"), endpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to get account nonce, as: %v", err)
	}
	nonceStr := new(string)
	json.Unmarshal(resp.Result, nonceStr)
	nonce := uint64(hexstr2dec(*nonceStr))
	return nonce, nil
}

// SendRawTX send raw transaction to endpoint.
func SendRawTX(rawTx []byte, endpoint string) (string, error) {
	reqData := fmt.Sprintf(
		`{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0x%x"],"id":1}`, rawTx)
	resp, err := doPost(reqData, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to send raw transaction, as: %v", err)
	}
	txHash := new(string)
	json.Unmarshal(resp.Result, txHash)
	return fmt.Sprintf("%s", *txHash), nil
}

// GetContractAddress get contract address by tx hash.
func GetContractAddress(txHash, endpoint string) (string, error) {
	reqData := fmt.Sprintf(
		`{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["%s"],"id":1}`, txHash)
	resp, err := doPost(reqData, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction receipt, as: %v", err)
	}
	results := make(map[string]string)
	json.Unmarshal(resp.Result, &results)
	if val, ok := results["contractAddress"]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("failed to get transaction's receipt")
	}
}

// doPost is a tool function used to talk to justitia API.
func doPost(reqData string, endpoint string) (*RPCResponse, error) {
	request, err := http.NewRequest("POST", endpoint, strings.NewReader(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request as: %v", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request to endpoint as: %v", err)
	}
	defer response.Body.Close()
	blob, _ := ioutil.ReadAll(response.Body)
	recv := new(RPCResponse)
	json.Unmarshal(blob, recv)
	if recv.Error != nil {
		return nil, fmt.Errorf("%v", recv.Error)
	}
	return recv, nil
}

// hexstr2dec converts hexadecimal string to decimal integer(int64).
func hexstr2dec(hex string) int64 {
	var str string
	if hex[0:2] == "0x" {
		str = hex[2:]
	} else {
		str = hex
	}
	dec, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		panic(err)
	}
	return dec
}
