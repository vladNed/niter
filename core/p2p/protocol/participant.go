package protocol

import (
	"context"

	"github.com/indexone/niter/core/discovery/schemas"
)

type ParticipantState struct {
	// Managing the state of the Initiator
	ctx             context.Context
	cancel          context.CancelFunc
	eventChannel    chan SEvents
	sendPeerChannel chan schemas.SwapMessage

	// Offer Details
	receivingAmount   string
	receivingCurrency string
	sendingAmount     string
	sendingCurrency   string

	swapHeight uint64

	// State Machine
	events      []SEvents
	secret      []byte
	secretProof []byte
	peerProof   []byte
}

func (i *ParticipantState) RunEventHandler() {
	for {
		select {
		case <-i.ctx.Done():
			return
		case event := <-i.eventChannel:
			i.handleSwapEvent(event)
		}
	}
}

func (i *ParticipantState) handleSwapEvent(event SEvents) {
	switch event {
	case SInit:
		logger.Debug("InitiatorState: SInit")
	case SLockedEGLD:
		logger.Debug("InitiatorState: SLockedEGLD")
	case SLockedBTC:
		logger.Debug("InitiatorState: SLockedBTC")
	default:
		logger.Debug("InitiatorState: Unknown event")
	}
}

func (i *ParticipantState) Close() {
	i.cancel()
}

func (i *ParticipantState) Start() {
	logger.Debug("InitiatorState: Starting the state machine")
	go i.RunEventHandler()
	i.eventChannel <- SInit
}

func NewParticipantState(offerDetails *schemas.OfferDetails, sendPeerChannel chan schemas.SwapMessage) *ParticipantState {
	ctx, cancel := context.WithCancel(context.Background())
	return &ParticipantState{
		ctx:               ctx,
		cancel:            cancel,
		eventChannel:      make(chan SEvents),
		sendPeerChannel:   sendPeerChannel,
		receivingAmount:   offerDetails.ReceivingAmount,
		receivingCurrency: offerDetails.ReceivingCurrency,
		sendingAmount:     offerDetails.SendingAmount,
		sendingCurrency:   offerDetails.SendingCurrency,
		swapHeight:        0,
		events:            []SEvents{},
	}
}
