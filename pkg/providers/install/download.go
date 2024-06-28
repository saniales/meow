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

// Package install contains the installer for the cat dependencies.
package install

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/saniales/meow-cli/gen/bindata"
	"github.com/saniales/meow-cli/pkg/constants"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type onDownloadProgressFunc func(total int64, reader io.Reader)

type Installer struct {
	httpClient         httpClient
	onDownloadProgress onDownloadProgressFunc
}

// NewInstallerWithProgressFunc creates a new Installer with the given httpClient and onDownloadProgressFunc.
//
// Parameters:
// - httpClient: The httpClient to be used by the Installer.
// - onDownloadProgress: The function to be called during the download progress (for example used to update UI).
//
// Returns:
// - *Installer: A pointer to the newly created Installer.
// - error: An error if the httpClient is nil.
func NewInstallerWithProgressFunc(
	httpClient httpClient,
	onDownloadProgress onDownloadProgressFunc,
) (*Installer, error) {
	if httpClient == nil {
		return nil, ErrNilHTTPClient
	}

	return &Installer{
		httpClient:         httpClient,
		onDownloadProgress: onDownloadProgress,
	}, nil
}

type DownloadDockerDesktopInstallerConfig struct {
	ForceDownload bool
}

// DownloadDockerDesktopInstaller downloads the Docker Desktop installer from the specified URL and saves it to a file in the OS temp directory.
//
// It returns the path to the temporary file and an error if any occurred.
func (i *Installer) DownloadDockerDesktopInstaller(config DownloadDockerDesktopInstallerConfig) error {
	installerFileName := "docker-desktop-installer"
	if runtime.GOOS == "windows" {
		installerFileName += ".exe"
	}

	installerPath := fmt.Sprintf("%s%cmeow-cli%c%s", os.TempDir(), os.PathSeparator, os.PathSeparator, installerFileName)
	if _, err := os.Stat(installerPath); !config.ForceDownload && err == nil {
		slog.Debug("Docker Desktop installer already exists", slog.String("path", installerPath))
		return nil
	}

	dockerDesktopInstallerURL, err := constants.GetConstant("docker_desktop_installer_url")
	if err != nil {
		return err
	}

	slog.Debug("Downloading Docker Desktop installer...", slog.String("url", dockerDesktopInstallerURL))
	req, err := http.NewRequest("GET", dockerDesktopInstallerURL, nil)
	if err != nil {
		return err
	}

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrNetwork(resp.StatusCode)
	}
	slog.Debug("DONE")

	slog.Debug("Saving Docker Desktop installer to temp dir...", slog.String("path", installerPath))
	tempFile, err := os.Create(installerPath + ".tmp")
	if err != nil {
		return err
	}

	var downloadWriter io.Writer = tempFile

	if i.onDownloadProgress != nil {
		var tempBuffer bytes.Buffer
		downloadWriter = io.MultiWriter(tempFile, &tempBuffer)

		i.onDownloadProgress(resp.ContentLength, &tempBuffer)
	}

	_, err = io.Copy(downloadWriter, resp.Body)
	if err != nil {
		return err
	}

	err = tempFile.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tempFile.Name(), installerPath)
	if err != nil {
		return err
	}
	slog.Debug("DONE", slog.String("path", installerPath))

	return nil
}

type RunDockerDesktopInstallerConfig struct {
	Verbose bool
}

// RunDockerDesktopInstaller installs Docker Desktop on the system.
//
// It takes a RunDockerDesktopInstallerConfig struct as a parameter, which contains the configuration for the installer.
// The config parameter has a single field, Verbose, which indicates whether the installer should output verbose logs.
//
// The function returns an error if there was a problem installing Docker Desktop, or nil if the installation was successful.
func (i *Installer) RunDockerDesktopInstaller(config RunDockerDesktopInstallerConfig) error {
	if runtime.GOOS == "linux" {
		return ErrDockerDesktopInstallNotSupported
	}

	installerPath, err := constants.GetConstant("docker_desktop_installer_path")
	if err != nil {
		return err
	}

	installerCmd := exec.Command(installerPath, "--quiet", "--accept-license")

	if config.Verbose {
		installerCmd.Stdout = os.Stdout
		installerCmd.Stderr = os.Stderr
	}

	err = installerCmd.Start()
	if err != nil {
		return err
	}

	err = installerCmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

type InstallDockerConfig struct {
	Verbose bool
}

func (i *Installer) InstallDocker(config InstallDockerConfig) error {
	if runtime.GOOS != "linux" {
		return ErrDockerInstallNotSupported
	}

	installDockerScript, err := bindata.Asset("install-scripts/linux/install-docker.sh")
	if err != nil {
		return err
	}

	installDockerCmd := exec.Command("bash", "-c", string(installDockerScript))

	if config.Verbose {
		installDockerCmd.Stdout = os.Stdout
		installDockerCmd.Stderr = os.Stderr
	}

	err = installDockerCmd.Start()
	if err != nil {
		return err
	}

	err = installDockerCmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
