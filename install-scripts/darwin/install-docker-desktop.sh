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

# install docker desktop on mac checking apple silicon

if command -v docker &> /dev/null; then
    echo "Docker is already installed, skipping"
    exit 0
fi

if [ "$EUID" -ne 0 ]
  then echo "The install command requires root privileges, please run as root or with sudo"
  exit
fi

arch=$(uname -m)
echo "Detected architecture: $arch"

case $arch in
    x86_64)
        docker_install_url="https://desktop.docker.com/mac/main/amd64/Docker.dmg"
        ;;
    arm64)
        docker_install_url="https://desktop.docker.com/mac/main/arm64/Docker.dmg"
        ;;
    *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
esac

if [ -f "/tmp/Docker.dmg" ]; then
    echo "Docker Desktop installer already downloaded"
else    
    echo "Downloading Docker Desktop installer..."
    curl -L "$docker_install_url" -o /tmp/Docker.dmg
fi

echo "Installing Docker Desktop..."

hdiutil attach Docker.dmg
/Volumes/Docker/Docker.app/Contents/MacOS/install
hdiutil detach /Volumes/Docker

echo "Docker Desktop installed"