# SNISID — Terraform Module : Ceph Storage Cluster (Rook Operator)
# Classification: RESTREINT DEFENSE
# Role: Provision pools, StorageClasses, et multi-site RGW via Rook CRDs
# Dépendance: cluster Kubernetes existant + nodes labellisés storage

terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.11.0"
    }
  }
}

locals {
  ceph_cluster_name = "snisid-national-${var.region}"
  pools = {
    for p in var.pools : p.name => p
  }
}

# ───────────────────────────────────────────────
# Rook Ceph Operator (Helm)
# ───────────────────────────────────────────────
resource "helm_release" "rook_ceph_operator" {
  name       = "rook-ceph"
  repository = "https://charts.rook.io/release"
  chart      = "rook-ceph"
  version    = "1.13.0"
  namespace  = var.namespace
  create_namespace = true

  set {
    name  = "image.repository"
    value = "${var.registry}/rook/ceph"
  }
  set {
    name  = "image.tag"
    value = "v1.13.0"
  }
  set {
    name  = "resources.requests.memory"
    value = "512Mi"
  }
  set {
    name  = "resources.limits.memory"
    value = "1Gi"
  }
  set {
    name  = "enableDiscoveryDaemon"
    value = "true"
  }
  set {
    name  = "csi.enableCephfsDriver"
    value = "true"
  }
  set {
    name  = "csi.enableRbdDriver"
    value = "true"
  }
  set {
    name  = "csi.enableSnapshotter"
    value = "true"
  }
  set {
    name  = "csi.provisionerReplicas"
    value = "2"
  }

  depends_on = [kubernetes_namespace.rook_ceph]
}

resource "kubernetes_namespace" "rook_ceph" {
  metadata {
    name = var.namespace
    labels = {
      "snisid.gov/tier"    = var.tier
      "snisid.gov/region"  = var.region
      "snisid.gov/service" = "storage"
      "pod-security.kubernetes.io/enforce"        = "privileged"
      "pod-security.kubernetes.io/enforce-version" = "latest"
    }
  }
}

# ───────────────────────────────────────────────
# CephCluster CR (via kubernetes_manifest pour CRD)
# ───────────────────────────────────────────────
resource "kubernetes_manifest" "ceph_cluster" {
  manifest = {
    apiVersion = "ceph.rook.io/v1"
    kind       = "CephCluster"
    metadata = {
      name      = local.ceph_cluster_name
      namespace = var.namespace
      labels = {
        "snisid.gov/tier"    = var.tier
        "snisid.gov/region"  = var.region
        "snisid.gov/service" = "storage"
      }
    }
    spec = {
      cephVersion = {
        image            = "${var.registry}/ceph/ceph:v18.2.1"
        allowUnsupported = false
      }
      dataDirHostPath = "/var/lib/rook"
      mon = {
        count                = var.mon_count
        allowMultiplePerNode = false
      }
      mgr = {
        count = 2
        modules = [
          { name = "pg_autoscaler", enabled = true },
          { name = "rook", enabled = true }
        ]
      }
      dashboard = {
        enabled = true
        ssl     = true
      }
      monitoring = {
        enabled = true
      }
      network = {
        provider     = "host"
        connections = {
          requireMsgr2 = true
          compression = { enabled = true }
          encryption  = { enabled = true }
        }
      }
      storage = {
        useAllNodes       = false
        useAllDevices     = false
        nodes             = var.storage_nodes
        resources = {
          osd = {
            limits = { memory = "8Gi", cpu = "2000m" }
            requests = { memory = "4Gi", cpu = "1000m" }
          }
        }
      }
      placement = {
        osd = {
          nodeAffinity = {
            requiredDuringSchedulingIgnoredDuringExecution = {
              nodeSelectorTerms = [
                {
                  matchExpressions = [
                    {
                      key      = "snisid.gov/storage-node"
                      operator = "In"
                      values   = ["true"]
                    }
                  ]
                }
              ]
            }
          }
        }
      }
      disruptionManagement = {
        managePodBudgets      = true
        osdMaintenanceTimeout = 30
      }
      crashCollector = {
        disable = false
      }
      logCollector = {
        enabled    = true
        periodicity = "daily"
        maxLogSize = "500M"
      }
    }
  }

  depends_on = [helm_release.rook_ceph_operator]
}

