package mbit

import (
	"context"
	"os"

	"github.com/marklude/mbit/exchange"
	"github.com/marklude/mbit/feed"
	"github.com/marklude/mbit/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start trading...",
	Long:  "Start trading using config...",
	Run:   executeStartCmd,
}

func init() {
	RootCmd.AddCommand(startCmd)
}

func executeStartCmd(cmd *cobra.Command, args []string) {

	l := logger.NewLogger(os.Getenv("ENV"), os.Getenv("ELK"))
	arbitrage1, err := exchange.NewBinanceWSClient(context.Background(), os.Getenv("BinanceApiKey"), os.Getenv("BinanceSecretKey"))
	if err != nil {
		errors.Wrap(err, "New binance websocket client")
	}
	feeder := feed.NewFeed(context.Background(), feed.WithLogger(l), feed.WithBinance(arbitrage1))
	feeder.Start()
}
