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

param (
    [string]$DockerDesktopPath = "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    [string]$CatImageURL = "ghcr.io/cheshire-cat-ai/core"
    [string]$CatImageVersion = "latest"
    [int]$CatPort = 1865
    [string]$PluginsFolder = ".\cat-plugins"
    [string]$DataFolder = ".\cat-data"
    [string]$StaticFolder = ".\cat-static"
    [string]$DockerImageName = "cheshire-cat-ai"
)

$IsDockerRunning = Get-Process 'com.docker.proxy'
if ( !$isDockerRunning )
{
    Write-Progress -Activity "Starting Docker Desktop..."

    Start-Process "$DockerDesktopPath"
    Start-Sleep -Seconds 30

    Write-Progress -Activity "Starting Docker Desktop..." -Completed

    Write-Output "Check if Docker Desktop is running..."
    $IsDockerRunning = Get-Process 'com.docker.proxy'

    if ( !$isDockerRunning )
    {
        Write-Error "Docker Desktop did not start correctly, cannot start the cat"
        return -1
    }
}

# Pull the latest cat image
$CatImageFullURL = "$CatImageURL:$CatImageVersion"
docker pull "$CatImageFullURL"

# Run the cat image
$IsDockerRunSuccessful = docker run \
    --name "cheshire-cat-ai" \
    -d \
    -p "$CatPort:80" \
    -v "$PluginsFolder:/app/cat/plugins" \
    -v "$DataFolder:/app/cat/data" \
    -v "$StaticFolder:/app/cat/static" \
    "$CatImageFullURL"

if ( !$IsDockerRunSuccessful )
{
    Write-Error "Docker run failed"
    return -2
}