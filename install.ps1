# SNISID Platform One-Click Installer for Windows
# Version: 1.0.0

$ErrorActionPreference = "Stop"
Write-Host "Initializing SNISID National Identity Platform Installer..."

# 1. Administrator Check
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) { Write-Warning "Admin required."; exit }

# 2. Prerequisite Check (WSL2 & Docker)
Write-Host "Checking prerequisites..."

$wsl = Get-Command "wsl" -ErrorAction SilentlyContinue
if ($null -eq $wsl) { 
    Write-Host "Installing WSL2..."
    wsl --install --no-distribution
    Write-Host "WSL2 installed. Reboot required."
    exit 
}

$docker = Get-Command "docker" -ErrorAction SilentlyContinue
if ($null -eq $docker) { 
    Write-Host "Docker Desktop not found."
    exit 
}

# 3. K3D Installation
$k3d = Get-Command "k3d" -ErrorAction SilentlyContinue
if ($null -eq $k3d) { 
    Write-Host "Installing k3d..."
    winget install k3d 
}

# 4. Bootstrap Cluster
Write-Host "Bootstrapping Local Cluster..."
if (Test-Path ".\scripts\bootstrap.ps1") { 
    powershell -ExecutionPolicy Bypass -File .\scripts\bootstrap.ps1 
}

# 5. Offline Mode Check
if (Test-Path ".\offline_images.tar") { 
    Write-Host "Found offline images."
    k3d image import .\offline_images.tar -c snisid 
}

# 6. Desktop Shortcut
try { 
    $WshShell = New-Object -ComObject WScript.Shell
    $Shortcut = $WshShell.CreateShortcut("$HOME\Desktop\SNISID Dashboard.lnk")
    $Shortcut.TargetPath = "http://localhost"
    $Shortcut.Save() 
} catch {}

Write-Host "Installation Complete!"
