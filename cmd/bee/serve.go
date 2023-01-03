package main

import (
	"context"
	"errors"
	"fmt"
	server2 "github.com/luminosita/bee/common/server"
	"github.com/luminosita/bee/common/server/adapters"
	"github.com/luminosita/bee/internal/bee"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

type ServerOptions struct {
	// Flags
	Env     server2.Environment
	BaseUrl string
}

func commandServe() *cobra.Command {
	options := ServerOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] environment",
		Short:   "Launch Bee",
		Example: "bee serve environment",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			options.Env = server2.EnvironmentFromString(args[0])

			return runServe(&options, cmd.Flags())
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.BaseUrl, "baseUrl", "", "Base URL")

	return cmd
}

func runServe(options *ServerOptions, pflags *pflag.FlagSet) error {
	ctx := context.Background()

	boot := rkboot.NewBoot()
	boot.Bootstrap(ctx)

	viper, err := setupViper(options, pflags)
	if err != nil {
		return err
	}

	server := adapters.NewFiberServerTemplate(options.Env, bee.NewBeeServer(&bee.Config{}))

	server.Run(ctx, viper)

	boot.WaitForShutdownSig(ctx)

	return nil
}

func setupViper(options *ServerOptions, pflags *pflag.FlagSet) (*viper.Viper, error) {
	viper := rkentry.GlobalAppCtx.GetConfigEntry(fmt.Sprintf("%s-config", options.Env))

	if viper == nil {
		//TODO: Externalize
		return nil, errors.New("Unable to load connfiguration. Check the configuration file path")
	}

	if options.BaseUrl != "" {
		viper.BindPFlag("config.server.baseUrl", pflags.Lookup("baseUrl"))
	}

	//TOOD: Environnment vars are not working
	viper.SetEnvPrefix("bee") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return viper.Viper, nil
}
