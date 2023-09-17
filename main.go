package main

import (
	instana "github.com/instana/go-sensor"
	"github.com/pharmaeasy/tokenizer/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	//instana initalization
	instana.InitSensor(&instana.Options{
		Service:           "tokenizer",
		LogLevel:          instana.Debug,
		EnableAutoProfile: true,
	})

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Error(err.Error())
		os.Exit(1)
	}
}
