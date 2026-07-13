# SNISID: National ABIS Architecture

The National Automated Biometric Identification System (ABIS) is the foundational "Deduplication Engine" of SNISID, ensuring the principle of **"One Person, One Identity"** across the entire population.

---

## 1. Population-Scale 1:N Matching

The ABIS performs exhaustive 1:N (one-to-many) searches to identify duplicate enrollments or fraudulent identity claims.

- **Vector Similarity Search**: Biometric templates are stored as high-dimensional vectors in a distributed vector database (e.g., **Milvus**, **Qdrant**, or **Elasticsearch with Vector Search**).
- **Indexing Strategy**: Uses **HNSW (Hierarchical Navigable Small World)** or **IVF_PQ (Inverted File with Product Quantization)** indexing to achieve sub-second search latency across tens of millions of records.
- **Cascaded Matching**:
  1. **Primary Filter**: Rapid, low-precision vector search to identify the top 100 potential matches.
  2. **Secondary Matcher**: High-precision, modality-specific comparison (e.g., Minutiae-based fingerprint matching) for the top candidates.

---

## 2. Distributed Deduplication Engine

- **Horizontal Scalability**: The matching engine is sharded across multiple nodes. Each node handles a subset of the national biometric database.
- **GPU-Accelerated Searching**: Utilizes GPU clusters to perform massively parallel distance calculations (Cosine Similarity or Euclidean Distance).
- **Kafka-Driven Ingestion**: New enrollments are queued in Kafka, allowing the ABIS to process deduplication requests asynchronously and maintain a consistent throughput.

---

## 3. Adjudication & Exception Management

When the ABIS identifies a potential duplicate, the system triggers a formal **Adjudication Workflow**.

- **Automatic Hits**: If the match score exceeds a very high threshold (e.g., 0.999), the system automatically flags the enrollment as a "Duplicate" and denies the request.
- **Manual Adjudication**: If the score falls within a "Gray Zone," the case is escalated to a human **Biometric Expert** for visual verification.
- **Auditability**: Every adjudication decision (Automatic or Manual) is cryptographically signed and stored in the **Sovereign Audit Ledger**, including the raw similarity scores and the expert's credentials.

---

## 4. Secure Biometric Storage

- **Data Partitioning**: Biometric vectors are stored separately from demographic data. Only the **National ID (NID)** acts as the link between the two.
- **Encryption at Rest**: Templates are encrypted using **AES-256-GCM** with keys managed by the national HSM cluster.
- **Anonymized Vectors**: The vectors themselves are "De-Identified" so that even a database leak would not reveal the original biometric traits without the matching algorithm's specific parameters.

---

## 5. Performance & Reliability

- **Latency Target**: < 2 seconds for a full 1:N national search.
- **Accuracy Target**: Equal to or better than **NIST FRVT/MINEX** standards.
- **High Availability**: Multi-region deployment with real-time replication of the vector index.
