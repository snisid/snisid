# SNISID — Terraform Module : HashiCorp Vault HA Raft + HSM Auto-Unseal
# Classification: TOP-SECRET
# Role: Déploiement Vault national via Helm sur cluster Identity (Tier-0)
# HSM: Thales Luna 7 (PKCS#11) — jamais unseal manuel en production

terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.11.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23.0"
    }
    vault = {
      source  = "hashicorp/vault"
      version = ">= 3.20.0"
    }
  }
}

locals {
  vault_fullname = "${var.deployment_name}-vault"
  vault_labels = {
    "snisid.gov/tier"    = "tier-0"
    "snisid.gov/region"  = var.region
    "snisid.gov/service" = "vault"
    "app.kubernetes.io/name"     = "vault"
    "app.kubernetes.io/instance" = var.deployment_name
    "app.kubernetes.io/component" = "server"
  }
}

# ───────────────────────────────────────────────
# Namespace Tier-0
# ───────────────────────────────────────────────
resource "kubernetes_namespace" "vault_ns" {
  metadata {
    name = var.namespace
    labels = merge(local.vault_labels, {
      "pod-security.kubernetes.io/enforce"         = "privileged"
      "pod-security.kubernetes.io/enforce-version"  = "latest"
    })
    annotations = {
      "snisid.gov/classification" = "TOP-SECRET"
    }
  }
}

# ───────────────────────────────────────────────
# ConfigMap Vault (HSM + Raft configuration)
# ───────────────────────────────────────────────
resource "kubernetes_config_map_v1" "vault_config" {
  metadata {
    name      = "${local.vault_fullname}-config"
    namespace = kubernetes_namespace.vault_ns.metadata[0].name
    labels    = local.vault_labels
  }
  data = {
    "vault.hcl" = <<-EOT
      ui = true

      listener "tcp" {
        address         = "[::]:8200"
        cluster_address = "[::]:8201"
        tls_disable     = false
        tls_cert_file   = "/vault/tls/server.crt"
        tls_key_file    = "/vault/tls/server.key"
        tls_min_version = "tls13"

        telemetry {
          prometheus_retention_time = "30s"
          disable_hostname = true
        }
      }

      storage "raft" {
        path    = "/vault/data"
        node_id = "NODE_ID"

        retry_leader_election = true

        autopilot {
          cleanup_dead_servers      = true
          last_contact_threshold    = "10s"
          max_trailing_logs         = 250
          min_quorum                = 3
          server_stabilization_time = "10s"
        }
      }

      seal "pkcs11" {
        lib            = "${var.hsm_lib_path}"
        slot           = "${var.hsm_slot}"
        key_label      = "${var.hsm_key_label}"
        hmac_key_label = "${var.hsm_hmac_key_label}"
        generate_key   = "true"
      }

      service_registration "kubernetes" {}

      telemetry {
        prometheus_retention_time = "30s"
        disable_hostname = true
      }

      audit_device "file" {
        path        = "/vault/audit/vault_audit.log"
        log_raw     = false
        hmac_accessor = true
        mode        = "0644"
        format      = "json"
      }
    EOT
  }
}

# ───────────────────────────────────────────────
# Secret TLS (cert-manager ou manuel)
# ───────────────────────────────────────────────
resource "kubernetes_secret_v1" "vault_tls" {
  metadata {
    name      = "${local.vault_fullname}-tls"
    namespace = kubernetes_namespace.vault_ns.metadata[0].name
    labels    = local.vault_labels
  }
  type = "kubernetes.io/tls"
  data = {
    "tls.crt" = var.vault_tls_cert_pem
    "tls.key" = var.vault_tls_key_pem
    "ca.crt"  = var.vault_ca_cert_pem
  }
}

# ───────────────────────────────────────────────
# ServiceAccount (no automount unless injector needs it)
# ───────────────────────────────────────────────
resource "kubernetes_service_account_v1" "vault" {
  metadata {
    name      = local.vault_fullname
    namespace = kubernetes_namespace.vault_ns.metadata[0].name
    labels    = local.vault_labels
  }
  automount_service_account_token = true
}

