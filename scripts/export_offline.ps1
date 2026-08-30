# SNISID Air-Gapped Image Exporter
# Usage: Run on an online machine to prepare for offline install

Write-Host "📦 Preparing Offline SOC Images..." -ForegroundColor Cyan

$images = @(
    "snisid/identity-api:latest",
    "snisid/fraud-engine:latest",
    "snisid/ai-face:latest",
    "snisid/web:latest",
    "postgres:16",
    "bitnami/kafka:latest",
    "neo4j:5.22"
)

foreach ($img in $images) {
    Write-Host "📥 Pulling $img..."
    docker pull $img
}

Write-Host "💾 Exporting images to offline_images.tar..."
docker save $images -o offline_images.tar

Write-Host "✅ Export complete. Copy offline_images.tar to the target machine." -ForegroundColor Green
