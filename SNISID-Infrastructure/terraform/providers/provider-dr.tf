# SNISID — Terraform DR Provider Configurations
# Classification: RESTREINT DEFENSE
# Scope: Secondary DR DC — providers Proxmox DR, Kubernetes DR, Helm DR

terraform {
  required_version = ">= 1.6.0"

  required_providers {
    proxmox = {
      source  = "Telmate/proxmox"
      version = ">= 2.9.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.11.0"
    }
  }

  backend "s3" {
    bucket                      = "snisid-terraform-state"
    key                         = "dr/terraform.tfstate"
    region                      = "national"
    endpoint                    = "https://s3.interne.snisid.gouv.local"
    encrypt                     = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    use_path_style              = true
  }
}

provider "proxmox" {
  alias               = "dr"
  pm_api_url          = "https://pve-api-dr.snisid.gouv.local:8006/api2/json"
  pm_user             = "snisid-terraform@pve"
  pm_tls_insecure     = false
  pm_tls_cert_path    = "/etc/ssl/certs/snisid-ca.crt"
  pm_parallel         = 4
  pm_timeout          = 600
}

provider "kubernetes" {
  alias                  = "dr"
  host                   = "https://api.dr.snisid.gouv.local:6443"
  cluster_ca_certificate = file("../../.secrets/dr-k8s-ca.crt")
  client_certificate     = file("../../.secrets/dr-k8s-client.crt")
  client_key             = file("../../.secrets/dr-k8s-client.key")
}

provider "helm" {
  alias = "dr"
  kubernetes {
    host                   = "https://api.dr.snisid.gouv.local:6443"
    cluster_ca_certificate = file("../../.secrets/dr-k8s-ca.crt")
    client_certificate     = file("../../.secrets/dr-k8s-client.crt")
    client_key             = file("../../.secrets/dr-k8s-client.key")
  }
}
