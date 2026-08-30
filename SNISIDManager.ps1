# SNISID Desktop Management Utility
# Usage: .\SNISIDManager.ps1 [start|stop|restart|status]

param (
    [Parameter(Mandatory=$true)]
    [ValidateSet("start", "stop", "restart", "status")]
    $Action
)

function Get-PlatformStatus {
    Write-Host "📊 SNISID Platform Status:" -ForegroundColor Cyan
    k3d cluster list snisid
    kubectl get pods -n snisid
}

switch ($Action) {
    "start" {
        Write-Host "🚀 Starting SNISID Cluster..."
        k3d cluster start snisid
    }
    "stop" {
        Write-Host "🛑 Stopping SNISID Cluster..."
        k3d cluster stop snisid
    }
    "restart" {
        Write-Host "🔄 Restarting SNISID Cluster..."
        k3d cluster stop snisid
        k3d cluster start snisid
    }
    "status" {
        Get-PlatformStatus
    }
}
