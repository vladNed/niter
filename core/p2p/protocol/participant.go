package protocol

import (
	"context"
)


type ParticipantState struct {
	ctx          context.Context
	cancel       context.CancelFunc
	eventChannel chan SEvents
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

func NewParticipantState() *ParticipantState {
	ctx, cancel := context.WithCancel(context.Background())
	return &ParticipantState{
		ctx:          ctx,
		cancel:       cancel,
		eventChannel: make(chan SEvents),
	}
}
