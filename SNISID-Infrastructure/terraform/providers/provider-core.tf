# SNISID — Terraform Core Provider Configurations
# Classification: RESTREINT DEFENSE
# Scope: Core DC — providers Proxmox, Kubernetes, Helm, Vault (init)

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
    vault = {
      source  = "hashicorp/vault"
      version = ">= 3.20.0"
    }
    dns = {
      source  = "hashicorp/dns"
      version = ">= 3.3.2"
    }
  }

  backend "s3" {
    bucket                      = "snisid-terraform-state"
    key                         = "core/terraform.tfstate"
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

# Provider Proxmox Core DC
provider "proxmox" {
  pm_api_url          = "https://pve-api-core.snisid.gouv.local:8006/api2/json"
  pm_user             = "snisid-terraform@pve"
  pm_tls_insecure     = false  # CA nationale validée
  pm_tls_cert_path    = "/etc/ssl/certs/snisid-ca.crt"
  pm_otp              = ""     # Smartcard/OTP via bastion
  pm_debug            = false
  pm_parallel         = 4
  pm_timeout          = 600
}

# Provider Kubernetes (Core cluster post-bootstrap)
provider "kubernetes" {
  host                   = "https://api.core.snisid.gouv.local:6443"
  cluster_ca_certificate = file("../../.secrets/core-k8s-ca.crt")
  client_certificate     = file("../../.secrets/core-k8s-client.crt")
  client_key             = file("../../.secrets/core-k8s-client.key")
}

# Provider Helm (Core cluster)
provider "helm" {
  kubernetes {
    host                   = "https://api.core.snisid.gouv.local:6443"
    cluster_ca_certificate = file("../../.secrets/core-k8s-ca.crt")
    client_certificate     = file("../../.secrets/core-k8s-client.crt")
    client_key             = file("../../.secrets/core-k8s-client.key")
  }
}

# Provider Vault (pour init secrets provisioning uniquement — pas de state dans Vault)
provider "vault" {
  address = "https://vault.snisid-identity.svc.cluster.local:8200"
  ca_cert_file = "../../.secrets/snisid-ca.crt"
  auth_login {
    path = "auth/cert/login"
    parameters = {
      name = "terraform-core"
    }
  }
}

# Provider DNS interne souverain
provider "dns" {
  update {
    server        = "coredns-core.snisid.gouv.local"
    key_name      = "snisid-tsig-key."
    key_algorithm = "hmac-sha256"
    key_secret    = var.dns_tsig_secret_core
  }
}
