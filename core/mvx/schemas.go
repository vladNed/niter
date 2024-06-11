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

type SmartContractResult struct {
	Hash string `json:"hash"`
	Data string `json:"data"`
}

type TransactionData struct {
	Type                 string                `json:"type"`
	Hash                 string                `json:"hash"`
	Nonce                int                   `json:"nonce"`
	Round                int                   `json:"round"`
	Epoch                int                   `json:"epoch"`
	Value                string                `json:"value"`
	Function             string                `json:"function"`
	SmartContractResults []SmartContractResult `json:"smartContractResults"`
}

type TransactionResponse struct {
	Data struct {
		Transaction TransactionData `json:"transaction"`
	} `json:"data"`
	Code  string `json:"code"`
	Error string `json:"error"`
}

type ContractStorageKeys struct {
	Data struct {
		Pairs map[string]string `json:"pairs"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
