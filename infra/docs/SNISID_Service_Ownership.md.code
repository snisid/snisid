# SNISID: Service Ownership & Accountability Model

To maintain operational agility and prevent organizational bottlenecks, SNISID employs a strict Service Ownership Model. This ensures every microservice has a distinct domain owner, clear boundaries, and definitive accountability for business logic and data.

---

## 1. Service Ownership Matrix

The following matrix maps each microservice to its organizational owner, primary business logic, and authoritative data store.

| Microservice | Owning Domain / Team | Primary Business Logic | Authoritative Data Store |
| :--- | :--- | :--- | :--- |
| **API Gateway Service** | Edge / DevSecOps Team | External ingress, TLS termination, routing, rate limiting, and WAF inspection. | *None (Stateless)* |
| **Authentication Service** | Identity & Access Team | Credential validation, JWT generation, mTLS certificate verification. | Keycloak DB / Vault |
| **Identity Service** | Core Business Team | Citizen lifecycle management (Enroll, Update, Suspend), demographic logic. | **PostgreSQL** (Identity State) |
| **Streaming Service** | Data Engineering Team | Topic orchestration, guaranteed event delivery, schema registry enforcement. | **Kafka** (Event Log) |
| **Fraud Detection Service** | Risk & Compliance Team | Executing business rules, velocity checks, and calculating aggregate risk scores. | **Redis** (Scoring State) |
| **AI Inference Service** | MLOps / AI Team | Executing ML models for biometric matching (1:N, 1:1) and liveness detection. | *Model Weights Storage* |
| **Graph Intelligence Service** | Intelligence / Data Team | Ingesting events to map familial, financial, and geographic relationships. | **Neo4j** (Relationship Graph) |
| **SOC Alert Service** | Cybersecurity (SOC) Team | Translating system anomalies into MITRE ATT&CK alerts and triggering SOAR. | Elastic / SIEM |

---

## 2. Responsibility Boundaries (Anti-Overlap)

To prevent the "God Service" anti-pattern and avoid overlapping logic, strict boundaries are enforced.

### Identity vs. Authentication
*   **Identity Service** knows *about* the citizen (name, address, status). It does **not** know how the citizen logs in.
*   **Authentication Service** knows *how* to verify the citizen (passwords, OTPs). It does **not** care about the citizen's address or demographic status.

### Fraud Detection vs. AI Inference
*   **AI Inference** performs pure math. It returns a score stating "These two faces match with 98% confidence." It does **not** know if 98% is high enough to approve a passport.
*   **Fraud Detection** applies the business logic. It receives the 98% confidence score from the AI service and applies the national policy rule (e.g., "Passports require > 99% confidence"), thereby making the final "Approve/Reject" decision.

### Identity vs. Graph Intelligence
*   **Identity Service** stores flat, relational records. It knows "John Doe lives at 123 Main St."
*   **Graph Intelligence Service** maps structural context. It knows "15 different identities have claimed 123 Main St in the last 24 hours, forming a high-risk cluster." The Identity Service is entirely unaware of this clustering.

---

## 3. Conflict Resolution Rules Between Services

When business requirements seem to cross service boundaries, the following strict architectural rules dictate how the conflict is resolved:

### Rule 1: The "Source of Truth" Preeminence
No service is allowed to directly read or write to another service's database. If the Fraud Service needs a citizen's date of birth, it **must** query the Identity Service API. If the database schema changes, only the Identity team needs to update their code.

### Rule 2: Synchronous vs. Asynchronous Decoupling
If an operation requires multiple services, prefer asynchronous orchestration unless immediate consistency is required for safety.
*   *Conflict:* "We shouldn't issue an ID if the fraud score is too high."
*   *Resolution:* Do not make the Identity Service wait synchronously for the Fraud Service (which might be doing heavy AI tasks). Instead, the Identity Service creates the ID in a `PENDING_REVIEW` state and publishes an event. The Fraud Service completes its check and sends an `ApproveIdentityCommand` back to the Identity Service to change the state to `ACTIVE`.

### Rule 3: The BFF (Backend-For-Frontend) / Orchestrator Rule
If a new user interface or government agency requires a complex payload containing Identity data, Fraud data, and Graph data, **do not** add custom logic to the Identity service to fetch the other data.
*   *Resolution:* The API Gateway or a dedicated Orchestrator service fetches the data independently from Identity, Fraud, and Graph, aggregates the JSON, and sends it to the client.

### Rule 4: API Contract Immutability
Service owners are completely autonomous regarding their internal tech stack (e.g., the AI team can switch from Python to C++). However, their external API contracts (OpenAPI/gRPC) are treated as legally binding agreements. Breaking changes to an API contract require a formal deprecation period and approval from the Inter-Agency Coordination Board.
