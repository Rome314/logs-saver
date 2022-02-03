package events

import (
	"time"
)

type Config struct {
	EventsTopic             string
	BufferSize              int
	BufferAutoClearDuration time.Duration
}
