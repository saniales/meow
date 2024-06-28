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

// Package constants contains the OS-mapped constants definitions.
package constants

import (
	"fmt"
	"os"
	"runtime"
)

var (
	constants map[string]map[string]map[string]string = map[string]map[string]map[string]string{
		"windows": {
			"amd64": {
				"docker_desktop_installer_url":  "https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe",
				"default_docker_installer_path": fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, "docker-desktop-installer.exe"),
			},
		},
		"linux": {
			"amd64": {
				"docker_installer_url":          "",
				"default_docker_installer_path": fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, "docker-desktop-installer"),
			},
		},
		"darwin": {
			"amd64": {
				"docker_desktop_installer_url":  "https://desktop.docker.com/mac/main/amd64/Docker.dmg",
				"default_docker_installer_path": fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, "docker-desktop-installer"),
			},
			"arm64": {
				"docker_desktop_installer_url":  "https://desktop.docker.com/mac/main/arm64/Docker.dmg",
				"default_docker_installer_path": fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, "docker-desktop-installer"),
			},
		},
	}
)

// GetConstant retrieves the value of a constant based on its name.
//
// Parameters:
// - name: the name of the constant to retrieve.
//
// Returns:
// - string: the value of the constant.
// - error: an error if the constant is not found or the operating system and architecture are not supported.
func GetConstant(name string) (string, error) {
	definedConstants, exists := constants[runtime.GOOS][runtime.GOARCH]
	if !exists {
		return "", ErrUnsupportedOSArch
	}

	constant, exists := definedConstants[name]
	if !exists {
		return "", ErrConstantNotFound
	}

	return constant, nil
}
