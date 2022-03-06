package exchange

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/marklude/logger"
)

const (
	websocketLimit = 5
	symbol         = "BNBUSDT"
)

type binanceWrapper struct {
	ctx    context.Context
	api    *binance.Client
	logger logger.Logger
}

type Binance interface {
	GetPrice() (float64, error)
}

func NewBinanceWSClient(ctx context.Context, apiKey, secretKey string) (Binance, error) {
	binance.UseTestnet = true
	client := binance.NewClient(apiKey, secretKey)
	client.NewServerTimeService().Do(context.Background())
	return &binanceWrapper{
		ctx: ctx,
		api: client,
	}, nil
}

func (b *binanceWrapper) GetPrice() (price float64, err error) {
	var kline []*binance.Kline
	ticker := time.NewTicker(3 * time.Minute)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				fmt.Print("done")
				return
			case <-ticker.C:
				kline, err = b.api.NewKlinesService().Symbol(symbol).Interval("3m").Do(b.ctx)
				if err != nil {
					b.logger.Error(err)

				}

				fmt.Printf("BNBUSDT - %s\n", kline[len(kline)-1].Close)

			}
		}
	}()

	// p, err := strconv.ParseFloat(kline[0].Close, 64)
	// if err != nil {
	// 	errors.Wrap(err, "[Get Price] Parse float error")
	// }
	//done <- true
	return 0.0, nil

}
