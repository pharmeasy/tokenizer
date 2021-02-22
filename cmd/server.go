package cmd

import (
	"context"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"bitbucket.org/pharmaeasyteam/goframework/server"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/config"
	"bitbucket.org/pharmaeasyteam/tokenizer/internal/cryptography"
	"github.com/spf13/cobra"
)

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts Tokenizer App",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := config.Load(configFile)
			if err != nil {
				logging.GetLogger().Error("Could not load configurations from file ")
			}

			//logger.Info("Final configuration", zap.Any("config", &globalConfig))
			return RunServerStart(ctx, cfg)
		},
	}

	return cmd
}

//RunServerStart run server
func RunServerStart(ctx context.Context, cfg *config.TokenizerConfig) error {

	svr := server.New(
		server.WithGlobalConfig(&cfg.Server),
	)

	// Add the crypto module
	svr.AddModule("crypto", &cryptography.ModuleCrypto{
		LoadModule: cfg.LoadModule,
	})

	svr.Start(ctx)
	logging.GetLogger().Info("Shutting down")

	return nil
}
