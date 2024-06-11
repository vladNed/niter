package protocol

const (
	// The data channel official name. This is the name of the atomic swap protocol channel used for communication
	// With this, all other connections are ignored
	DATA_CHANNEL_LABEL = "atomic-swap-data-channel"

	// Peer Human Readable Part
	HRP = "nit"

	// The atomic swap protocol decimals
	EGLD_DECIMALS = 1e18
	BTC_DECIMALS  = 1e8

	// The atomic swap protocol version
	VERSION = "1"

	CLAIM_COMMITMENT_KEY  = "claim_commitment"
	REFUND_COMMITMENT_KEY = "refund_commitment"
	SWAP_STATE_KEY        = "swap"
	// TODO: Add more keys

)

type Coin string

const (
	EGLD Coin = "EGLD"
	BTC  Coin = "BTC"
)

func (c Coin) String() string {
	return string(c)
}
