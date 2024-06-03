package transactions

type TxType int
const (
	Multiversx TxType = iota
	Bitcoin
)

type Tx interface {
	Broadcast() error
	Serialize() map[string]interface{}
}

// Multiversx transaction mainly are for interacting with smart contracts
type MvxTx struct {
	FuncName string   `json:"funcName"`
	Value    string   `json:"value"`
	Args     []string `json:"args"`
}

func (tx *MvxTx) Broadcast() error {
	return nil
}

// Bitcoin transaction mainly are for locking and claiming transactions
type BtcTx struct{
	Raw string `json:"raw"`
}

func (tx *BtcTx) Broadcast() error {
	return nil
}
