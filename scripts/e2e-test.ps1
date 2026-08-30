#Requires -Version 7.0
<#
.SYNOPSIS
    SNISID E2E Integration Tests — verifies key inter-service flows.
.DESCRIPTION
    Runs 4 inter-service flow tests against a running SNISID environment.
    Uses Invoke-RestMethod with JSON payloads.
.PARAMETER BaseUrl
    Base URL for all services (default http://localhost).
.PARAMETER PortOffset
    Port offset for e2e override ports (default 0).
#>

param(
    [string]$BaseUrl = "http://localhost",
    [int]$PortOffset = 0
)

$passed = 0
$failed = 0
$errors = @()

function Invoke-Api {
    param([string]$Service, [int]$Port, [string]$Method = "GET", [string]$Endpoint, $Body = $null)
    $url = "$BaseUrl`:$($Port + $PortOffset)$Endpoint"
    $params = @{ Uri = $url; Method = $Method; ContentType = "application/json" }
    if ($Body) { $params["Body"] = ($Body | ConvertTo-Json -Depth 10) }
    try {
        $result = Invoke-RestMethod @params -TimeoutSec 15 -ErrorAction Stop
        return $result
    } catch {
        throw "Request failed: $Method $url -> $_"
    }
}

function Assert-Equal {
    param([string]$TestName, $Expected, $Actual)
    if ($Expected -ne $Actual) {
        throw "Assertion failed [$TestName]: expected '$Expected', got '$Actual'"
    }
    Write-Host "  [OK] $TestName" -ForegroundColor Green
}

function Run-TestFlow {
    param([string]$Name, [scriptblock]$Script)
    Write-Host "`n=== Test Flow: $Name ===" -ForegroundColor Cyan
    try {
        & $Script
        Write-Host "[PASS] $Name" -ForegroundColor Green
        $script:passed++
    } catch {
        Write-Host "[FAIL] $Name : $_" -ForegroundColor Red
        $script:failed++
        $script:errors += "[$Name] $_"
    }
}

# --- Flow 1: Identity → Civil → FPR ---
Run-TestFlow -Name "Identity → Civil → FPR" -Script {
    # 1a. Enroll citizen via id-core
    $citizen = Invoke-Api -Service "id-core" -Port 8201 -Method POST -Endpoint "/api/v1/citizens" -Body @{
        firstName = "Jean"
        lastName  = "Dupont"
        dob       = "1990-01-15"
        nationalId = "HT-1990-0001"
    }
    $citizenId = $citizen.id
    Write-Host "  Citizen enrolled: $citizenId"

    # 1b. Create birth record via civil-ht
    $birthRecord = Invoke-Api -Service "civil-ht" -Port 8202 -Method POST -Endpoint "/api/v1/birth-records" -Body @{
        citizenId   = $citizenId
        registrant  = "Jean Dupont"
        dateOfBirth = "1990-01-15"
        placeOfBirth = "Port-au-Prince"
    }
    $recordId = $birthRecord.id
    Write-Host "  Birth record created: $recordId"

    # 1c. Issue warrant via fpr-ht
    $warrant = Invoke-Api -Service "fpr-ht" -Port 8205 -Method POST -Endpoint "/api/v1/warrants" -Body @{
        subjectId        = $citizenId
        warrantType      = "Criminal_Investigation"
        issuingAuthority = "Tribunal de Premiere Instance"
        reason           = "Fraud investigation"
    }
    $warrantId = $warrant.id
    Write-Host "  Warrant issued: $warrantId"

    # 1d. Check citizen appears in FPR
    $fprRecord = Invoke-Api -Service "fpr-ht" -Port 8205 -Method GET -Endpoint "/api/v1/persons/$citizenId"
    Assert-Equal -TestName "Citizen found in FPR" -Expected $citizenId -Actual $fprRecord.id
}

# --- Flow 2: SIGINT → FISA → HUMINT ---
Run-TestFlow -Name "SIGINT → FISA → HUMINT" -Script {
    # 2a. Create FISA warrant via fisa-court-svc
    $fisaWarrant = Invoke-Api -Service "fisa-court-svc" -Port 8312 -Method POST -Endpoint "/api/v1/fisa-warrants" -Body @{
        targetName      = "Foreign Agent X"
        targetNationalId = "HT-FX-2026"
        warrantType     = "FISA_Title_I"
        courtDocket     = "FISA-2026-0042"
        expirationDate  = "2027-01-01"
    }
    $fisaId = $fisaWarrant.id
    Write-Host "  FISA warrant created: $fisaId"

    # 2b. Create interception target via sigint-ht
    $target = Invoke-Api -Service "sigint-ht" -Port 8301 -Method POST -Endpoint "/api/v1/interception-targets" -Body @{
        fisaWarrantId = $fisaId
        targetIdentifier = "comm-channel-42"
        commType         = "Satellite"
        authorizee       = "DGSI"
    }
    $targetId = $target.id
    Write-Host "  Interception target created: $targetId"

    # 2c. Record intercepted communication
    $intercept = Invoke-Api -Service "sigint-ht" -Port 8301 -Method POST -Endpoint "/api/v1/intercepted-comms" -Body @{
        targetId     = $targetId
        contentHash  = "sha256:a" + "b" * 62
        metadata     = @{ frequency = "12.4GHz"; timestamp = "2026-06-20T10:00:00Z" }
        classification = "TOP_SECRET"
    }
    $commId = $intercept.id
    Write-Host "  Intercepted communication recorded: $commId"

    # 2d. Create HUMINT report cross-referencing SIGINT
    $report = Invoke-Api -Service "humint-ht" -Port 8302 -Method POST -Endpoint "/api/v1/intelligence-reports" -Body @{
        title       = "Source Report re Foreign Agent X"
        sourceRef   = "HUMINT-2026-071"
        sigintRefs  = @($commId)
        fisaWarrantRef = $fisaId
        summary     = "Source confirms identity of Foreign Agent X"
        classification = "TOP_SECRET"
    }
    $reportId = $report.id
    Write-Host "  HUMINT report created: $reportId"

    # Verify cross-reference is stored
    $verifyReport = Invoke-Api -Service "humint-ht" -Port 8302 -Method GET -Endpoint "/api/v1/intelligence-reports/$reportId"
    Assert-Equal -TestName "Report references SIGINT comm" -Expected $commId -Actual $verifyReport.sigintRefs[0]
}

# --- Flow 3: Air Defense → Military C2 ---
Run-TestFlow -Name "Air Defense → Military C2" -Script {
    # 3a. Create radar contact via air-defense-ht
    $radarContact = Invoke-Api -Service "air-defense-ht" -Port 8303 -Method POST -Endpoint "/api/v1/radar-contacts" -Body @{
        trackId     = "TRK-2026-8912"
        bearing     = 270.5
        range       = 45.2
        altitude    = 10000
        speed       = 450
        classification = "UNKNOWN"
    }
    $contactId = $radarContact.id
    Write-Host "  Radar contact created: $contactId"

    # 3b. Open military operation via mil-c2-ht
    $operation = Invoke-Api -Service "mil-c2-ht" -Port 8304 -Method POST -Endpoint "/api/v1/operations" -Body @{
        opName   = "Operation Sentinel Shield"
        opType   = "AIR_DEFENSE"
        priority = "HIGH"
        region   = "Northern Airspace"
    }
    $opId = $operation.id
    Write-Host "  Military operation opened: $opId"

    # 3c. Link radar contact to operation
    $link = Invoke-Api -Service "mil-c2-ht" -Port 8304 -Method POST -Endpoint "/api/v1/operations/$opId/contacts" -Body @{
        contactId = $contactId
        role       = "HOSTILE_TRACK"
    }
    Write-Host "  Radar contact linked to operation"

    # 3d. Submit SITREP
    $sitrep = Invoke-Api -Service "mil-c2-ht" -Port 8304 -Method POST -Endpoint "/api/v1/sitreps" -Body @{
        operationId = $opId
        reportType  = "SITREP"
        summary     = "Unknown track intercepted at 270 bearing, 45.2nm range. Intercept initiated."
        classification = "SECRET"
    }
    $sitrepId = $sitrep.id
    Write-Host "  SITREP submitted: $sitrepId"

    # Verify operation has the contact
    $opDetails = Invoke-Api -Service "mil-c2-ht" -Port 8304 -Method GET -Endpoint "/api/v1/operations/$opId"
    Assert-Equal -TestName "Operation has radar contact" -Expected $contactId -Actual $opDetails.contacts[0].contactId
}

# --- Flow 4: Biosurveillance → Critical Infrastructure ---
Run-TestFlow -Name "Biosurveillance → Critical Infrastructure" -Script {
    # 4a. Report disease outbreak via bio-surveillance-ht
    $outbreak = Invoke-Api -Service "bio-surveillance-ht" -Port 8305 -Method POST -Endpoint "/api/v1/outbreaks" -Body @{
        disease      = "Cholera"
        region       = "Artibonite Dept"
        confirmedCases = 47
        suspectedCases  = 120
        severity     = "EMERGENCY"
        reportedBy   = "MSPP"
    }
    $outbreakId = $outbreak.id
    Write-Host "  Disease outbreak reported: $outbreakId"

    # 4b. Check health facility stock
    $stock = Invoke-Api -Service "bio-surveillance-ht" -Port 8305 -Method GET -Endpoint "/api/v1/facilities/stocks" -Body @{
        facilityId  = "HOP-ART-001"
        region      = "Artibonite"
    }
    Write-Host "  Health facility stock checked"

    # 4c. Create infrastructure incident for hospital
    $incident = Invoke-Api -Service "critical-infra-protection-ht" -Port 8311 -Method POST -Endpoint "/api/v1/infrastructure-incidents" -Body @{
        facilityId   = "HOP-ART-001"
        category     = "HEALTH"
        incidentType = "SUPPLY_SHORTAGE"
        severity     = "CRITICAL"
        description  = "Cholera outbreak - IV fluids and oral rehydration salts below emergency threshold"
        linkedOutbreakId = $outbreakId
    }
    $incidentId = $incident.id
    Write-Host "  Infrastructure incident created: $incidentId"

    # Verify outbreak linked to incident
    $verifyIncident = Invoke-Api -Service "critical-infra-protection-ht" -Port 8311 -Method GET -Endpoint "/api/v1/infrastructure-incidents/$incidentId"
    Assert-Equal -TestName "Incident linked to outbreak" -Expected $outbreakId -Actual $verifyIncident.linkedOutbreakId
}

Write-Host "`n============================================"
Write-Host "E2E Test Results: $passed passed, $failed failed" -ForegroundColor $(if ($failed -eq 0) { "Green" } else { "Red" })
Write-Host "============================================"

if ($errors.Count -gt 0) {
    Write-Host "`nErrors:" -ForegroundColor Red
    $errors | ForEach-Object { Write-Host "  - $_" }
    exit 1
}
exit 0
