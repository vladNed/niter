package protocol

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/logging"
	"github.com/indexone/niter/core/mvx"
	"github.com/indexone/niter/core/utils"
)

var logger = logging.NewLogger(config.Config.LogLevel)

type InitiatorState struct {

	// Managing the state of the Initiator
	ctx             context.Context
	cancel          context.CancelFunc
	eventChannel    chan SEvents
	sendPeerChannel chan schemas.SwapMessage

	mvxAddress      string

	// Offer Details
	receivingAmount   string
	receivingCurrency string
	sendingAmount     string
	sendingCurrency   string
	isCreator         bool

	swapHeight uint64

	// State Machine
	events      []SEvents
	secret      []byte
	secretProof []byte
	peerProof   []byte
}

func NewInitiatorState(
	offerDetails *schemas.OfferDetails,
	sendPeerChannel chan schemas.SwapMessage,
	mvxAddress string,
	isCreator bool,
) *InitiatorState {
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
		mvxAddress:        mvxAddress,
	}
}

func (i *InitiatorState) RunEventHandler() {
	for {
		select {
		case <-i.ctx.Done():
			logger.Debug("InitiatorState: Context cancelled")
			return
		case event := <-i.eventChannel:
			go i.handleSwapEvent(event)
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
		err := i.handleSInitDone()
		if err != nil {
			logger.Error("Error handling SInitDone event: ", err.Error())
			// TODO: Handle error in events
		}
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

func (i *InitiatorState) GetEvents() []SEvents {
	return i.events
}

func (i *InitiatorState) GetTransactionDetails(requestType TransactionRequestType) (map[string]interface{}, error) {
	switch requestType {
	case CreateSwap:
		var amount string
		if i.sendingCurrency == EGLD.String() {
			amount = i.sendingAmount
		} else {
			amount = i.receivingAmount
		}
		data := map[string]interface{}{
			"claimProof": hex.EncodeToString(i.peerProof),
			"refundProof": hex.EncodeToString(i.secretProof),
			"amount": amount,
		}
		return data, nil
	default:
		return nil, errors.New("invalid request type")
	}
}

func (i *InitiatorState) checkEnoughBalance() (error) {
	mvxBalance, err := mvx.GetAddressBalance(i.mvxAddress)
	if err != nil {
		logger.Error("Error getting balance: ", err.Error())
		return err
	}
	if i.isCreator && i.sendingCurrency == EGLD.String() {
		sendingAmount := utils.ToBigInt(i.sendingAmount)
		if mvxBalance.Cmp(sendingAmount) == -1 {
			logger.Error("Insufficient balance")
			return errors.New("insufficient balance")
		}
		return nil
	}

	if !i.isCreator && i.receivingCurrency == EGLD.String() {
		receivingAmount := utils.ToBigInt(i.receivingAmount)
		if mvxBalance.Cmp(receivingAmount) == -1 {
			logger.Error("Insufficient balance")
			return errors.New("insufficient balance")
		}
		return nil
	}

	return nil
}

func (i *InitiatorState) handleSInit() error {
	logger.Debug("SInit event handling started")
	err := i.checkEnoughBalance()
	if err != nil {
		logger.Error("Error checking balance")
		return err
	}
	secret, err := utils.GenerateSeed()
	if err != nil {
		logger.Error("Error generating seed")
		return err
	}
	proof := utils.Hash(secret)
	proofData, err := hex.DecodeString(proof)
	if err != nil {
		logger.Error("Error decoding proof")
		return err
	}
	i.secret = secret
	i.secretProof = proofData
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
	i.events = append(i.events, SInit)
	i.eventChannel <- SInitDone
	return nil
}

func (i *InitiatorState) handleSInitDone() error {
	logger.Debug("SInitDone event handling started")

	return nil
}
