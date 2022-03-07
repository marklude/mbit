package mbit

import (
	"context"
	"os"

	"github.com/marklude/logger"
	"github.com/marklude/mbit/exchange"
	"github.com/marklude/mbit/feed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	// initialize logger
	l := logger.NewLogger(os.Getenv("ENV"))

	// initialize new binance
	arbitrage1 := exchange.NewBinanceWSClient(context.Background(), os.Getenv("BinanceApiKey"), os.Getenv("BinanceSecretKey"))

	// initialize new cex
	cexBaseWsURL := viper.GetString("mbit.arbitrage2.websocket-api.endpoint")
	arbitrage2 := exchange.NewCex(os.Getenv("CexApiKey"), os.Getenv("CexSecretKey"), cexBaseWsURL)

	// initialize feed
	feeder := feed.NewFeed(context.Background(), feed.WithLogger(l), feed.WithBinance(arbitrage1), feed.WithCex(arbitrage2))
	feeder.Start()
}
