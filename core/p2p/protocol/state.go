package protocol

// The Peer state machine enumeration
type PeerState int

const (

	// PeerIdle is the initial state of the peer
	PeerIdle PeerState = iota

	// PeerInitiator is the state where the peer has initiated a connection
	// created a data channel and waits for a responder to join
	PeerInitiator

	// PeerResponder is the state where the peer has received an offer
	// and is ready to create an answer
	PeerResponder

	// PeerNegotiating is the state where the peer has create a new offer
	// and is waiting for the other peers to respond
	PeerNegotiating

	// PeerConnected is the state where the peer has successfully connected
	PeerConnected

	// PeerCommunicating is the state where the peer is able to communicate
	PeerAuthenticating

	// Communication is now possible
	PeerCommunicating
)

func (p PeerState) String() string {
	switch p {
	case PeerIdle:
		return "PeerIdle"
	case PeerInitiator:
		return "PeerInitiator"
	case PeerResponder:
		return "PeerResponder"
	case PeerNegotiating:
		return "PeerNegotiating"
	case PeerConnected:
		return "PeerConnected"
	case PeerCommunicating:
		return "PeerCommunicating"
	case PeerAuthenticating:
		return "PeerAuthenticating"
	default:
		return "Unknown"
	}
}

// PeerEvents is the enumeration of the events that can be triggered on the peer
type PeerEvents int

const (
	UnknownEvent PeerEvents = iota
	ResponderICECandidate
	InitiatorICECandidate
)




