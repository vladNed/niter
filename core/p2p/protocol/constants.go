package protocol


// The data channel official name. This is the name of the atomic swap protocol channel used for communication
// With this, all other connections are ignored
const DATA_CHANNEL_LABEL = "atomic-swap-data-channel"

// Peer Human Readable Part
const HRP = "nit"

type Coin string
const (
	EGLD Coin = "EGLD"
	BTC Coin = "BTC"
)

func (c Coin) String() string {
	return string(c)
}

const EGLD_DECIMALS = 1e18
const BTC_DECIMALS = 1e8
