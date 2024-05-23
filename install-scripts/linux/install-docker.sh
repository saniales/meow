#!/bin/bash

# Copyright © 2024 Alessandro Sanino <alessandro@sanino.dev>

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.

if command -v docker &> /dev/null; then
    echo "Docker is already installed"
    exit 0
fi

if [ "$EUID" -ne 0 ]
  then echo "The install command requires root privileges, please run as root or with sudo"
  exit
fi

$docker_install_flags=""
while [ $# -gt 0 ]; do
	case "$1" in
		--dry-run)
            ;;&
		--version)
            $docker_install_flags="$docker_install_flags $1"
            shift
            ;;
		*)
			echo "Illegal option $1"
			;;
	esac
	shift $(( $# > 0 ? 1 : 0 ))
done

echo "Installing docker..."

# import the flags into docker install script
curl -fsSL https://get.docker.com -o get-docker.sh $docker_install_flags
sh get-docker.sh
rm get-docker.sh

exit 0
