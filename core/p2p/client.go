package p2p

import (
	"context"
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/pion/webrtc/v4"

	"github.com/indexone/niter/core/bitcoin"
	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/crypto"
	"github.com/indexone/niter/core/discovery"
	msgSchemas "github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/logging"
	"github.com/indexone/niter/core/mvx"
	"github.com/indexone/niter/core/p2p/protocol"
	"github.com/indexone/niter/core/utils"
)

var logger = logging.NewLogger(config.Config.LogLevel)

type RemotePeerInfo struct {
	Id string
}

type Peer struct {

	// WebRTC Peer Data
	LocalConnection   *webrtc.PeerConnection
	DataChannel       *webrtc.DataChannel
	State             protocol.PeerState
	KeyPair           *crypto.NetworkKey
	RemotePeer        *RemotePeerInfo
	p2pEventsChannel  chan protocol.PeerEvents
	swapEventsChannel chan protocol.SEventMessage
	msgChannel        chan msgSchemas.Message

	// Swap Data
	ActiveOfferId string
	swapChannel   chan msgSchemas.SwapMessage
	SwapState     protocol.SwapState

	// Context
	ctx    context.Context
	cancel context.CancelFunc

	// Wallets
	btcWallet *bitcoin.Wallet
	mvxWallet *mvx.Wallet
}

func NewPeer(
	p2pEventsChannel chan protocol.PeerEvents,
	swapEventsChannel chan protocol.SEventMessage,
	msgChannel chan msgSchemas.Message,
	btcWallet *bitcoin.Wallet,
	mvxWallet *mvx.Wallet,
) (*Peer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	keyPair, err := crypto.GenerateKey()
	if err != nil {
		logger.Warn("Error generating key pair: ", err.Error())
		cancel()
		return nil, err
	}
	peer := &Peer{
		LocalConnection:  nil,
		DataChannel:      nil,
		State:            protocol.PeerIdle,
		KeyPair:          keyPair,
		RemotePeer:       nil,
		p2pEventsChannel: p2pEventsChannel,
		swapEventsChannel: swapEventsChannel,
		ActiveOfferId:    "",
		msgChannel:       msgChannel,
		SwapState:        nil,
		ctx:              ctx,
		cancel:           cancel,
		swapChannel:      make(chan msgSchemas.SwapMessage),
		btcWallet:        btcWallet,
		mvxWallet:        mvxWallet,
	}
	logger.Debug("Peer initialized")
	go peer.MessageHandler()
	return peer, nil
}

func (p *Peer) Id() string {
	commitment := p.KeyPair.Commitment()
	encodedCommitment := base32.StdEncoding.EncodeToString([]byte(commitment))[:20]
	return fmt.Sprintf("%s1%s", protocol.HRP, encodedCommitment)
}

// StartInitiator starts the peer as an initiator node
// This is mainly used to start the peer that handles the data channel creation
// and the offer creation.
func (p *Peer) StartInitiator() error {
	webRtcConfig := config.GetICEConfiguration()
	peerConn, err := webrtc.NewPeerConnection(webRtcConfig)
	if err != nil {
		logger.Warn("Error creating peer connection: ", err.Error())
		return err
	}
	p.LocalConnection = peerConn
	p.State = protocol.PeerInitiator
	p.setupConnectionCallbacks()
	p.LocalDataChannelHandlers()
	return nil
}

// StartResponder starts the peer as a responder node
// This is mainly used to start the peer that handles the remote received offer
// and the data channel.
func (p *Peer) StartResponder() error {
	webRtcConfig := config.GetICEConfiguration()
	peerConn, err := webrtc.NewPeerConnection(webRtcConfig)
	if err != nil {
		logger.Warn("Error creating peer connection: ", err.Error())
		return err
	}
	p.LocalConnection = peerConn
	p.State = protocol.PeerResponder
	p.setupConnectionCallbacks()
	p.RemoteDataChannelHandlers()
	return nil
}

