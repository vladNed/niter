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

func GetAddressBalance(address string) (*big.Int, error) {
	resp, err := http.Get(config.Config.MvxGatewatURL + "/address/" + address + "/balance")
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
	var balanceResponse BalanceResponse
	err = json.Unmarshal(jsonData, &balanceResponse)
	if err != nil {
		return nil, err
	}

	balance := utils.ToBigInt(balanceResponse.Data.Balance)

	return balance, nil
}
