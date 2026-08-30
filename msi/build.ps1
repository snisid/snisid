param(
    [string]$WixPath = "${env:ProgramFiles(x86)}\WiX Toolset v3.11\bin",
    [string]$Source = "snisid.wxs",
    [string]$OutputDir = ".\out"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $WixPath)) {
    throw "WiX Toolset not found at $WixPath. Install from https://wixtoolset.org"
}

$candle = "$WixPath\candle.exe"
$light  = "$WixPath\light.exe"

New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null

Write-Host ">>> Compiling WiX source with candle.exe" -ForegroundColor Cyan
& $candle -arch x64 -out "$OutputDir\snisid.wixobj" $Source
if (-not $?) { throw "Candle compilation failed" }

Write-Host ">>> Linking MSI with light.exe" -ForegroundColor Cyan
& $light -out "$OutputDir\snisid-1.0.0.msi" "$OutputDir\snisid.wixobj" -ext WixUIExtension -cultures:en-us
if (-not $?) { throw "Light linking failed" }

Write-Host ">>> MSI created: $OutputDir\snisid-1.0.0.msi" -ForegroundColor Green
