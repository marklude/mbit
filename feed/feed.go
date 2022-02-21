package feed

import (
	"context"

	"github.com/marklude/mbit/datastore"
	"github.com/marklude/mbit/exchange"
	"github.com/marklude/mbit/logger"
)

type feed struct {
	ctx     context.Context
	logger  logger.Logger
	binance exchange.Binance
	redis   datastore.Redis
}

type options func(*feed)

func NewFeed(ctx context.Context, opts ...options) *feed {
	fd := &feed{ctx: ctx}
	fd.WithOptions(opts...)
	return fd
}

func (fd *feed) WithOptions(opts ...options) {
	for _, opt := range opts {
		opt(fd)
	}
}

func WithLogger(l logger.Logger) options {
	return func(f *feed) {
		f.logger = l
	}
}

func WithBinance(b exchange.Binance) options {
	return func(f *feed) {
		f.binance = b
	}
}

func WithRedis(r datastore.Redis) options {
	return func(f *feed) {
		f.redis = r
	}
}

func (fd *feed) Start() {
	fd.logger.Info("Starting feeder...")
	defer fd.logger.Info("Stopping feeder...")

	//init feeders
	go fd.binance.GetPrice()

	<-fd.ctx.Done()

}
