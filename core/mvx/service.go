package mvx

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/utils"
)

func getCall(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %d", resp.StatusCode)
	}
	jsonData, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func GetAddressBalance(address string) (*big.Int, error) {
	resp, err := getCall(config.Config.MvxGatewatURL + "/address/" + address + "/balance")
	if err != nil {
		return nil, err
	}
	var balanceResponse BalanceResponse
	err = json.Unmarshal(resp, &balanceResponse)
	if err != nil {
		return nil, err
	}

	balance := utils.ToBigInt(balanceResponse.Data.Balance)

	return balance, nil
}

func GetTransactionResult(txHash string) (*TransactionData, error) {
	resp, err := getCall(config.Config.MvxGatewatURL + "/transaction/" + txHash + "?withResults=true")
	if err != nil {
		return nil, err
	}
	var transactionResponse TransactionResponse
	err = json.Unmarshal(resp, &transactionResponse)
	if err != nil {
		return nil, err
	}
	result := transactionResponse.Data.Transaction

	return &result, nil
}


func GetContractStorageKeys(contractAddress string) (map[string]string, error) {
	resp, err := getCall(config.Config.MvxGatewatURL + "/address/" + contractAddress + "/keys")
	if err != nil {
		return nil, err
	}

	var contractStorageKeys ContractStorageKeys
	err = json.Unmarshal(resp, &contractStorageKeys)
	if err != nil {
		return nil, err
	}

	return contractStorageKeys.Data.Pairs, nil
}
