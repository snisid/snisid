# SNISID Interoperability Standards

## 1. Protocol Standards
- **REST**: Default for most integrations. Must follow RESTful principles.
- **gRPC**: Used for high-performance, low-latency internal communications.
- **Async (Events)**: Kafka-based for decoupled systems.

## 2. Data Formats
- **JSON**: Primary data interchange format.
- **Protobuf**: For gRPC services.
- **UTF-8**: Encoding for all text data.

## 3. API Design Rules
- **Versioning**: Versioning in the URL (`/v1/`).
- **Naming**: Use kebab-case for URLs (`/api/identity-records`).
- **HTTP Methods**:
  - `GET`: Read
  - `POST`: Create
  - `PUT`: Update (Full)
  - `PATCH`: Update (Partial)
  - `DELETE`: Remove
- **Error Handling**: Standard HTTP status codes (200, 201, 400, 401, 403, 404, 500).

## 4. Documentation
- **OpenAPI 3.0**: Mandatory for all REST APIs.
- **AsyncAPI**: Mandatory for event-driven systems.
