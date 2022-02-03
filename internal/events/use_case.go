package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/pkg/logging"
)

type uc struct {
	logger          *logging.Entry
	repo            eventEntities.Repository
	sub             message.Subscriber
	mx              sync.Mutex
	buffer          []*eventEntities.Event
	cfg             Config
	autoclearTicker *time.Ticker
}

func NewUseCase(logger *logging.Entry, repo eventEntities.Repository, sub message.Subscriber, config Config) *uc {

	return &uc{
		logger:          logger,
		repo:            repo,
		sub:             sub,
		mx:              sync.Mutex{},
		buffer:          []*eventEntities.Event{},
		cfg:             config,
		autoclearTicker: time.NewTicker(config.BufferAutoClearDuration),
	}
}

func (u *uc) Run(ctx context.Context) error {
	messages, err := u.sub.Subscribe(ctx, u.cfg.EventsTopic)
	if err != nil {
		return errors.WithMessage(err, "subscribing to messages")
	}

	go u.listener(ctx, messages)
	return nil
}

func (u *uc) listener(ctx context.Context, msgs <-chan *message.Message) {
	logger := u.logger.WithMethod("listener")
	for {
		select {
		case _ = <-ctx.Done():
			u.autoclearTicker.Stop()
			return
		case msg := <-msgs:
			u.handleMessage(msg)
		case _ = <-u.autoclearTicker.C:
			u.mx.Lock()
			logger.WithPlace("autoclear").Info("Autoclear time reached")
			err := u.clearBuffer()
			if err != nil {
				logger.WithPlace("autoclear").Error(err)
			}
			u.mx.Unlock()
		}
	}
}

func (u *uc) handleMessage(msg *message.Message) {
	u.mx.Lock()
	defer u.mx.Unlock()

	logger := u.logger.WithMethod("handleMessage")

	rawEvent := eventEntities.RawEvent{}

	if err := json.Unmarshal(msg.Payload, &rawEvent); err != nil {
		logger.WithPlace("read_message").Error(err)
		msg.Ack()
		return
	}

	event, err := RawToEvent(rawEvent)
	if err != nil {
		logger.WithPlace("validate_message").Error(err)
		msg.Ack()
		return
	}

	u.buffer = append(u.buffer, event)
	u.resetAutoclearTicker()
	if len(u.buffer) < u.cfg.BufferSize {
		msg.Ack()
		return
	}

	logger.Info("Buffer max len reached")
	err = u.clearBuffer()
	if err != nil {
		logger.WithPlace("clearBuffer").Error(err)
	}
	msg.Ack()

	return

}

func (u *uc) resetAutoclearTicker() {
	u.autoclearTicker.Reset(u.cfg.BufferAutoClearDuration)
}

func (u *uc) clearBuffer() error {
	logger := u.logger.WithMethod("clearBuffer")
	defer u.resetAutoclearTicker()
	logger.Info("Starting...")
	defer logger.Info("Finish...")

	bufferLen := int64(len(u.buffer))
	if bufferLen == 0 {
		logger.Info("buffer is empty, nothing to insert")
		return nil
	}
	logger.Infof("current buffer len: %d cap %d", len(u.buffer), cap(u.buffer))
	stored, err := u.repo.StoreMany(u.buffer...)
	if err != nil {
		return errors.WithMessage(err, "storing into db")
	}

	if stored == 0 {
		return errors.New("nothing stored")
	}

	if stored < bufferLen {
		logger.Warnf("could store %d/%d events", stored, bufferLen)
	}
	u.buffer = u.buffer[:0]
	return nil
}
