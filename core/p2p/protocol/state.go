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

// The events that can be triggered during a swap
type SEvents int

const (
	SInit SEvents = iota
	SInitDone
	SLockedEGLD
	SLockedBTC
	SRefund
	SClaimed
	SOk
	SFailed
	Unknown
)

func (se *SEvents) String() string {
	switch *se {
	case SInit:
		return "SInit"
	case SInitDone:
		return "SInitDone"
	case SLockedEGLD:
		return "SLockedEGLD"
	case SLockedBTC:
		return "SLockedBTC"
	case SRefund:
		return "SRefund"
	case SClaimed:
		return "SClaimed"
	case SOk:
		return "SOk"
	case SFailed:
		return "SFailed"
	default:
		return "Unknown"
	}
}

func SEventsFromString(s string) SEvents {
	switch s {
	case "SInit":
		return SInit
	case "SInitDone":
		return SInitDone
	case "SLockedEGLD":
		return SLockedEGLD
	case "SLockedBTC":
		return SLockedBTC
	case "SRefund":
		return SRefund
	case "SClaimed":
		return SClaimed
	case "SOk":
		return SOk
	case "SFailed":
		return SFailed
	default:
		return Unknown
	}
}

// The message that is passed to thee state machine from the event emitter
type SEventMessage struct {
	Event SEvents
	Data  string
}

type SLockedEGLDData struct {
	TxHash string `json:"hash"`
}

type TransactionRequestType string

const (
	CreateSwap   TransactionRequestType = "CreateSwap"
	SetReadySwap TransactionRequestType = "SetReadySwap"
	ClaimSwap    TransactionRequestType = "ClaimSwap"
	RefundSwap   TransactionRequestType = "RefundSwap"
)

func (t TransactionRequestType) String() string {
	return string(t)
}

var transactionRequestTypeMap = map[string]TransactionRequestType{
	"CreateSwap":   CreateSwap,
	"SetReadySwap": SetReadySwap,
	"ClaimSwap":    ClaimSwap,
	"RefundSwap":   RefundSwap,
}

func TransactionRequestTypeFromString(s string) TransactionRequestType {
	return transactionRequestTypeMap[s]
}

type SwapState interface {
	Start()
	Close()
	RunEventHandler()
	handleSwapEvent(event SEvents, data string)
	GetEvents() []SEvents
	GetTransactionDetails(requestType TransactionRequestType) (map[string]interface{}, error)
}
