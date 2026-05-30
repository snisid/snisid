# SNISID: Master Deployment Makefile

.PHONY: bootstrap deploy clean test

bootstrap:
	@echo "🚀 Bootstrapping SNISID Infrastructure..."
	@powershell ./scripts/bootstrap.sh

deploy:
	@echo "🏗️ Deploying SNISID to Kubernetes..."
	@kubectl create ns snisid || true
	@helm install snisid-platform ./deployments/helm/snisid -n snisid
	@kubectl apply -f deployments/gitops/app-of-apps.yaml

test:
	@echo "🧪 Running Platform Tests..."
	@go test ./backend/...
	@pytest ai/tests/

clean:
	@echo "🧹 Cleaning up environment..."
	@helm uninstall snisid-platform -n snisid
	@kubectl delete ns snisid
