package websocket

import "time"

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication
	ack        AckType
	ackTimeout time.Duration
	concurrency int
	pattern string
	maxConnectionIdle time.Duration
}

func newServerOptions(opts ...ServerOptions) serverOption {
	o := serverOption{
		Authentication: new(authentication),
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		pattern:        "/ws",
		concurrency:     defaultConcurrency,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
func WithAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}
func WithHandlerPattern(pattern string) ServerOptions {
	return func(opt *serverOption) {
		opt.pattern = pattern
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}
