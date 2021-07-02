package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A CLI",
	Long:  "A CLI to execute arbitrary Linux commands",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(queryStatusCmd)
	rootCmd.AddCommand(startServerCmd)
	startCmd.Flags().StringP("dir", "d", "./internal/api/certs", "Provide a directory for your client certificates")
	stopCmd.Flags().StringP("dir", "d", "./internal/api/certs", "Provide a directory for your client certificates")
	logCmd.Flags().StringP("dir", "d", "./internal/api/certs", "Provide a directory for your client certificates")
	queryStatusCmd.Flags().StringP("dir", "d", "./internal/api/certs", "Provide a directory for your client certificates")
}
