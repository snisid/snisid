# Script Windows PowerShell pour générer un .env local de test SNISID.
# Usage :
#   npm run setup:local

$ErrorActionPreference = 'Stop'

function New-HexSecret([int]$Bytes) {
  node -e "console.log(require('crypto').randomBytes($Bytes).toString('hex'))"
}

function New-Base64Secret([int]$Bytes) {
  node -e "console.log(require('crypto').randomBytes($Bytes).toString('base64'))"
}

$jwtSecret = New-HexSecret 64
$encryptionKey = New-Base64Secret 32
$apiPepper = New-HexSecret 32
$govApiKey = "gov_$(New-HexSecret 32)"

@"
NODE_ENV=development
PORT=3001
MCP_TRANSPORT=http
SNISID_DEMO_MODE=true

JWT_ISSUER=snisid-iam
JWT_AUDIENCE=snisid-mcp
JWT_SECRET=$jwtSecret
JWT_EXPIRES_IN=15m
ENCRYPTION_KEY_B64=$encryptionKey
API_KEY_PEPPER=$apiPepper

API_GATEWAY_BASE_URL=http://localhost:4000
ONI_API_BASE_URL=http://localhost:4000/oni
DGI_API_BASE_URL=http://localhost:4000/dgi
PNH_API_BASE_URL=http://localhost:4000/pnh
MJSP_API_BASE_URL=http://localhost:4000/mjsp
IMMIGRATION_API_BASE_URL=http://localhost:4000/immigration
ANH_API_BASE_URL=http://localhost:4000/anh
PASSPORT_API_BASE_URL=http://localhost:4000/passport
BIOMETRIC_API_BASE_URL=http://localhost:4000/biometric
EDUCATION_API_BASE_URL=http://localhost:4000/education
INTELLIGENCE_API_BASE_URL=http://localhost:4000/intelligence
GOV_API_KEY=$govApiKey

REDIS_URL=redis://localhost:6379
QDRANT_URL=http://localhost:6333
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
LOG_LEVEL=info

OPENAI_API_KEY=
OPENAI_BASE_URL=https://api.openai.com/v1
DEEPSEEK_API_KEY=
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
ANTHROPIC_API_KEY=
ANTHROPIC_BASE_URL=https://api.anthropic.com
GOOGLE_ANTIGRAVITY_API_KEY=
GOOGLE_ANTIGRAVITY_BASE_URL=
ARENA_AI_API_KEY=
ARENA_AI_BASE_URL=
CLAUDE_API_KEY=
CURSOR_BRIDGE_URL=http://localhost:39001
VSCODE_BRIDGE_URL=http://localhost:39002
GITHUB_COPILOT_BRIDGE_URL=http://localhost:39003
"@ | Set-Content -Encoding UTF8 .env

Write-Host "✅ .env local généré avec succès." -ForegroundColor Green
Write-Host "✅ GOV_API_KEY générée pour les APIs mock." -ForegroundColor Green
Write-Host "Tu peux maintenant lancer : npm run mock:gov" -ForegroundColor Cyan
