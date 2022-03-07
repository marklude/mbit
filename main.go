package main

import (
	mbit "github.com/marklude/mbit/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		errors.Wrap(err, "Error reading config file")
	}

	mbit.Execute()
}
