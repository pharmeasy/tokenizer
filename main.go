package main

import (
	"bitbucket.org/pharmaeasyteam/tokenizer/cmd"
	instana "github.com/instana/go-sensor"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	instana.InitSensor(&instana.Options{
		Service:           "tokenizer-gateway",
		LogLevel:          instana.Debug,
		EnableAutoProfile: true,
	})
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Error(err.Error())
		os.Exit(1)
	}
}
