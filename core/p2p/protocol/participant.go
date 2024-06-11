package protocol

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/mvx"
	"github.com/indexone/niter/core/utils"
)

type ParticipantState struct {
	// Managing the state of the Initiator
	ctx                context.Context
	cancel             context.CancelFunc
	eventChannel       chan SEventMessage
	swapSendChannel    chan schemas.SwapMessage
	swapReceiveChannel chan schemas.SwapMessage

	// Offer Details
	receivingAmount     string
	receivingCurrency   string
	sendingAmount       string
	sendingCurrency     string
	swapContractAddress string

	swapHeight uint64

	// State Machine
	events      []SEvents
	secret      []byte
	secretProof []byte
	peerProof   []byte
}

func NewParticipantState(
	offerDetails *schemas.OfferDetails,
	swapSendChannel chan schemas.SwapMessage,
	swapReceiveChannel chan schemas.SwapMessage,
) *ParticipantState {
	ctx, cancel := context.WithCancel(context.Background())
	return &ParticipantState{
		ctx:                ctx,
		cancel:             cancel,
		eventChannel:       make(chan SEventMessage),
		swapSendChannel:    swapSendChannel,
		swapReceiveChannel: swapReceiveChannel,
		receivingAmount:    offerDetails.ReceivingAmount,
		receivingCurrency:  offerDetails.ReceivingCurrency,
		sendingAmount:      offerDetails.SendingAmount,
		sendingCurrency:    offerDetails.SendingCurrency,
		swapHeight:         0,
		events:             []SEvents{},
	}
}

func (s *ParticipantState) RunEventHandler() {
	for {
		select {
		case <-s.ctx.Done():
			logger.Error("ParticipantState: Context cancelled")
			return
		case eventMessage := <-s.eventChannel:
			go s.handleSwapEvent(eventMessage.Event, eventMessage.Data)
		}
	}
}

func (s *ParticipantState) handleSwapEvent(event SEvents, data string) {
	switch event {
	case SInit:
		s.events = append(s.events, event)
		err := s.handleSInit()
		if err != nil {
			s.eventChannel <- SEventMessage{Event: SFailed, Data: ""}
		}
	case SInitDone:
		s.events = append(s.events, event)
		err := s.handleSInitDone()
		if err != nil {
			s.eventChannel <- SEventMessage{Event: SFailed, Data: ""}
		}
	case SLockedEGLD:
		s.events = append(s.events, event)
		err := s.handleSLockedEGLD()
		if err != nil {
			s.eventChannel <- SEventMessage{Event: SFailed, Data: ""}
		}
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
	s.eventChannel <- SEventMessage{Event: SInit, Data: ""}
}

func (s *ParticipantState) GetEvents() []SEvents {
	return s.events
}

func (s *ParticipantState) GetTransactionDetails(requestType TransactionRequestType) (map[string]interface{}, error) {
	return nil, nil
}

func (s *ParticipantState) handleSInit() error {
	logger.Debug("[P] SInit event handling started")

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
	s.swapSendChannel <- schemas.SwapMessage{
		Type:    schemas.Secret,
		Payload: s.secretProof,
	}
	peerMessage := <-s.swapReceiveChannel
	if peerMessage.Type != schemas.Secret {
		logger.Error("Expected Secret message from peer")
		return errors.New("invalid secret proof received")
	}
	s.peerProof = peerMessage.Payload
	s.events = append(s.events, SInit)
	s.eventChannel <- SEventMessage{Event: SInitDone, Data: ""}

	return nil
}

func (s *ParticipantState) checkSwapContractRequirements() error {
	storageKeys, err := mvx.GetContractStorageKeys(s.swapContractAddress)
	if err != nil {
		return err
	}

	// Check swap state
	swapKey := hex.EncodeToString([]byte(SWAP_STATE_KEY))
	if _, ok := storageKeys[swapKey]; ok {
		return errors.New("swap state not new")
	}

	// Check refund commitment
	refundKey := hex.EncodeToString([]byte(REFUND_COMMITMENT_KEY))
	refundStorage, ok := storageKeys[refundKey]
	if !ok {
		return errors.New("refund commitment not found")
	}
	strRefundProof, _ := hex.DecodeString(refundStorage)
	strPeerProof := hex.EncodeToString(s.peerProof)
	if string(strRefundProof) != strPeerProof {
		return errors.New("refund commitment mismatch")
	}

	// Check claim commitment
	claimKey := hex.EncodeToString([]byte(CLAIM_COMMITMENT_KEY))
	claimStorage, ok := storageKeys[claimKey]
	if !ok {
		return errors.New("claim commitment not found")
	}
	strClaimProof, _ := hex.DecodeString(claimStorage)
	strSecretProof := hex.EncodeToString(s.secretProof)
	if string(strClaimProof) != strSecretProof {
		return errors.New("claim commitment mismatch")
	}

	return nil
}

func (s *ParticipantState) handleSInitDone() error {
	logger.Debug("[P] SInitDone event handling started")
	peerMessage := <-s.swapReceiveChannel
	if peerMessage.Type != schemas.ContractCreated {
		logger.Error(" [P] Expected ContractCreated message from peer")
		return errors.New("invalid contract created message received")
	}

	logger.Debug(string(peerMessage.Payload))

	var eventDataJson SLockedEGLDData
	err := json.Unmarshal(peerMessage.Payload, &eventDataJson)
	if err != nil {
		logger.Error("Error unmarshalling SLockedEGLDData")
		return err
	}

	transaction, err := mvx.GetTransactionResult(eventDataJson.TxHash)
	if err != nil {
		logger.Error("Error getting transaction result")
		return err
	}

	result := transaction.SmartContractResults[0]
	deployedContractAddress, err := mvx.ParseDeployResult(&result)
	if err != nil {
		logger.Error("Error parsing deploy result")
		return err
	}

	s.swapContractAddress = deployedContractAddress
	if err := s.checkSwapContractRequirements(); err != nil {
		logger.Error("Checking smart contract failed:" + err.Error())
		return err
	}

	s.eventChannel <- SEventMessage{Event: SLockedEGLD, Data: ""}
	return nil
}

func (s *ParticipantState) handleSLockedEGLD() error {
	logger.Debug("[P] SLockedEGLD event handling started")
	return nil
}
