package p2p

import (
	"os"

	"github.com/pion/webrtc/v4"

	"github.com/indexone/niter/core/config"
	"github.com/indexone/niter/core/logging"
)

var logger = logging.NewLogger(config.Config.LogLevel)

type Peer struct {
	LocalConnection *webrtc.PeerConnection
	DataChannel     *webrtc.DataChannel
	State           PeerState
}

func NewPeer() (*Peer, error) {
	webRtcConfig := config.GetICEConfiguration()
	peerConn, err := webrtc.NewPeerConnection(webRtcConfig)
	if err != nil {
		logger.Warn("Error creating peer connection: ", err.Error())
		return nil, err
	}

	dataChannel, err := peerConn.CreateDataChannel(DATA_CHANNEL_LABEL, nil)
	if err != nil {
		logger.Warn("Error creating data channel: ", err.Error())
		return nil, err
	}

	peer := &Peer{
		LocalConnection: peerConn,
		DataChannel:     dataChannel,
		State:           PeerIdle,
	}

	peer.setupConnectionCallbacks()
	peer.setupDataChannelProtocol()

	logger.Debug("Peer initialized")

	return peer, nil
}

func (p *Peer) setupConnectionCallbacks() {
	p.LocalConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		switch connectionState {
		case webrtc.PeerConnectionStateConnected:
			logger.Debug("Peer connection state is connected")
		case webrtc.PeerConnectionStateDisconnected:
			logger.Debug("Peer connection state is disconnected")
		case webrtc.PeerConnectionStateFailed:
			logger.Debug("Peer connection state has failed")
		case webrtc.PeerConnectionStateClosed:
			logger.Debug("Peer connection state is closed")
		}
	})
	p.LocalConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		logger.Debug("ICE connection state changed to: ", connectionState.String())
	})
	p.LocalConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		logger.Debug("ICE candidate found: ", candidate.String())
	})
}

func (p *Peer) onInvalidDataChannel(d *webrtc.DataChannel) {
	logger.Debug("Ignoring data channel: ", d.Label())
	err := d.Close()
	if err != nil {
		logger.Warn("Could not close data channel: ", err.Error())
		os.Exit(1) // TODO: Handle this better
	}

	err = p.LocalConnection.Close()
	if err != nil {
		logger.Warn("Could not close peer connection: ", err.Error())
		os.Exit(1) // TODO: Handle this better
	}
}

func (p *Peer) setupDataChannelProtocol() {
	// Set data channel protocol
	p.DataChannel.OnOpen(func() {
		logger.Debug("Data channel is open")

		candidatePair, err := p.LocalConnection.SCTP().Transport().ICETransport().GetSelectedCandidatePair()
		if err != nil {
			logger.Warn("Could not get selected candidate pair: ", err.Error())
			return
		}

		logger.Debug("Selected candidate pair: ", candidatePair.String())
		// TODO: Add on open channel protocol logic
	})
	p.DataChannel.OnClose(func() {
		logger.Debug("Data channel is closed")
		// TODO: Add on close channel protocol logic
	})
	p.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		logger.Debug("Message received: ", string(msg.Data))
		// TODO: Add message handling logic
	})

	// Set receiving data channel
	p.LocalConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		logger.Debug("Data channel received: ", d.Label())

		if d.Label() != DATA_CHANNEL_LABEL {
			p.onInvalidDataChannel(d)
			return
		}

		d.OnOpen(func() {
			logger.Debug("Accepted data channel. %s - %d\n", d.Label())
			// TODO: Add authentication logic
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			logger.Debug("Message received: ", string(msg.Data))
			// TODO: Add message handling logic
		})
	})
}

func (p *Peer) CreateOffer() (*webrtc.SessionDescription, error) {
	offer, err := p.LocalConnection.CreateOffer(nil)
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