param(
    [string]$BuildPath = "..\..\build",
    [string]$SiteName = "SNISID",
    [string]$AppPoolName = "SNISID",
    [string]$Bindings = "https://+:443",
    [string]$CertThumbprint = "",
    [switch]$SkipIIS
)

$ErrorActionPreference = "Stop"

function Write-Step { Write-Host ">>> $($args[0])" -ForegroundColor Cyan }

# --- IIS Setup ---
if (-not $SkipIIS) {
    Write-Step "Installing IIS features"
    $features = @(
        "Web-WebServer",
        "Web-Common-Http",
        "Web-Default-Doc",
        "Web-Dir-Browsing",
        "Web-Http-Errors",
        "Web-Static-Content",
        "Web-Http-Logging",
        "Web-Http-Tracing",
        "Web-Request-Monitor",
        "Web-Stat-Compression",
        "Web-Dyn-Compression",
        "Web-Asp-Net45",
        "Web-ISAPI-Ext",
        "Web-ISAPI-Filter",
        "Web-Mgmt-Console",
        "Web-Scripting-Tools"
    )
    $features | ForEach-Object { Install-WindowsFeature -Name $_ }
}

# --- App Pool ---
Write-Step "Creating application pool"
if (-not (Get-IISAppPool -Name $AppPoolName -ErrorAction SilentlyContinue)) {
    New-IISAppPool -Name $AppPoolName -ManagedRuntimeVersion "v4.0" -Force
}
Set-ItemProperty -Path "IIS:\AppPools\$AppPoolName" -Name enable32BitAppOnWin64 -Value $false
Set-ItemProperty -Path "IIS:\AppPools\$AppPoolName" -Name managedPipelineMode -Value "Integrated"

# --- Web App Deploy ---
Write-Step "Deploying backend API web app"
$webPath = "C:\inetpub\wwwroot\$SiteName"
if (Test-Path $webPath) { Remove-Item "$webPath\*" -Recurse -Force }
else { New-Item -ItemType Directory -Path $webPath -Force }

Copy-Item "$BuildPath\api\*" $webPath -Recurse -Force

# Create IIS site
if (-not (Get-IISSite -Name $SiteName -ErrorAction SilentlyContinue)) {
    New-IISSite -Name $SiteName -PhysicalPath $webPath -BindingInformation $Bindings -Force
} else {
    Set-ItemProperty -Path "IIS:\Sites\$SiteName" -Name physicalPath -Value $webPath
}
Set-ItemProperty -Path "IIS:\Sites\$SiteName" -Name applicationPool -Value $AppPoolName

# --- HTTPS Binding ---
if ($CertThumbprint) {
    Write-Step "Configuring HTTPS binding"
    $existing = Get-IISSiteBinding -Name $SiteName | Where-Object { $_.Protocol -eq "https" }
    if ($existing) { Remove-IISSiteBinding -Name $SiteName -BindingInformation $existing.BindingInformation }
    New-IISSiteBinding -Name $SiteName -Protocol https -BindingInformation "https://*:443" -CertificateThumbprint $CertThumbprint -CertStoreLocation "Cert:\LocalMachine\My"
}

# --- Windows Services for Go backends ---
Write-Step "Setting up Windows Services"
$services = @(
    @{ Name = "snisid-core";  Exe = "snisid-svc.exe" },
    @{ Name = "snisid-sync";  Exe = "snisid-sync.exe" }
)
foreach ($svc in $services) {
    $binPath = Join-Path "$BuildPath\bin" $svc.Exe
    if (Get-Service $svc.Name -ErrorAction SilentlyContinue) {
        Stop-Service $svc.Name -Force; sc.exe delete $svc.Name
    }
    New-Service -Name $svc.Name -BinaryPathName $binPath -DisplayName "SNISID $($svc.Name)" -StartupType Automatic
    Start-Service $svc.Name
}

# --- Firewall ---
Write-Step "Configuring firewall rules"
New-NetFirewallRule -DisplayName "SNISID API (443)" -Direction Inbound -Protocol TCP -LocalPort 443 -Action Allow

Write-Step "Deployment complete" -ForegroundColor Green
