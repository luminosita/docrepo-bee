package main

import (
	"context"
	"fmt"
	"github.com/luminosita/bee/common/server/adapters"
	"github.com/luminosita/bee/internal/bee"
	"github.com/spf13/cobra"
)

type serveOptions struct {
	// Flags
	webHTTPAddr  string
	webHTTPSAddr string
}

func commandServe() *cobra.Command {
	options := serveOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] config-file",
		Short:   "Launch Bee",
		Example: "bee serve",
		//		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			//						options.config = args[0]

			return runServe()
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.webHTTPAddr, "web-http-addr", "", "Web HTTP address")
	flags.StringVar(&options.webHTTPSAddr, "web-https-addr", "", "Web HTTPS address")

	return cmd
}

func runServe() error {
	server, err := adapters.NewFiberServerAdapter(context.Background(), bee.NewBeeServer())
	if err != nil {
		return fmt.Errorf("failed to initialize server: %v", err)
	}

	return server.Run()
}
