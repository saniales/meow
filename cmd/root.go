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

// Package cmd represents the commands of the Meow CLI
package cmd

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/saniales/meow-cli/pkg/providers/install"
)

var cfgFile string

var cliVersion string = "0.0.0-unstable"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meow",
	Short: "The paw-friendly Cheshire Cat AI Command Line Interface",
	Long: `The paw-friendly Cheshire Cat AI Command Line Interface.
	
The Cheshire Cat is an open-source, hackable and production-ready framework that
allows developing intelligent personal AI assistant agents on top of Large Language Models (LLMs).

You can find more information at https://cheshirecat.ai and https://cheshire-cat-ai.github.io/docs.

You can use this paw-friendly CLI to manage the cat installs, call the API and more!`,
	Example:           "meow help",
	Version:           cliVersion,
	PersistentPreRunE: globalPreRunE,
}

var globalFlags struct {
	configFile string
	verbose    bool
	quiet      bool
	json       bool
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of meow",
	Long:  `Print the version number of meow`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("meow version %s\n", cliVersion)
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the cat and its dependencies on the current machine",
	Long:  `Installs the cat and its dependencies on the current machine`,
	Run:   executeInstall,
}

var installCmdFlags struct {
	reinstall bool
	dryRun    bool
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
	rootCmd.AddCommand(versionCmd)

	// global flags
	rootCmd.PersistentFlags().StringVar(&globalFlags.configFile, "config", "", "config file (default is $HOME/.meow-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Enable verbose output (default is false) - Incompatible with --quiet")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.quiet, "quiet", "q", false, "Enables output only on errors (default is false) - Incompatible with --verbose")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.json, "json", false, "Enables JSON formatted output (default is false)")

	// install flags
	installCmd.Flags().BoolVar(&installCmdFlags.reinstall, "reinstall", false, "Force install even if (default is false)")
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

	viper.SetEnvPrefix("CCAT")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func globalPreRunE(cmd *cobra.Command, args []string) error {
	// TODO: check for updates for the CLI from github releases

	if globalFlags.verbose && globalFlags.quiet {
		return errors.New("--verbose and --quiet flags are incompatible")
	}

	defaultLogHandlerOptions := new(slog.HandlerOptions)
	if globalFlags.verbose {
		defaultLogHandlerOptions.Level = slog.LevelDebug
	} else if globalFlags.quiet {
		defaultLogHandlerOptions.Level = slog.LevelError
	} else {
		defaultLogHandlerOptions.Level = slog.LevelInfo
	}

	var defaultLogHandler slog.Handler
	if globalFlags.json {
		defaultLogHandler = slog.NewJSONHandler(os.Stdout, defaultLogHandlerOptions)
	} else {
		defaultLogHandler = slog.NewTextHandler(os.Stdout, defaultLogHandlerOptions)
	}
	defaultLogger := slog.New(defaultLogHandler)
	slog.SetDefault(defaultLogger)

	return nil
}

// executeInstall performs the "install" logic.
func executeInstall(cmd *cobra.Command, args []string) {
	downloadProgressFunc := progressbarFunc

	httpClient := new(http.Client)
	installer, err := install.NewInstallerWithProgressFunc(httpClient, downloadProgressFunc)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	switch runtime.GOOS {
	case "linux":
		slog.Info(
			"Downloading and installing Docker Engine...",
			slog.String("os", runtime.GOOS),
			slog.String("arch", runtime.GOARCH),
		)
		err := installer.InstallDocker(install.InstallDockerConfig{
			Verbose: globalFlags.verbose,
		})
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	case "darwin":
	case "windows":
		slog.Info(
			"Downloading Docker Desktop installer...",
			slog.String("os", runtime.GOOS),
			slog.String("arch", runtime.GOARCH),
		)
		err := installer.DownloadDockerDesktopInstaller(install.DownloadDockerDesktopInstallerConfig{
			ForceDownload: installCmdFlags.reinstall,
		})
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		slog.Info(
			"Running Docker Desktop installer...",
			slog.String("os", runtime.GOOS),
			slog.String("arch", runtime.GOARCH),
		)
		err = installer.RunDockerDesktopInstaller(install.RunDockerDesktopInstallerConfig{
			Verbose: globalFlags.verbose,
		})
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}

func progressbarFunc(total int64, reader io.Reader) {
	progressBar := pb.Full.Start64(total)
	pbProxy := progressBar.NewProxyReader(reader)
	io.Copy(io.Discard, pbProxy)
	progressBar.Finish()
}

func ttyFunc(total int64, reader io.Reader) {
	var totalRead int64
	buffer := make([]byte, total/50)
	start := 0
	end := 50
	for totalRead <= total {
		n, _ := reader.Read(buffer)
		percentage := float64(totalRead) / float64(total) * 100

		progressTTY := strings.Repeat("#", start) + strings.Repeat(".", end)
		progressMessage := fmt.Sprintf("%s - %d%% completed", progressTTY, int(percentage))
		totalRead += int64(n)

		fmt.Printf("\r%s", progressMessage)

		start++
		end--
	}

	fmt.Println(strings.Repeat("#", 50))
}
