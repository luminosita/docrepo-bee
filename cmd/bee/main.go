package bee

import (
	"fmt"
	"github.com/luminosita/docrepo-bee/internal/bee"
	"github.com/luminosita/honeycomb/pkg/cmd"
	"os"

	"github.com/spf13/cobra"
)

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "bee",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(2)
		},
	}
	rootCmd.AddCommand(cmd.CommandServe(bee.NewBeeServer(&bee.Config{})))
	rootCmd.AddCommand(cmd.CommandVersion())
	return rootCmd
}

func main() {
	if err := commandRoot().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
