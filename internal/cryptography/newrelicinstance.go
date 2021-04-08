package cryptography

import (
	"os"
	"strings"

	"bitbucket.org/pharmaeasyteam/goframework/config"
	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/metrics"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

func setNewRelic(globalConfig *config.ServerConfig) {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	napp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(globalConfig.ServiceName),
		newrelic.ConfigLicense(globalConfig.NewRelic.LicenseKey),
		func(cfg *newrelic.Config) {
			cfg.Enabled = true
			cfg.Labels["hostname"] = hostname
			cfg.Labels["env"] = strings.ToLower(globalConfig.Profile)
			cfg.Labels["region"] = strings.ToLower(globalConfig.Region)
		},
	)
	if err != nil {
		logging.GetLogger().Error("Failed to initialize newrelic ", zap.Error(err))
	} else {
		logging.GetLogger().Info("Setting newrelic app")
		logging.GetLogger().Info("Setting newrelic app", zap.String("service_name", globalConfig.ServiceName))
		metrics.SetNewRelicApp(napp)
	}
}

/*
os.Setenv("SERVER_NEWRELIC_ENABLED", "true")
	os.Setenv("SERVER_NEWRELIC_LICENSEKEY", "b451e7f2a22397867187fb27b2c127931b67e44b")
*/
