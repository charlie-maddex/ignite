/*
Copyright © 2022 Charlie Maddex (charlie@multi.sh)
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ignite",
	Short: "A CLI for system administrators.",
	Long: `A CLI for system administors.
	This includes networking, storage, and Office 365 administration.
	The CLI is written in Go and uses Cobra for the CLI framework.
	It is designed to be cross-platform and run on Windows, Linux, and macOS.
	If you have any questions, please contact the developer at charlie@multi.sh.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ignite-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
