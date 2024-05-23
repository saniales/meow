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

if ! command -v docker &> /dev/null
then
    echo "Docker command not found, maybe you need to install docker first?"
    echo "You can do it with the 'meow install' command"
    exit
fi

cat_image_url="ghcr.io/cheshire-cat-ai/core"
cat_image_version="latest"
cat_port=1865
plugins_folder="./plugins"
data_folder="./data"
static_folder="./static"
docker_image_name="cheshire-cat-ai"

while [ $# -gt 0 ]; do
	case "$1" in
        --cat-image-url)
            shift
            cat_image_url="$1"
            shift
            ;;
		--cat-image-version)
            shift
            cat_image_version="$1"
            shift
            ;;
        --cat-port)
            shift
            cat_port="$1"
            shift
            ;;
        --plugins-folder)
            shift
            plugins_folder="$1"
            shift
            ;;
        --data-folder)
            shift
            data_folder="$1"
            shift
            ;;
        --static-folder)
            shift
            static_folder="$1"
            shift
            ;;
        --docker-image-name)
            shift
            docker_image_name="$1"
            shift
            ;;
		*)
			echo "Illegal option $1"
			;;
	esac
	shift $(( $# > 0 ? 1 : 0 ))
done

cat_image_full_url="$cat_image_url:$cat_image_version"
docker run --name $docker_image_name \
    -d \
    --rm \
    -p "$cat_port:80" \
    -v "$plugins_folder:/app/cat/plugins" \
    -v "$data_folder:/app/cat/data" \
    -v "$static_folder:/app/cat/static" \
    "$cat_image_full_url"

if [ $? -ne 0 ]; then
    echo "Failed to run docker container with cat"
    exit -1
fi