// ResetPeer resets the peer to its initial state
func (p *Peer) ResetPeer() {
	p.LocalConnection = nil
	p.DataChannel = nil
	p.State = protocol.PeerIdle
	p.RemotePeer = nil
	p.cancel()
	p.ctx, p.cancel = context.WithCancel(context.Background())
	go p.MessageHandler()
}

func (p *Peer) setupConnectionCallbacks() {
	p.LocalConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		if p.State > protocol.PeerIdle {
			return
		}
		switch connectionState {
		case webrtc.PeerConnectionStateConnected:
			logger.Debug("Peer connection state is connected")
			p.State = protocol.PeerConnected
		case webrtc.PeerConnectionStateDisconnected:
			logger.Debug("Peer connection state is disconnected")
			go p.ResetPeer()
		case webrtc.PeerConnectionStateFailed:
			logger.Debug("Peer connection state has failed")
			go p.ResetPeer()
		case webrtc.PeerConnectionStateNew:
			logger.Debug("Peer connection state is new")
			p.State = protocol.PeerNegotiating
		}
	})
	p.LocalConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			return
		}

		switch p.State {
		case protocol.PeerInitiator:
			logger.Debug("Gathering complete")
			p.p2pEventsChannel <- protocol.InitiatorICECandidate
		case protocol.PeerResponder:
			logger.Debug("Gathering complete")
			p.p2pEventsChannel <- protocol.ResponderICECandidate
		default:
			logger.Warn("Unknown peer state 1")
		}
	})
}

func (p *Peer) CreateOffer() (*webrtc.SessionDescription, error) {
	offer, err := p.LocalConnection.CreateOffer(&webrtc.OfferOptions{ICERestart: true})
	if err != nil {
		logger.Warn("Error creating offer: ", err.Error())
		return nil, err
	}
	err = p.LocalConnection.SetLocalDescription(offer)
	if err != nil {
		logger.Warn("Error setting local description: ", err.Error())
		return nil, err
	}
	return &offer, nil
}

func (p *Peer) CreateAnswer() (*webrtc.SessionDescription, error) {
	answer, err := p.LocalConnection.CreateAnswer(nil)
	if err != nil {
		logger.Warn("Error creating answer: ", err.Error())
		return nil, err
	}
	err = p.LocalConnection.SetLocalDescription(answer)
	if err != nil {
		logger.Warn("Error setting local description: ", err.Error())
		return nil, err
	}
	return &answer, nil
}

func (p *Peer) SetOffer(encodedSDP string) error {
	offer, err := utils.DecodeSDP(encodedSDP)
	if err != nil {
		logger.Warn("Error decoding SDP: ", err.Error())
		return err
	}
	err = p.LocalConnection.SetRemoteDescription(*offer)
	if err != nil {
		logger.Warn("Error setting remote description: ", err.Error())
		return err
	}

	return nil
}

func (p *Peer) SendData(data []byte) error {
	err := p.DataChannel.Send(data)
	if err != nil {
		logger.Warn("Error sending data: ", err.Error())
		return err
	}

	return nil
}

func (p *Peer) peerAuthentication(msgData webrtc.DataChannelMessage) error {
	peerId := string(msgData.Data)
	if peerId[:3] != protocol.HRP {
		logger.Debug("Invalid peer hrp")
		return errors.New("invalid peer id")
	}
	if len(peerId[4:]) != 20 {
		logger.Debug("Invalid peer data")
		return errors.New("invalid peer id")
	}
	p.RemotePeer = &RemotePeerInfo{Id: peerId}
	offer, ok := discovery.Cache.GetOffer(p.ActiveOfferId)
	if !ok {
		logger.Debug("Offer not found")
		return errors.New("offer not found")
	}
	go p.swapMessageHandler()
	if offer.OfferDetails.SwapCreator == p.Id() {
		if offer.OfferDetails.SendingCurrency == protocol.EGLD.String() {
			p.SwapState = protocol.NewInitiatorState(
				p.ctx,
				&offer.OfferDetails,
				p.swapChannel,
				p.swapEventsChannel,
				p.mvxWallet.Address,
				true,
			)
		} else {
			p.SwapState = protocol.NewParticipantState(&offer.OfferDetails, p.swapChannel)
		}
	} else {
		if offer.OfferDetails.ReceivingCurrency == protocol.EGLD.String() {
			p.SwapState = protocol.NewInitiatorState(
				p.ctx,
				&offer.OfferDetails,
				p.swapChannel,
				p.swapEventsChannel,
				p.mvxWallet.Address,
				false,
			)
		} else {
			p.SwapState = protocol.NewParticipantState(&offer.OfferDetails, p.swapChannel)
		}
	}
	p.SwapState.Start()
	p.State = protocol.PeerCommunicating
	return nil
}

