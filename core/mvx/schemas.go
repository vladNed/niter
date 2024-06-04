package mvx

type BalanceResponse struct {
	Data struct {
		Balance    string `json:"balance"`
		BlockNonce struct {
			Nonce    int    `json:"nonce"`
			Hash     string `json:"hash"`
			RootHash string `json:"rootHash"`
		} `json:"blockInfo"`
	} `json:"data"`
	Code  string `json:"code"`
	Error string `json:"error"`
}
