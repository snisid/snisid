# SNISID: Digital Identity Civil Graph

The Digital Identity Civil Graph is the "Social Backbone" of SNISID, responsible for modeling the complex web of family relationships, civil dependencies, and life-cycle events that define a citizen's legal standing within the nation.

---

## 1. Civil Relationship Modeling (Neo4j)

The graph uses **Neo4j** to represent civil connections as a rich network of nodes and edges.

- **Entity Nodes**:
  - `Person`: The core identity node (NID, Biometric Template ID).
  - `Document`: Birth certificates, marriage licenses, death certificates.
  - `Location`: Birthplace, current residence, death location.
- **Relationship Edges**:
  - `PARENT_OF` / `CHILD_OF`: Defining family lineage.
  - `SPOUSE_OF`: Representing legal marriage/unions.
  - `GUARDIAN_OF`: Defining legal responsibility for minors or vulnerable persons.
  - `WITNESS_OF`: Linking identities to civil events (e.g., witnesses at a wedding).

---

## 2. Identity Lifecycle Tracking

The graph tracks every major civil event from registration to certification.

- **Birth Registration**: The "Genesis Event" that creates the initial `Person` node and its primary `PARENT_OF` links.
- **Biometric Maturation**: As a citizen grows, the graph tracks the evolution of their biometric templates, linking them to their original civil identity.
- **Status Transitions**: Real-time updates for status changes: `Single` -> `Married`, `Minor` -> `Adult`, `Alive` -> `Deceased`.
- **Death Certification**: The "Final Event" that marks the identity as `INACTIVE` and triggers the cascade of legal closures (e.g., pension stops, document cancellation).

---

## 3. Civil Intelligence & Fraud Prevention

The civil graph is a powerful tool for detecting sophisticated identity fraud.

- **Lineage Consistency**: Automatically flagging anomalous family trees (e.g., a person having three biological mothers or parents younger than their children).
- **Marriage Fraud Detection**: Identifying "Ghost Marriages" or high-frequency marriage/divorce cycles within specific relationship clusters.
- **Inheritance Fraud Protection**: Ensuring that death certificates are cryptographically linked to a verified identity before any asset transfers are permitted.
- **Guardian Abuse Prevention**: Monitoring the relationships between guardians and multiple non-related minors for potential exploitation patterns.

---

## 4. Integration & Data Synchronization

- **Kafka Ingestion**: Every civil event (from the Ministry of Justice or Civil Registry) is pushed to the `civil.events.v1` topic.
- **Neo4j Streams**: Changes in the civil graph are streamed to the **National Intelligence Fusion Center** for real-time risk enrichment.
- **Sovereign Audit Ledger**: Every modification to a relationship edge is logged as a "Civil Transaction" with full cryptographic traceability.

---

## 5. Security & Privacy

- **RBAC on Relationships**: Access to sensitive family links (e.g., adoption records) is restricted to specialized judicial analysts.
- **Graph Anonymization**: For demographic research, the graph can be exported in an "Isomorphic Anonymized" state where nodes and edges are preserved but PII is masked.
- **Consent-Based Access**: Certain relationship data is only accessible to third parties (e.g., banks) with the explicit, digitally signed consent of the citizen.
