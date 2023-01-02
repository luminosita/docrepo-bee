package main

import (
	"context"
	"fmt"
	"github.com/luminosita/bee/common/server/adapters"
	"github.com/luminosita/bee/internal/bee"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type serveOptions struct {
	// Flags
	baseUrl string
	version string
}

func commandServe() *cobra.Command {
	options := serveOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] config-file",
		Short:   "Launch Bee",
		Example: "bee serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			flags := cmd.Flags()

			err := flags.Parse(args)

			if err != nil {
				return err
			}

			//TODO : Not working for flags and environment
			viper.BindPFlag("server.version", flags.Lookup("version"))
			viper.BindPFlag("server.baseUrl", flags.Lookup("baseUrl"))

			viper.SetEnvPrefix("bee") // will be uppercased automatically
			viper.AutomaticEnv()

			return runServe()
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.baseUrl, "baseUrl", "", "Base URL")
	flags.StringVar(&options.version, "version", "", "Version")

	return cmd
}

func runServe() error {
	server, err := adapters.NewFiberServerAdapter(context.Background(), bee.NewBeeServer())
	if err != nil {
		return fmt.Errorf("failed to initialize server: %v", err)
	}

	return server.Run()
}
