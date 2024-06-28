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

# Download bindata if not already downloaded
go install -a -v github.com/go-bindata/go-bindata/...@latest

# Generate bindata assets
go-bindata -nomemcopy -pkg bindata -o ./gen/bindata/bindata.go ./install-scripts/...

# get latest slug from github actions env variable, or defaults to commit hash
if [ -z "$GITHUB_REF_NAME" ]; then
    head_describe=$(git describe --tags)
    meow_version=${GITHUB_REF_NAME:-${head_describe}}
else
    meow_version=${GITHUB_REF_NAME}
fi

echo "Meow version: ${meow_version}"

os_list=("linux" "windows")
arch_list=("amd64")

buildFor() {
    os="$1"
    arch="$2"
    extension=""

    if [ "${os}" == "windows" ]; then
        extension=".exe"
    fi

    real_meow_version=${meow_version}-${os}-${arch}
    GOOS=${os} GOARCH=${arch} go build \
        -trimpath \
        -ldflags "-s -w -X 'github.com/saniales/meow-cli/cmd.cliVersion=${real_meow_version}'" \
        -o "./bin/meow-${os}-${arch}${extension}"
}

for os in "${os_list[@]}"; do
    for arch in "${arch_list[@]}"; do
        echo "Compiling for ${os} ${arch}"
        buildFor "${os}" "${arch}"
    done
done

# Cross compile darwin executables
os="darwin"
arch_list=("amd64" "arm64")
for arch in "${arch_list[@]}"; do
    echo "Compiling for ${os} ${arch}"
    buildFor "${os}" "${arch}"
done 