# ───────────────────────────────────────────────
# CephBlockPools (dynamiques selon variables)
# ───────────────────────────────────────────────
resource "kubernetes_manifest" "ceph_blockpool" {
  for_each = local.pools

  manifest = {
    apiVersion = "ceph.rook.io/v1"
    kind       = "CephBlockPool"
    metadata = {
      name      = each.value.name
      namespace = var.namespace
    }
    spec = {
      failureDomain = each.value.failure_domain
      replicated = {
        size                  = each.value.replicated_size
        requireSafeReplicaSize = true
      }
      deviceClass       = lookup(each.value, "device_class", "ssd")
      compressionMode   = lookup(each.value, "compression_mode", "aggressive")
      application       = "rbd"
    }
  }

  depends_on = [kubernetes_manifest.ceph_cluster]
}

# ───────────────────────────────────────────────
# StorageClasses RBD
# ───────────────────────────────────────────────
resource "kubernetes_storage_class" "ceph_rbd" {
  for_each = local.pools

  metadata {
    name = each.value.storage_class_name
    annotations = {
      "storageclass.kubernetes.io/is-default-class" = lookup(each.value, "is_default", false) ? "true" : "false"
    }
  }
  storage_provisioner    = "rook-ceph.rbd.csi.ceph.com"
  reclaim_policy         = lookup(each.value, "reclaim_policy", "Delete")
  allow_volume_expansion = true
  volume_binding_mode    = "WaitForFirstConsumer"
  parameters = {
    clusterID                 = local.ceph_cluster_name
    pool                      = each.value.name
    imageFormat               = "2"
    imageFeatures             = "layering,fast-diff,deep-flatten,object-map"
    "csi.storage.k8s.io/provisioner-secret-name"      = "rook-csi-rbd-provisioner"
    "csi.storage.k8s.io/provisioner-secret-namespace"  = var.namespace
    "csi.storage.k8s.io/node-stage-secret-name"         = "rook-csi-rbd-node"
    "csi.storage.k8s.io/node-stage-secret-namespace"   = var.namespace
  }

  depends_on = [kubernetes_manifest.ceph_blockpool]
}

# ───────────────────────────────────────────────
# CephFilesystem (partagé national)
# ───────────────────────────────────────────────
resource "kubernetes_manifest" "ceph_filesystem" {
  count = var.enable_filesystem ? 1 : 0

  manifest = {
    apiVersion = "ceph.rook.io/v1"
    kind       = "CephFilesystem"
    metadata = {
      name      = "${local.ceph_cluster_name}-fs"
      namespace = var.namespace
    }
    spec = {
      metadataPool = {
        replicated = { size = 3 }
      }
      dataPools = [
        {
          name = "replicated"
          replicated = { size = 3 }
          compressionMode = "aggressive"
        }
      ]
      metadataServer = {
        activeCount   = 2
        activeStandby = true
      }
    }
  }

  depends_on = [kubernetes_manifest.ceph_cluster]
}

# ───────────────────────────────────────────────
# CephObjectStore (S3 compatible national)
# ───────────────────────────────────────────────
resource "kubernetes_manifest" "ceph_object_store" {
  count = var.enable_object_store ? 1 : 0

  manifest = {
    apiVersion = "ceph.rook.io/v1"
    kind       = "CephObjectStore"
    metadata = {
      name      = "${local.ceph_cluster_name}-s3"
      namespace = var.namespace
    }
    spec = {
      metadataPool = {
        failureDomain = "rack"
        replicated = { size = 3 }
      }
      dataPool = {
        failureDomain = "rack"
        erasureCoded = {
          dataChunks   = 4
          codingChunks = 2
        }
      }
      gateway = {
        instances = 2
        placement = {
          podAntiAffinity = {
            requiredDuringSchedulingIgnoredDuringExecution = [
              {
                labelSelector = {
                  matchLabels = {
                    app = "rook-ceph-rgw"
                  }
                }
                topologyKey = "topology.kubernetes.io/zone"
              }
            ]
          }
        }
        resources = {
          limits = { memory = "4Gi", cpu = "2000m" }
          requests = { memory = "2Gi", cpu = "1000m" }
        }
      }
      zone = {
        name = "${local.ceph_cluster_name}-zone"
      }
    }
  }

  depends_on = [kubernetes_manifest.ceph_cluster]
}
