package websocket

import (
	"math"
	"time"
)

const (
	infinity = time.Duration(math.MaxInt64)
	defaultMaxConnectionIdle = infinity
	defaultAckTimeout        = 30 * time.Second
)
