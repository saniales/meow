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
    [string]$DockerImageName = "cheshire-cat-ai"
)

$IsDockerRunning = Get-Process 'com.docker.proxy'
if ( !$isDockerRunning )
{
    Write-Progress -Activity "Docker Desktop not running, no need to stop the cat"
    return -1
}

# Stop the cat image
$IsDockerStopSuccessful = docker stop \
    --name "cheshire-cat-ai"

if ( !$IsDockerStopSuccessful )
{
    Write-Error "Docker stop failed"
    return -2
}