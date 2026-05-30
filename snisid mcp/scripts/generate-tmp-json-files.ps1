$ErrorActionPreference = "Stop"

Set-Location "C:\Users\sopil\Desktop\snisid system\snisid mcp"

New-Item -ItemType Directory -Force ".\tmp" | Out-Null

Write-Host "Generation des tokens JWT/MFA frais..." -ForegroundColor Cyan

$auth = npx tsx scripts/create-test-token.ts | ConvertFrom-Json

$authObj = [ordered]@{
  accessToken   = [string]$auth.accessToken
  mfaToken      = [string]$auth.mfaToken
  deviceId      = "device-test-001"
  purpose       = "verification legale de test"
  correlationId = "corr-test-0001"
  sessionId     = "sess-test-0001"
}

function Save-JsonFile {
  param(
    [string]$FileName,
    [object]$Object
  )

  $path = Join-Path ".\tmp" $FileName
  $Object | ConvertTo-Json -Depth 50 | Set-Content -Encoding UTF8 $path
}

$visaFormats = [ordered]@{
  usa = [ordered]@{
    country = "US"
    exampleVisaNumber = "12345678"
    exampleCaseNumber = "USA-2026-12345678"
    note = "Numero de test representatif pour visa USA. Ne pas utiliser comme donnee reelle."
  }
  canada = [ordered]@{
    country = "CA"
    exampleVisaNumber = "V123456789"
    exampleUci = "1234-5678"
    note = "Numero de test representatif pour visa Canada. Ne pas utiliser comme donnee reelle."
  }
  france = [ordered]@{
    country = "FR"
    exampleVisaNumber = "FRA123456789"
    exampleSchengenSticker = "FR123456789"
    note = "Numero de test representatif pour visa France/Schengen."
  }
  schengen = [ordered]@{
    country = "SCHENGEN"
    exampleVisaNumber = "SCH123456789"
    note = "Format de test generique pour espace Schengen."
  }
  unitedKingdom = [ordered]@{
    country = "GB"
    exampleVisaNumber = "UK123456789"
    note = "Numero de test representatif pour visa Royaume-Uni."
  }
  dominicanRepublic = [ordered]@{
    country = "DO"
    exampleVisaNumber = "DO123456789"
    note = "Numero de test representatif pour visa Republique Dominicaine."
  }
  mexico = [ordered]@{
    country = "MX"
    exampleVisaNumber = "MX123456789"
    note = "Numero de test representatif pour visa Mexique."
  }
  brazil = [ordered]@{
    country = "BR"
    exampleVisaNumber = "BR123456789"
    note = "Numero de test representatif pour visa Bresil."
  }
  chile = [ordered]@{
    country = "CL"
    exampleVisaNumber = "CL123456789"
    note = "Numero de test representatif pour visa Chili."
  }
  panama = [ordered]@{
    country = "PA"
    exampleVisaNumber = "PA123456789"
    note = "Numero de test representatif pour visa Panama."
  }
}

$passportReference = [ordered]@{
  rule = "1 ou 2 lettres suivies de 9 chiffres"
  examples = @(
    "A123456789",
    "AB123456789"
  )
}

$haitiSubjectRefRecommendations = [ordered]@{
  note = "References de test uniquement. Ne pas utiliser de donnees personnelles reelles."
  recommendedSubjectRefs = @(
    "HT-NID-2026-000001",
    "HT-PASSPORT-A123456789",
    "HT-NIF-123-456-789-0",
    "HT-PHONE-509-0000-0000",
    "HT-BUSINESS-123-456-789-0"
  )
}

Save-JsonFile "auth.current.json" $authObj

Save-JsonFile "visa.formats.reference.json" $visaFormats
Save-JsonFile "passport.formats.reference.json" $passportReference
Save-JsonFile "haiti.subjectRef.recommendations.json" $haitiSubjectRefRecommendations

Save-JsonFile "identity.birthCertificate.json" ([ordered]@{
  nationalId = "HT12345"
  certificateNumber = "A1234567890"
  qrCode = "A1234567890"
  numeroDemande = "D2026-2035/1234567890"
  auth = $authObj
})

Save-JsonFile "identity.passportLookup.json" ([ordered]@{
  passportNumber = "A123456789"
  passportFormatRule = $passportReference
  auth = $authObj
})

Save-JsonFile "justice.criminalRecord.json" ([ordered]@{
  nationalId = "HT12345"
  caseScope = "SUMMARY"
  availableCaseScopes = @("SUMMARY", "FULL")
  auth = $authObj
})

Save-JsonFile "justice.criminalRecord.SUMMARY.json" ([ordered]@{
  nationalId = "HT12345"
  caseScope = "SUMMARY"
  auth = $authObj
})

Save-JsonFile "justice.criminalRecord.FULL.json" ([ordered]@{
  nationalId = "HT12345"
  caseScope = "FULL"
  auth = $authObj
})

Save-JsonFile "immigration.borderAlerts.json" ([ordered]@{
  nationalId = "HT12345"
  passportNumber = "A123456789"
  passportFormatRule = $passportReference
  auth = $authObj
})

