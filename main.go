package main

import (
	"bitbucket.org/pharmaeasyteam/tokenizer/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Error(err.Error())
		os.Exit(1)
	}
}
