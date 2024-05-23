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
    [string]$DockerDesktopInstallerURL = "https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe"
    [string]$DockerDesktopInstallerPath = Join-Path $Env:Temp DockerDesktopInstaller.exe
    [string]$DockerDesktopPath = "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    [boolean]$ForceInstall = false
)

$IsDockerInstalled = Test-Path $DockerDesktopPath -PathType Leaf
if ( $IsDockerInstalled && !$ForceInstall )
{
    Write-Output "Docker Desktop already installed, skipping"
    return 0
}

$DockerDesktopInstallerAlreadyDownloaded = Test-Path $DockerDesktopInstallerPath -PathType Leaf

if ( !$DockerDesktopInstallerAlreadyDownloaded )
{
    Write-Output "Docker Desktop for Windows installer already downloaded"
}
else 
{
    ## Download the installer
    Start-BitsTransfer -DisplayName "Downloading Docker Desktop for Windows installer" -Source $DockerDesktopInstallerURL -Destination $DockerDesktopInstallerPath
    Write-Output "Docker Desktop for Windows installer downloaded"
}

## Run the installer
Write-Output "Installing Docker Desktop for Windows..."
Start-Process "$DockerDesktopInstallerPath" -Wait "install --quiet --accept-license"
Write-Output "Docker Desktop for Windows installed"