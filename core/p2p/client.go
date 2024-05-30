package p2p

import (
	"encoding/base32"
	"fmt"

	"github.com/pion/webrtc/v4"

	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/crypto"
	"github.com/indexone/niter/core/logging"
	"github.com/indexone/niter/core/p2p/protocol"
	"github.com/indexone/niter/core/utils"
)

var logger = logging.NewLogger(config.Config.LogLevel)

type PeerData struct {
	Side      string `json:"side"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

type RemotePeerInfo struct {
	Id string
}

type Peer struct {
	LocalConnection *webrtc.PeerConnection
	DataChannel     *webrtc.DataChannel
	State           protocol.PeerState
	ExchangeData    []PeerData
	KeyPair         *crypto.NetworkKey
	RemotePeer      *RemotePeerInfo
	eventsChannel   chan protocol.PeerEvents
}

func NewPeer(eventChannel chan protocol.PeerEvents) (*Peer, error) {
	keyPair, err := crypto.GenerateKey()
	if err != nil {
		logger.Warn("Error generating key pair: ", err.Error())
		return nil, err
	}
	peer := &Peer{
		LocalConnection: nil,
		DataChannel:     nil,
		State:           protocol.PeerIdle,
		KeyPair:         keyPair,
		RemotePeer:      nil,
		eventsChannel:   eventChannel,
	}
	logger.Debug("Peer initialized")
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

func (p *Peer) ResetPeer() {
	p.LocalConnection = nil
	p.DataChannel = nil
	p.State = protocol.PeerIdle
	p.ExchangeData = nil
	p.RemotePeer = nil
}

func (p *Peer) setupConnectionCallbacks() {
	p.LocalConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
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
		if candidate == nil {
			switch p.State {
			case protocol.PeerInitiator:
				logger.Debug("Gathering complete")
				p.eventsChannel <- protocol.InitiatorICECandidate
			case protocol.PeerResponder:
				logger.Debug("Gathering complete")
				p.eventsChannel <- protocol.ResponderICECandidate
			default:
				logger.Warn("Unknown peer state")
			}
			return
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

	dataStruct := PeerData{
		Side:      "local",
		Data:      string(data),
		Timestamp: utils.GetTimestamp(),
	}
	p.ExchangeData = append(p.ExchangeData, dataStruct)

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
			p.RemotePeer = &RemotePeerInfo{
				Id: string(msg.Data),
			}
			p.State = protocol.PeerCommunicating
		case protocol.PeerCommunicating:
			p.ExchangeData = append(p.ExchangeData, PeerData{
				Side:      "remote",
				Data:      string(msg.Data),
				Timestamp: utils.GetTimestamp(),
			})
		default:
			logger.Warn("Unknown peer state")
		}
	})

	dataChannel.OnClose(func() {
		logger.Debug("[RECEIVING] Data channel closed")
		p.DataChannel = nil
		p.ExchangeData = nil
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
				p.RemotePeer = &RemotePeerInfo{
					Id: string(msg.Data),
				}
				p.State = protocol.PeerCommunicating
			case protocol.PeerCommunicating:
				p.ExchangeData = append(p.ExchangeData, PeerData{
					Side:      "remote",
					Data:      string(msg.Data),
					Timestamp: utils.GetTimestamp(),
				})
			default:
				logger.Warn("Unknown peer state")
			}
		})

		d.OnClose(func() {
			logger.Debug("[RECEIVING] Data channel closed")
			p.DataChannel = nil
			p.ExchangeData = nil
		})
	})
}
