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

package install

import (
	"fmt"
)

var (
	ErrNilHTTPClient                    = fmt.Errorf("nil HTTP client provided")
	ErrDockerInstallNotSupported        = fmt.Errorf("docker install not supported on this operating system, please install docker manually or perform the automatic docker desktop installation")
	ErrDockerDesktopInstallNotSupported = fmt.Errorf("docker desktop install not supported on linux, please install docker desktop manually or perform the automatic docker installation")
)

func ErrNetwork(statusCode int) error {
	return fmt.Errorf("request failed with status code %d", statusCode)
}
