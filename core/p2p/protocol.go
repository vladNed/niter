package p2p

// The data channel official name. This is the name of the atomic swap protocol channel used for communication
// With this, all other connections are ignored
const DATA_CHANNEL_LABEL = "atomic-swap-data-channel"

// The Peer state machine enumeration
type PeerState int

const (

	// PeerIdle is the initial state of the peer
	PeerIdle PeerState = iota

	// PeerNegotiating is the state where the peer has create a new offer
	// and is waiting for the other peers to respond
	PeerNegotiating
)