func (p *Peer) LocalDataChannelHandlers() {
	dataChannel, err := p.LocalConnection.CreateDataChannel(protocol.DATA_CHANNEL_LABEL, nil)
	if err != nil {
		logger.Warn("Error creating data channel: ", err.Error())
		return
	}

	dataChannel.OnOpen(func() {
		p.State = protocol.PeerAuthenticating
		dataChannel.Send([]byte(p.Id()))
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		switch p.State {
		case protocol.PeerAuthenticating:
			err := p.peerAuthentication(msg)
			if err != nil {
				logger.Warn("Error authenticating peer: ", err.Error())
				go p.ResetPeer()
				return
			}
		case protocol.PeerCommunicating:
			msgData, err := msgSchemas.DeserializeSwapMessage(msg.Data)
			if err != nil {
				logger.Warn("Error deserializing swap message: ", err.Error())
				go p.ResetPeer()
				return
			}
			p.swapChannel <- *msgData
		default:
			logger.Warn("Unknown peer state 2")
		}
	})

	dataChannel.OnClose(func() {
		logger.Debug("[RECEIVING] Data channel closed")
		p.DataChannel = nil
	})

	p.DataChannel = dataChannel
}

func (p *Peer) RemoteDataChannelHandlers() {
	p.LocalConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		logger.Debug("[RECEIVING] Data channel received: ", d.Label())
		p.DataChannel = d

		d.OnOpen(func() {
			p.State = protocol.PeerAuthenticating
			d.Send([]byte(p.Id()))
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			switch p.State {
			case protocol.PeerAuthenticating:
				err := p.peerAuthentication(msg)
				if err != nil {
					logger.Warn("Error authenticating peer: ", err.Error())
					go p.ResetPeer()
					return
				}
			case protocol.PeerCommunicating:
				msgData, err := msgSchemas.DeserializeSwapMessage(msg.Data)
				if err != nil {
					logger.Warn("Error deserializing swap message: ", err.Error())
					go p.ResetPeer()
					return
				}
				p.swapChannel <- *msgData
			default:
				logger.Warn("Unknown peer state 3")
			}
		})

		d.OnClose(func() {
			logger.Debug("[RECEIVING] Data channel closed")
			p.DataChannel = nil
		})
	})
}

// The message handler listens for messages from the websocket connection
// and processes them accordingly.
func (p *Peer) MessageHandler() {
	for {
		select {
		case <-p.ctx.Done():
			logger.Debug("Disconnecting peer message handler closed")
			return
		case signallingMessage := <-p.msgChannel:
			switch msgType := signallingMessage.(type) {
			case *msgSchemas.AnswerMessage:
				p.SetOffer(msgType.AnswerSDP)
			default:
				logger.Warn("Unknown message type")
			}
		}
	}
}

// The swap message handler listens for swap message from the swap state
// manager and sends them to the remote peer.
func (p *Peer) swapMessageHandler() {
	for {
		select {
		case <-p.ctx.Done():
			logger.Debug("Swap message handler closed")
			return
		case swapCtxMessage := <-p.swapChannel:
			data := swapCtxMessage.Serialize()
			err := p.SendData(data)
			if err != nil {
				logger.Warn("Error sending swap message: ", err.Error())
				continue
			}
		}
	}
}
