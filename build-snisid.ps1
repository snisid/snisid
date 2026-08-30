$ErrorActionPreference = "Continue"
$env:GOOS = "windows"; $env:GOARCH = "amd64"; $env:CGO_ENABLED = "0"

$services = @(
    @{ name = "api-platform";        path = ".\services\api-platform\cmd" },
    @{ name = "identity-api";        path = ".\services\identity-api\cmd" },
    @{ name = "federation-gateway";  path = ".\services\federation-gateway\cmd" },
    @{ name = "enrollment-service";  path = ".\services\enrollment-service\cmd" },
    @{ name = "gateway";             path = ".\gateway\cmd" }
)

$results = @()
foreach ($s in $services) {
    $exe = "$($s.name).exe"
    if (Test-Path $exe) { Remove-Item $exe -Force }
    Write-Host "`n--- Building $($s.name) ---" -ForegroundColor Cyan
    & go build -o $exe $s.path 2>&1 | Out-Host
    $ok = Test-Path $exe
    $results += [pscustomobject]@{ Service = $s.name; Status = if ($ok) { "OK" } else { "FAIL" } }
}
$results | Format-Table -AutoSize
