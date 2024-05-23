/*
Copyright © 2024 Alessandro Sanino <alessandro@sanino.dev>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meow",
	Short: "The paw-friendly Cheshire Cat AI Command Line Interface",
	Long: `The paw-friendly Cheshire Cat AI Command Line Interface.
	
The Cheshire Cat is an open-source, hackable and production-ready framework that
allows developing intelligent personal AI assistant agents on top of Large Language Models (LLMs).

You can find more information at https://cheshirecat.ai and https://cheshire-cat-ai.github.io/docs.

You can use this paw-friendly CLI to manage the cat installs, call the API and more!`,
	Example: "meow help",
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the cat and its dependencies on the current machine",
	Long:  `Installs the cat and its dependencies on the current machine`,
	RunE:  executeInstall,
}

var installCmdFlags struct {
	dryRun bool
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(installCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.meow-cli.yaml)")

	// install flags
	installCmd.Flags().BoolVar(&installCmdFlags.dryRun, "dry-run", false, "Print only the install steps without performing them (default is false)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".meow-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".meow-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// executeInstall performs the "install" logic.
func executeInstall(cmd *cobra.Command, args []string) error {
	return errors.New("command not yet implemented")
}
