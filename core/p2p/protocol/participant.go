package protocol

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/utils"
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

func (s *ParticipantState) RunEventHandler() {
	for {
		select {
		case <-s.ctx.Done():
			logger.Error("ParticipantState: Context cancelled")
			return
		case event := <-s.eventChannel:
			go s.handleSwapEvent(event)
		}
	}
}

func (s *ParticipantState) handleSwapEvent(event SEvents) {
	switch event {
	case SInit:
		err := s.handleSInit()
		if err != nil {
			logger.Error("Error handling SInit event: ", err.Error())
			// TODO: Handle error in events
		}
	case SInitDone:
		logger.Debug("InitiatorState: SInitDone")
	default:
		logger.Debug("InitiatorState: Unknown event")
	}
}

func (s *ParticipantState) Close() {
	s.cancel()
}

func (s *ParticipantState) Start() {
	logger.Debug("ParticipantState: Starting the state machine")
	go s.RunEventHandler()
	s.eventChannel <- SInit
}

func (s *ParticipantState) GetEvents() []SEvents {
	return s.events
}

func (s *ParticipantState) handleSInit() error {
	logger.Debug("SInit event handling started")

	// TODO: Check the balance of the participant to lock the funds

	secret, err := utils.GenerateSeed()
	if err != nil {
		logger.Error("Error generating secret: ", err.Error())
		return err
	}
	proof := utils.Hash(secret)
	proofData, err := hex.DecodeString(proof)
	if err != nil {
		logger.Error("Error decoding proof")
		return err
	}
	s.secret = secret
	s.secretProof = proofData
	s.sendPeerChannel <- schemas.SwapMessage{
		Type:    schemas.Secret,
		Payload: s.secretProof,
	}
	peerMessage := <-s.sendPeerChannel
	if peerMessage.Type != schemas.Secret {
		logger.Error("Expected Secret message from peer")
		return errors.New("invalid secret proof received")
	}
	s.peerProof = peerMessage.Payload
	s.events = append(s.events, SInit)
	s.eventChannel <- SInitDone

	return nil
}