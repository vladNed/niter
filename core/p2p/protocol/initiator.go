package protocol

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/logging"
	"github.com/indexone/niter/core/utils"
)

var logger = logging.NewLogger(config.Config.LogLevel)

type InitiatorState struct {

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

func (i *InitiatorState) RunEventHandler() {
	for {
		select {
		case <-i.ctx.Done():
			return
		case event := <-i.eventChannel:
			i.handleSwapEvent(event)
		}
	}
}

func (i *InitiatorState) handleSwapEvent(event SEvents) {
	switch event {
	case SInit:
		err := i.handleSInit()
		if err != nil {
			logger.Error("Error handling SInit event: ", err.Error())
			// TODO: Handle error in events
		}
	case SInitDone:
		logger.Debug("InitiatorState: SInitDone event received")
	default:
		logger.Debug("InitiatorState: Unknown event")
	}
}

func (i *InitiatorState) Close() {
	i.cancel()
}

func (i *InitiatorState) Start() {
	logger.Debug("InitiatorState: Starting the state machine")
	go i.RunEventHandler()
	i.eventChannel <- SInit
}

func NewInitiatorState(offerDetails *schemas.OfferDetails, sendPeerChannel chan schemas.SwapMessage) *InitiatorState {
	ctx, cancel := context.WithCancel(context.Background())
	return &InitiatorState{
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

func (i *InitiatorState) handleSInit() error {
	logger.Debug("SInit event handling started")

	// TODO: Check the balance of the initiator is enough to lock the funds
	// TODO: Give a go to the peer

	secret, err := utils.GenerateSeed()
	if err != nil {
		logger.Error("Error generating seed")
		return err
	}
	proof := utils.Hash(secret)
	proofData, err := hex.DecodeString(proof)
	i.secret = secret
	i.secretProof = proofData

	if err != nil {
		logger.Error("Error decoding proof")
		return err
	}
	i.sendPeerChannel <- schemas.SwapMessage{
		Type:    schemas.Secret,
		Payload: proofData,
	}
	peerMessage := <-i.sendPeerChannel
	if peerMessage.Type != schemas.Secret {
		logger.Error("Invalid message type received")
		return errors.New("invalid message type received")
	}
	i.peerProof = peerMessage.Payload
	i.events = append(i.events, SInitDone)
	i.eventChannel <- SInitDone
	return nil
}
