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

# Generate bindata assets
go-bindata -nomemcopy -pkg bindata -o ./gen/bindata/bindata.go ./install-scripts/...

# get latest slug from github actions env variable, or defaults to commit hash
MEOW_VERSION=${GITHUB_REF_NAME:-$(git rev-parse HEAD)}

OS_LIST=("linux" "windows")
ARCH_LIST=("386" "amd64" "arm64")

for OS in "${OS_LIST[@]}"; do
    for ARCH in "${ARCH_LIST[@]}"; do
        echo "Compiling for ${OS} ${ARCH}"
        GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-s -w" -o ./bin/meow-${OS}-${ARCH} .
    done
done

# Cross compile darwin executables
ARCH_LIST=("amd64" "arm64")
for ARCH in "${ARCH_LIST[@]}"; do
    echo "Compiling for darwin ${ARCH}"
    GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-s -w" -o ./bin/meow-darwin-${ARCH} .
done 