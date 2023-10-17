package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lats",
	Short: "Lats simplifies disaster recovery in AWS",
	Long: `Lats simplifies disaster recovery in AWS"
                Complete documentation is available at https://latscli.io/documentation/`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi From Lats")
	},
}

// Execute is the main function for cobra that we will run in our main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("error starting lats", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(CreateRDSSnapshotCmd)
	rootCmd.AddCommand(CopyRDSSnapshotCmd)
	rootCmd.AddCommand(RestoreRDSSnapshotCmd)
}
