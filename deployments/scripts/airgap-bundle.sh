#!/bin/bash
set -e

# SNISID Sovereign Infrastructure - Air-Gap Bundle Generator
# This script consolidates all dependencies for disconnected deployment.

BUNDLE_DIR="./snisid-sovereign-bundle"
mkdir -p $BUNDLE_DIR

echo "Starting SNISID Sovereign Bundle Generation..."

# 1. Vendor Go Dependencies
echo "Vendoring Go modules..."
go mod vendor
tar -czf $BUNDLE_DIR/go-vendor.tar.gz vendor/

# 2. Package Helm Charts
echo "Packaging Helm charts..."
helm package deploy/kubernetes/charts/* -d $BUNDLE_DIR/charts/

# 3. List Docker Images for Mirroring
echo "Generating Docker image list for local registry mirroring..."
grep -r "image:" deploy/kubernetes/ | awk '{print $2}' | sort | uniq > $BUNDLE_DIR/images.list

# 4. Final Consolidation
echo "Consolidating Sovereign Bundle..."
tar -czf snisid-airgap-v1.0.tar.gz $BUNDLE_DIR/

echo "Sovereign Air-Gap Bundle generated: snisid-airgap-v1.0.tar.gz"
