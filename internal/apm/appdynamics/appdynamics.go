package apm

import (
	"fmt"
	"os"
	"time"

	"bitbucket.org/pharmaeasyteam/tokenizer/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/lib/apm/appdynamics"
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
	appDconfig.Controller.AccessKey = os.Getenv("APPDYNAMICS_AGENT_ACCOUNT_ACCESS_KEY")

	// App Context
	appDconfig.AppName = vaultConfig.AppName
	appDconfig.TierName = vaultConfig.TierName
	appDconfig.NodeName = os.Getenv("APPDYNAMICS_AGENT_NODE_NAME")

	// misc
	appDconfig.InitTimeoutMs = 1000

	// init the SDK
	if err := appdynamics.InitSDK(&appDconfig); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}
}

func RunBTs() {
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