Save-JsonFile "immigration.visaLookup.json" ([ordered]@{
  passportNumber = "A123456789"
  country = "US"
  visaFormats = $visaFormats
  auth = $authObj
})

Save-JsonFile "immigration.visaLookup.USA.json" ([ordered]@{
  passportNumber = "A123456789"
  country = "US"
  visaNumber = "12345678"
  visaFormats = $visaFormats.usa
  auth = $authObj
})

Save-JsonFile "immigration.visaLookup.CANADA.json" ([ordered]@{
  passportNumber = "A123456789"
  country = "CA"
  visaNumber = "V123456789"
  visaFormats = $visaFormats.canada
  auth = $authObj
})

Save-JsonFile "immigration.visaLookup.FRANCE.json" ([ordered]@{
  passportNumber = "A123456789"
  country = "FR"
  visaNumber = "FRA123456789"
  visaFormats = $visaFormats.france
  auth = $authObj
})

Save-JsonFile "immigration.entryExit.json" ([ordered]@{
  passportNumber = "A123456789"
  visaNumbers = $visaFormats
  fromDate = "2026-01-01"
  toDate = "2026-12-31"
  auth = $authObj
})

Save-JsonFile "immigration.watchlistScan.json" ([ordered]@{
  nationalId = "HT12345"
  passportNumber = "A123456789"
  name = "TEST PERSON"
  visaNumbers = $visaFormats
  auth = $authObj
})

Save-JsonFile "education.diplomaVerification.json" ([ordered]@{
  diplomaNumber = "1234567890"
  nationalId = "HT12345"
  auth = $authObj
})

Save-JsonFile "tax.verifyNIF.json" ([ordered]@{
  nif = "123-456-789-0"
  auth = $authObj
})

Save-JsonFile "tax.businessRegistry.json" ([ordered]@{
  registrationNumber = "123-456-789-0"
  auth = $authObj
})

Save-JsonFile "intelligence.behaviorAnalysis.json" ([ordered]@{
  subjectRef = "HT-NID-2026-000001"
  scope = "LOW"
  recommendedSubjectRefs = $haitiSubjectRefRecommendations.recommendedSubjectRefs
  auth = $authObj
})

Save-JsonFile "intelligence.behaviorAnalysis.PASSPORT.json" ([ordered]@{
  subjectRef = "HT-PASSPORT-A123456789"
  scope = "LOW"
  auth = $authObj
})

Save-JsonFile "intelligence.behaviorAnalysis.NIF.json" ([ordered]@{
  subjectRef = "HT-NIF-123-456-789-0"
  scope = "LOW"
  auth = $authObj
})

Save-JsonFile "intelligence.behaviorAnalysis.PHONE.json" ([ordered]@{
  subjectRef = "HT-PHONE-509-0000-0000"
  scope = "LOW"
  auth = $authObj
})

@"
DOSSIER TMP SNISID MCP

Ce dossier contient des fichiers JSON de test pour MCP Inspector.

IMPORTANT :
- Les tokens expirent.
- Si tu recois jwt expired ou MFA_REQUIRED, relance :
  powershell -ExecutionPolicy Bypass -File .\scripts\generate-tmp-json-files.ps1

FICHIERS PRINCIPAUX :

identity.birthCertificate.json
identity.passportLookup.json

justice.criminalRecord.json
justice.criminalRecord.SUMMARY.json
justice.criminalRecord.FULL.json

immigration.borderAlerts.json
immigration.visaLookup.json
immigration.visaLookup.USA.json
immigration.visaLookup.CANADA.json
immigration.visaLookup.FRANCE.json
immigration.entryExit.json
immigration.watchlistScan.json

education.diplomaVerification.json

tax.verifyNIF.json
tax.businessRegistry.json

intelligence.behaviorAnalysis.json
intelligence.behaviorAnalysis.PASSPORT.json
intelligence.behaviorAnalysis.NIF.json
intelligence.behaviorAnalysis.PHONE.json

REFERENCES :
visa.formats.reference.json
passport.formats.reference.json
haiti.subjectRef.recommendations.json

EXEMPLE POUR COPIER UN FICHIER DANS LE PRESSE-PAPIERS :

Get-Content .\tmp\identity.birthCertificate.json -Raw | Set-Clipboard

Ensuite :
MCP Inspector -> Tools -> choisir le tool -> Input / Arguments -> CTRL + V -> Run Tool
"@ | Set-Content -Encoding UTF8 ".\tmp\README.txt"

Get-Content ".\tmp\identity.birthCertificate.json" -Raw | Set-Clipboard

Write-Host "Dossier tmp mis a jour avec succes." -ForegroundColor Green
Write-Host "Fichier copie dans le presse-papiers : identity.birthCertificate.json" -ForegroundColor Green
Write-Host "Repertoire : C:\Users\sopil\Desktop\snisid system\snisid mcp\tmp" -ForegroundColor Cyan

Get-ChildItem ".\tmp" | Select-Object Name, Length
