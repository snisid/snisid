param(
    [string]$WinPEKit = "C:\Program Files (x86)\Windows Kits\10\Assessment and Deployment Kit\Windows Preinstallation Environment",
    [string]$OutputISO = "SNISID-WinPE.iso",
    [string]$ToolsSource = "..\build\tools",
    [string]$Arch = "amd64"
)

$ErrorActionPreference = "Stop"
$PEBase = "$WinPEKit\$Arch"
$MountDir = "$env:TEMP\snisid_winpe_mount"
$ISODir = "$env:TEMP\snisid_iso"

Write-Host ">>> Cleaning previous build artifacts" -ForegroundColor Cyan
@($MountDir, $ISODir) | ForEach-Object {
    if (Test-Path $_) { Remove-Item $_ -Recurse -Force }
}

Write-Host ">>> Copying WinPE base files" -ForegroundColor Cyan
New-Item -ItemType Directory -Path $ISODir -Force
Copy-Item "$PEBase\*" $ISODir -Recurse -Force

Write-Host ">>> Copying SNISID offline tools" -ForegroundColor Cyan
$toolDest = "$ISODir\SNISID"
New-Item -ItemType Directory -Path $toolDest -Force
if (Test-Path $ToolsSource) {
    Copy-Item "$ToolsSource\*" $toolDest -Recurse -Force
} else {
    Write-Warning "Tools source not found at $ToolsSource — copying fallback scripts"
    @("snisid-offline.cmd", "snisid-viewer.exe") | ForEach-Object {
        $dummy = "$toolDest\$_"
        Set-Content -Path $dummy -Value "REM SNISID offline placeholder"
    }
}

Write-Host ">>> Placing startup script" -ForegroundColor Cyan
Copy-Item ".\startnet.cmd" "$ISODir\Windows\System32\startnet.cmd" -Force

Write-Host ">>> Creating ISO image" -ForegroundColor Cyan
$oscdimg = "$WinPEKit\..\..\Windows Kits\10\Assessment and Deployment Kit\Deployment Tools\$Arch\Oscdimg\oscdimg.exe"
if (-not (Test-Path $oscdimg)) {
    throw "oscdimg.exe not found at $oscdimg"
}

$cmd = "& '$oscdimg' -bootdata:2#p0,e,b$ISODir\boot\etfsboot.com#pEF,e,b$ISODir\efi\microsoft\boot\efisys.bin -u2 -udfver102 -lSNISID_WINPE '$ISODir' '$OutputISO'"
Invoke-Expression $cmd

Write-Host ">>> ISO created: $OutputISO" -ForegroundColor Green
