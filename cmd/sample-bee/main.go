package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "sample-bee",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(2)
		},
	}
	rootCmd.AddCommand(commandServe())
	rootCmd.AddCommand(commandVersion())
	return rootCmd
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	if err := commandRoot().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
