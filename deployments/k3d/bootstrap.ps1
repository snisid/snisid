$ErrorActionPreference = "Stop"

$cluster = "snisid"
k3d cluster create $cluster --agents 2 --servers 1 --port "8080:80@loadbalancer"
kubectl create namespace snisid
helm upgrade --install snisid ./deploy/helm/snisid -n snisid
Write-Host "SNISID deployed on k3d cluster '$cluster'."
