package apm

import (
	"fmt"
	"os"

	"bitbucket.org/pharmaeasyteam/tokenizer/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/lib/apm/appdynamics"
	"github.com/google/uuid"
)

// InitAppDynamics sets the configs from vault and env and starts AppDynamics
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
	appDconfig.InitTimeoutMs = 1000

	// init the SDK
	if err := appdynamics.InitSDK(&appDconfig); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}
}

// StartBT starts the AppD Business transaction, creates a UUID
// and stores it for future reference
func StartBT(transaction string) string {
	bt := appdynamics.StartBT(transaction, "")
	txUUID := uuid.New().String()
	appdynamics.StoreBT(bt, txUUID)

	return txUUID
}

// LogFatalBTError logs a business transaction error and ends the transaction
func LogFatalBTError(err error, appDTxnId string) {
	bt := appdynamics.GetBT(appDTxnId)
	appdynamics.AddBTError(bt, appdynamics.APPD_LEVEL_ERROR, err.Error(), true)
	appdynamics.EndBT(bt)
}

// LogFatalBTErrorMessage is the same as LogFatalBTError but it accepts the error as string
func LogFatalBTErrorMessage(err string, appDTxnId string) {
	bt := appdynamics.GetBT(appDTxnId)
	appdynamics.AddBTError(bt, appdynamics.APPD_LEVEL_ERROR, err, true)
	appdynamics.EndBT(bt)
}

// EndBT fetches the stored business transaction and ends it
func EndBT(appDTxnId string) {
	bt := appdynamics.GetBT(appDTxnId)
	appdynamics.EndBT(bt)
}
