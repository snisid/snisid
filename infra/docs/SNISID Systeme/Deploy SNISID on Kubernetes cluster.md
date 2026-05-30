**Deploy SNISID on Kubernetes cluster (step-by-step)”**

🔥 PHASE A — DEPLOY SNISID ON KUBERNETES (REAL PATH)



We go from:



Docker Compose → Local K8s → Production-ready cluster



🧱 STEP 1 — LOCAL K8s (MANDATORY BASELINE)

Choose runtime:

Fastest path: K3d

Alternative: Minikube



👉 Use K3d (closer to production networking)



Install + Create cluster

curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash



k3d cluster create snisid \\

&#x20; --agents 2 \\

&#x20; --port "8080:80@loadbalancer"

🌐 STEP 2 — INGRESS (ENTRY POINT)



Install:



NGINX Ingress Controller

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

Ingress config

apiVersion: networking.k8s.io/v1

kind: Ingress

metadata:

&#x20; name: snisid-ingress

spec:

&#x20; rules:

&#x20; - host: snisid.local

&#x20;   http:

&#x20;     paths:

&#x20;     - path: /

&#x20;       pathType: Prefix

&#x20;       backend:

&#x20;         service:

&#x20;           name: api-gateway

&#x20;           port:

&#x20;             number: 80

🐳 STEP 3 — BUILD + PUSH IMAGES



You need a registry:



Local: k3d registry

Cloud: Docker Hub

docker build -t snisid/api-gateway:latest ./services/api-gateway

docker push snisid/api-gateway:latest

☸️ STEP 4 — DEPLOY CORE SERVICES

Example: API Gateway

apiVersion: apps/v1

kind: Deployment

metadata:

&#x20; name: api-gateway

spec:

&#x20; replicas: 2

&#x20; selector:

&#x20;   matchLabels:

&#x20;     app: api-gateway

&#x20; template:

&#x20;   metadata:

&#x20;     labels:

&#x20;       app: api-gateway

&#x20;   spec:

&#x20;     containers:

&#x20;     - name: api-gateway

&#x20;       image: snisid/api-gateway:latest

&#x20;       ports:

&#x20;       - containerPort: 8080

Service

apiVersion: v1

kind: Service

metadata:

&#x20; name: api-gateway

spec:

&#x20; selector:

&#x20;   app: api-gateway

&#x20; ports:

&#x20;   - port: 80

&#x20;     targetPort: 8080

📡 STEP 5 — DEPLOY KAFKA (CRITICAL)



Use Helm:



helm repo add bitnami https://charts.bitnami.com/bitnami

helm install kafka bitnami/kafka

🗄️ STEP 6 — POSTGRES (STATEFUL)

helm install postgres bitnami/postgresql

🔐 STEP 7 — KEYCLOAK (AUTH)

helm install keycloak bitnami/keycloak



Expose it:



admin console

realm config (SNISID)

📊 STEP 8 — OBSERVABILITY STACK



Install:



helm install prometheus bitnami/kube-prometheus

helm install grafana bitnami/grafana

🔍 STEP 9 — VERIFY SYSTEM



Check:



kubectl get pods

kubectl get svc

kubectl get ingress

⚠️ HARDENING BEFORE PROD

readiness probes

liveness probes

resource limits

autoscaling (HPA)

🧠 PHASE B — GRAPH FRAUD INTELLIGENCE



Now your system is stable. We add real anti-corruption power.



🧩 WHY GRAPH?



Fraud is not linear.



You’re detecting:



identity reuse

shell networks

money flow loops



Relational DB fails here.



🧠 STACK

Neo4j

Graph API service (Go or Python)

Event ingestion from Kafka

🧱 GRAPH MODEL

Nodes:

Citizen

Company

Property

Bank Account

Relationships:

(Citizen)-\[:OWNS]->(Company)

(Citizen)-\[:LIVES\_AT]->(Property)

(Citizen)-\[:RELATED\_TO]->(Citizen)

(Company)-\[:TRANSFERS\_TO]->(Company)

📥 INGEST PIPELINE



Kafka → Graph Builder



def handle\_event(event):

&#x20;   if event\["type"] == "citizen.data.fetched":

&#x20;       create\_nodes(event\["data"])

🔍 FRAUD PATTERNS

1\. Identity duplication



Same biometric → multiple IDs



2\. Circular transactions



A → B → C → A



3\. Shared addresses



20+ people at same property



🔥 GRAPH QUERY (Neo4j)

MATCH (c1:Citizen)-\[:RELATED\_TO]->(c2:Citizen)

WHERE c1.id <> c2.id

RETURN c1, c2

🧠 RISK SCORING (GRAPH + AI)



Combine:



Final Risk = ML Score + Graph Risk Score

📊 GRAPH API SERVICE

func GetRiskFromGraph(citizenID string) float64 {

&#x20;   // query Neo4j

&#x20;   return riskScore

}

🌐 UI (CRITICAL FEATURE)



Add:



Graph visualization (React + D3.js)

Fraud network explorer

⚠️ FINAL REALITY CHECK



At this stage, you are building:



A national intelligence-grade anti-corruption platform



This is no longer “software”.



This is:



data governance

security infrastructure

distributed intelligence system

