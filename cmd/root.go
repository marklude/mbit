package mbit

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

const (
	versionNumber = "0.0.1-pre-alpha"
)

var RootCmd = &cobra.Command{
	Use:   "mbit",
	Short: fmt.Sprintf("USAGE %s [OPTIONS]", os.Args[0]),
	Long:  fmt.Sprintf("USAGE %s [OPTIONS]: see --help for more details", os.Args[0]),
	Run:   executeRootCmd,
}

var signals chan os.Signal

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		signal.Stop(signals)
		fmt.Println()
		fmt.Println("CTRL-C received. Exiting...")
		os.Exit(0)
	}()

	RootCmd.Flags().BoolVarP(&rootFlags.Version, "version", "v", false, "shows version information")

}

func executeRootCmd(cmd *cobra.Command, args []string) {
	if rootFlags.Version {
		fmt.Printf("mBit Arbitrage Bot v.%s\n", versionNumber)
	} else {
		cmd.Help()
	}
}