# ───────────────────────────────────────────────
# RBAC for Vault service registration in K8s
# ───────────────────────────────────────────────
resource "kubernetes_cluster_role_v1" "vault" {
  metadata {
    name = "${local.vault_fullname}-role"
    labels = local.vault_labels
  }
  rule {
    api_groups = [""]
    resources  = ["pods", "endpoints", "services", "nodes"]
    verbs      = ["get", "list", "watch", "update", "patch"]
  }
  rule {
    api_groups = ["apps"]
    resources  = ["statefulsets"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_cluster_role_binding_v1" "vault" {
  metadata {
    name = "${local.vault_fullname}-binding"
    labels = local.vault_labels
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role_v1.vault.metadata[0].name
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.vault.metadata[0].name
    namespace = kubernetes_namespace.vault_ns.metadata[0].name
  }
}

# ───────────────────────────────────────────────
# StatefulSet Vault HA Raft (5 replicas)
# ───────────────────────────────────────────────
resource "kubernetes_stateful_set_v1" "vault" {
  metadata {
    name      = local.vault_fullname
    namespace = kubernetes_namespace.vault_ns.metadata[0].name
    labels    = local.vault_labels
  }
  spec {
    replicas     = var.replicas
    service_name = "${local.vault_fullname}-internal"
    selector {
      match_labels = {
        app.kubernetes.io/name     = "vault"
        app.kubernetes.io/instance = var.deployment_name
      }
    }
    template {
      metadata {
        labels = merge(local.vault_labels, {
          app.kubernetes.io/name     = "vault"
          app.kubernetes.io/instance = var.deployment_name
        })
        annotations = {
          "prometheus.io/scrape" = "true"
          "prometheus.io/port"   = "8200"
          "snisid.gov/classification" = "TOP-SECRET"
        }
      }
      spec {
        service_account_name = kubernetes_service_account_v1.vault.metadata[0].name
        affinity {
          pod_anti_affinity {
            required_during_scheduling_ignored_during_execution {
              label_selector {
                match_labels = {
                  app.kubernetes.io/name     = "vault"
                  app.kubernetes.io/instance = var.deployment_name
                }
              }
              topology_key = "topology.kubernetes.io/zone"
            }
          }
        }
        toleration {
          key      = "snisid.gov/tier"
          operator = "Equal"
          value    = "tier-0"
          effect   = "NoSchedule"
        }
        node_selector = {
          "snisid.gov/tier" = "tier-0"
        }
        security_context {
          run_as_non_root     = true
          run_as_user         = 100
          run_as_group        = 1000
          fs_group            = 1000
          read_only_root_filesystem = true
          seccomp_profile {
            type = "RuntimeDefault"
          }
        }
        container {
          name  = "vault"
          image = "${var.registry}/hashicorp/vault:${var.vault_version}"
          command = ["/bin/sh", "-c"]
          args = ["sed -i 's/NODE_ID/${HOSTNAME}/' /vault/config/vault.hcl && vault server -config=/vault/config/vault.hcl"]
          port {
            container_port = 8200
            name           = "http"
            protocol       = "TCP"
          }
          port {
            container_port = 8201
            name           = "internal"
            protocol       = "TCP"
          }
          resources {
            limits = {
              memory = "8Gi"
              cpu    = "4000m"
            }
            requests = {
              memory = "4Gi"
              cpu    = "2000m"
            }
          }
          volume_mount {
            name       = "config"
            mount_path = "/vault/config"
            read_only  = true
          }
          volume_mount {
            name       = "tls"
            mount_path = "/vault/tls"
            read_only  = true
          }
          volume_mount {
            name       = "data"
            mount_path = "/vault/data"
          }
          volume_mount {
            name       = "audit"
            mount_path = "/vault/audit"
          }
          liveness_probe {
            http_get {
              path   = "/v1/sys/health?standbyok=true&sealedcode=200"
              port   = 8200
              scheme = "HTTPS"
            }
            initial_delay_seconds = 30
            period_seconds          = 15
            failure_threshold       = 3
          }
          readiness_probe {
            http_get {
              path   = "/v1/sys/health?standbyok=true&sealedcode=200&uninitcode=200"
              port   = 8200
              scheme = "HTTPS"
            }
            initial_delay_seconds = 10
            period_seconds          = 10
            failure_threshold       = 3
          }
          security_context {
            read_only_root_filesystem = true
            allow_privilege_escalation = false
            capabilities {
              drop = ["ALL"]
            }
          }
        }
        volume {
          name = "config"
          config_map {
            name         = kubernetes_config_map_v1.vault_config.metadata[0].name
            default_mode = "0644"
          }
        }
        volume {
          name = "tls"
          secret {
            secret_name = kubernetes_secret_v1.vault_tls.metadata[0].name
          }
        }
      }
    }
    volume_claim_template {
      metadata {
        name = "data"
        labels = local.vault_labels
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        storage_class_name = var.storage_class
        resources {
          requests = {
            storage = "100Gi"
          }
        }
      }
    }
    volume_claim_template {
      metadata {
        name = "audit"
        labels = local.vault_labels
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        storage_class_name = var.storage_class
        resources {
          requests = {
            storage = "50Gi"
          }
        }
      }
    }
  }
}

# ───────────────────────────────────────────────
# Service Headless (Raft cluster communication)
# ───────────────────────────────────────────────
resource "kubernetes_service_v1" "vault_internal" {
  metadata {
    name      = "${local.vault_fullname}-internal"
    namespace = kubernetes_namespace.vault_ns.metadata[0].name
    labels    = local.vault_labels
  }
  spec {
    selector = {
      app.kubernetes.io/name     = "vault"
      app.kubernetes.io/instance = var.deployment_name
    }
    port {
      name        = "http"
      port        = 8200
      target_port = 8200
      protocol    = "TCP"
    }
    port {
      name        = "internal"
      port        = 8201
      target_port = 8201
      protocol    = "TCP"
    }
    cluster_ip = "None"
    publish_not_ready_addresses = true
  }
}

# ───────────────────────────────────────────────
# Service LoadBalancer/ClusterIP pour accès Vault
# ───────────────────────────────────────────────
resource "kubernetes_service_v1" "vault" {
  metadata {
    name      = local.vault_fullname
    namespace = kubernetes_namespace.vault_ns.metadata[0].name
    labels    = local.vault_labels
  }
  spec {
    selector = {
      app.kubernetes.io/name     = "vault"
      app.kubernetes.io/instance = var.deployment_name
      app.kubernetes.io/component = "server"
    }
    port {
      name        = "http"
      port        = 8200
      target_port = 8200
      protocol    = "TCP"
    }
    type = "ClusterIP"
  }
}
