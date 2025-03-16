package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Peykar server with plugins",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	logger := slog.Default()

	_ = logger

}

func init() {
	rootCmd.AddCommand(serveCmd)
}
