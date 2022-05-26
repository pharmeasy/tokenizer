package apm

import (
	"fmt"
	"time"

	"bitbucket.org/pharmaeasyteam/tokenizer/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/lib/apm/appdynamics"
	"github.com/sirupsen/logrus"
)

func InitAppDynamics(cfg *config.TokenizerConfig) {

	vaultConfig := cfg.VaultModule.AppDynamicsConfig

	// Configure AppD
	appDconfig := appdynamics.Config{}

	// Controller
	appDconfig.Controller.Host = vaultConfig.AppDynamicsConfigController.Host
	appDconfig.Controller.Port = vaultConfig.AppDynamicsConfigController.Port
	appDconfig.Controller.UseSSL = vaultConfig.AppDynamicsConfigController.UseSSL
	appDconfig.Controller.Account = vaultConfig.AppDynamicsConfigController.Account
	appDconfig.Controller.AccessKey = vaultConfig.AppDynamicsConfigController.AccessKey

	// App Context
	appDconfig.AppName = vaultConfig.AppName
	appDconfig.TierName = vaultConfig.TierName
	appDconfig.NodeName = vaultConfig.NodeName

	// misc
	appDconfig.InitTimeoutMs = 1000

	// init the SDK
	if err := appdynamics.InitSDK(&appDconfig); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")

		logrus.Info("Initialized AppDynamics successfully\n")

		// Run some BTs
		maxBtCount := 20000
		btCount := 0

		fmt.Print("Doing something")
		for btCount < maxBtCount {
			// start the "Checkout" transaction
			btHandle := appdynamics.StartBT("MyTestGolangBT", "")

			// do something....
			fmt.Print(".")
			milliseconds := 250
			time.Sleep(time.Duration(milliseconds) * time.Millisecond)

			// end the transaction
			appdynamics.EndBT(btHandle)

		}
		fmt.Print("\n")

		// Stop/Clean up the AppD SDK.
		appdynamics.TerminateSDK()
	}
}
