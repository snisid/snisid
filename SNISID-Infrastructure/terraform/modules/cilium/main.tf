# SNISID — Terraform Module : Cilium CNI eBPF (Cluster-wide)
# Classification: SECRET
# Role: Déploiement Cilium via Helm sur clusters nationaux
# Security: Kube-proxy replacement, WireGuard encryption, Hubble, Tetragon

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
  }
}

locals {
  cilium_values = {
    kubeProxyReplacement = true
    k8sServiceHost       = var.k8s_api_host
    k8sServicePort       = var.k8s_api_port
    image = {
      repository = "${var.registry}/cilium/cilium"
      tag        = var.cilium_version
      useDigest  = false
    }
    operator = {
      replicas = var.operator_replicas
      image = {
        repository = "${var.registry}/cilium/operator"
        tag        = var.cilium_version
      }
      resources = {
        limits = { memory = "512Mi", cpu = "500m" }
        requests = { memory = "256Mi", cpu = "100m" }
      }
    }
    ipam = {
      mode = "kubernetes"
    }
    hubble = {
      enabled = var.enable_hubble
      relay = {
        enabled  = true
        replicas = 2
        resources = {
          limits = { memory = "1Gi", cpu = "1000m" }
          requests = { memory = "256Mi", cpu = "100m" }
        }
      }
      ui = {
        enabled = var.enable_hubble_ui
        ingress = {
          enabled      = true
          className    = "cilium"
          hosts        = ["hubble.${var.region}.snisid.gouv.local"]
          tls = [
            {
              secretName = "hubble-ui-tls"
              hosts      = ["hubble.${var.region}.snisid.gouv.local"]
            }
          ]
        }
      }
      metrics = {
        enabled = [
          "dns:query",
          "drop",
          "tcp",
          "flow",
          "icmp",
          "http"
        ]
        enableOpenMetrics = true
      }
    }
    prometheus = {
      enabled = true
      port    = 9090
    }
    bandwidthManager = {
      enabled = true
    }
    hostFirewall = {
      enabled = true
    }
    encryption = {
      enabled      = true
      type         = "wireguard"
      nodeEncryption = true
    }
    bpf = {
      masquerade = true
    }
    nodePort = {
      enabled = true
    }
    loadBalancer = {
      algorithm = "maglev"
      mode      = "dsr"
    }
    l7Proxy = true
    egressGateway = {
      enabled = true
    }
    tetragon = {
      enabled = var.enable_tetragon
      tetragon = {
        resources = {
          limits   = { memory = "1Gi", cpu = "500m" }
          requests = { memory = "256Mi", cpu = "100m" }
        }
      }
    }
    policyAuditMode = false
    devices         = "eth+"
    mtu             = 9000
    ipv4NativeRoutingCIDR = var.pod_cidr
    tunnelProtocol        = ""
    routingMode           = "native"
    autoDirectNodeRoutes  = true
    cluster = {
      name = "snisid-${var.region}"
      id   = var.cluster_id
    }
    securityContext = {
      privileged = false
      capabilities = {
        ciliumAgent      = ["CHOWN","KILL","NET_ADMIN","NET_RAW","IPC_LOCK","SYS_ADMIN","SYS_RESOURCE","DAC_OVERRIDE","FOWNER","SETGID","SETUID"]
        cleanCiliumState = ["NET_ADMIN","SYS_ADMIN","SYS_RESOURCE"]
      }
    }
  }
}

# ───────────────────────────────────────────────
# Namespace Cilium (si isolé)
# ───────────────────────────────────────────────
resource "kubernetes_namespace" "cilium" {
  count = var.create_namespace ? 1 : 0
  metadata {
    name = var.namespace
    labels = {
      "snisid.gov/tier"    = var.tier
      "snisid.gov/region"  = var.region
      "snisid.gov/service" = "networking"
    }
  }
}

# ───────────────────────────────────────────────
# Helm Release Cilium
# ───────────────────────────────────────────────
resource "helm_release" "cilium" {
  name       = "cilium"
  repository = "https://helm.cilium.io/"
  chart      = "cilium"
  version    = var.cilium_helm_version
  namespace  = var.namespace
  create_namespace = false

  values = [yamlencode(local.cilium_values)]

  depends_on = [kubernetes_namespace.cilium]
}
