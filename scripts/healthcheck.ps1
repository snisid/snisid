#Requires -Version 7.0
<#
.SYNOPSIS
    SNISID E2E Healthcheck — checks /healthz on all 29 HT/NS services.
.DESCRIPTION
    Checks /healthz on all 29 services with timeout and retry logic.
    Exits with error code 1 if any service is unhealthy.
.PARAMETER PortOffset
    Port offset to apply (e.g., 90000 for e2e override). Default is 0 (direct ports).
.PARAMETER RetryCount
    Number of retries per service. Default 3.
.PARAMETER RetryDelay
    Seconds between retries. Default 5.
.PARAMETER TimeoutSeconds
    HTTP request timeout per attempt. Default 10.
#>

param(
    [int]$PortOffset = 0,
    [int]$RetryCount = 3,
    [int]$RetryDelay = 5,
    [int]$TimeoutSeconds = 10
)

$HT_SERVICES = @(
    @{ Name = "id-core";              Port = 8201 }
    @{ Name = "civil-ht";             Port = 8202 }
    @{ Name = "bio-ht";               Port = 8203 }
    @{ Name = "card-ht";              Port = 8204 }
    @{ Name = "pki-ht";               Port = 8206 }
    @{ Name = "iam-ht";               Port = 8207 }
    @{ Name = "interop-ht";           Port = 8208 }
    @{ Name = "infra-ht";             Port = 8209 }
    @{ Name = "cyber-ht";             Port = 8210 }
    @{ Name = "offline-ht";           Port = 8211 }
    @{ Name = "field-ht";             Port = 8212 }
    @{ Name = "data-ht";              Port = 8213 }
    @{ Name = "api-ht";               Port = 8214 }
    @{ Name = "foves-ht";             Port = 8215 }
    @{ Name = "lapi-ht";              Port = 8216 }
    @{ Name = "fpr-ht";               Port = 8205 }
    @{ Name = "sigint-ht";            Port = 8301 }
    @{ Name = "humint-ht";            Port = 8302 }
    @{ Name = "air-defense-ht";       Port = 8303 }
    @{ Name = "mil-c2-ht";            Port = 8304 }
    @{ Name = "bio-surveillance-ht";  Port = 8305 }
    @{ Name = "executive-protection-ht"; Port = 8306 }
    @{ Name = "transport-security-ht";   Port = 8307 }
    @{ Name = "radiation-safety-svc";     Port = 8308 }
    @{ Name = "all-source-fusion-ht";     Port = 8309 }
    @{ Name = "counterintel-ht";          Port = 8310 }
    @{ Name = "critical-infra-protection-ht"; Port = 8311 }
    @{ Name = "fisa-court-svc";       Port = 8312 }
    @{ Name = "classification-mgmt-ht";   Port = 8313 }
)

$unhealthy = @()

foreach ($svc in $HT_SERVICES) {
    $url = "http://localhost:$($svc.Port + $PortOffset)/healthz"
    $ok = $false

    for ($attempt = 1; $attempt -le $RetryCount; $attempt++) {
        try {
            $result = Invoke-WebRequest -Uri $url -Method GET -TimeoutSec $TimeoutSeconds -SkipCertificateCheck -ErrorAction Stop
            if ($result.StatusCode -eq 200) {
                Write-Host "[PASS] $($svc.Name) -> $($result.StatusCode)" -ForegroundColor Green
                $ok = $true
                break
            }
        } catch {
            if ($attempt -lt $RetryCount) {
                Write-Warning "[RETRY $attempt/$RetryCount] $($svc.Name) -> $_ (waiting ${RetryDelay}s)"
                Start-Sleep -Seconds $RetryDelay
            }
        }
    }

    if (-not $ok) {
        Write-Host "[FAIL] $($svc.Name) at $url" -ForegroundColor Red
        $unhealthy += $svc.Name
    }
}

Write-Host "`n--- Summary ---"
Write-Host "Total: $($HT_SERVICES.Count) | Healthy: $($HT_SERVICES.Count - $unhealthy.Count) | Unhealthy: $($unhealthy.Count)"

if ($unhealthy.Count -gt 0) {
    Write-Host "Unhealthy services: $($unhealthy -join ', ')" -ForegroundColor Red
    exit 1
}

Write-Host "All services healthy." -ForegroundColor Green
exit 0
