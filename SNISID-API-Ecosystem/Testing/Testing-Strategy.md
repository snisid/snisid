# SNISID National API Testing Platform

## Test Suites

### 1. Functional Testing
- **Tools**: Postman, Newman, Pytest.
- **Scope**: Verify that API endpoints return correct data for valid inputs.

### 2. Security Testing (DAST/SAST)
- **Tools**: OWASP ZAP, Snyk.
- **Scope**: Check for SQL injection, XSS, insecure headers, and broken authentication.

### 3. Load & Performance Testing
- **Tools**: k6, JMeter.
- **Scope**: Ensure the API can handle peaks of 10,000 requests per second.

### 4. Chaos Testing
- **Tools**: Chaos Mesh, Gremlin.
- **Scope**: Simulate network failures, pod crashes, and high latency to verify resilience.

### 5. Contract Testing
- **Tools**: Pact.
- **Scope**: Ensure that changes in a producer's API don't break consumer services.

## CI/CD Integration
- No API can be deployed to Production without passing all automated tests in the pipeline.
