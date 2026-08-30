# PROMPT 267: AUTOMATED BUILD SYSTEM

This architecture defines the high-performance, secure automated build system for all SNISID microservices, ensuring reproducible and verifiable artifacts.

---

## 1. Build Architecture (Containerized)

SNISID uses a **Multi-Stage Build** model to minimize image size and maximize security.

- **Build Runners**: Ephemeral GitHub Actions runners with SSD-backed storage for high-speed compilation.
- **Languages Supported**:
    - **Go**: Optimized with `CGO_ENABLED=0` for static binaries.
    - **Python**: ML-ready builds with pre-installed GPU drivers and optimized wheel caches.
    - **Java/Kotlin**: Maven/Gradle builds with centralized dependency proxies.

---

## 2. Artifact Workflows

1.  **Trigger**: Code merge to `main` or a tagged release.
2.  **Linting & Formatting**: Enforcement of national coding standards (e.g., `gofmt`, `flake8`).
3.  **Compilation**: Multi-threaded build process using language-specific optimizations.
4.  **Packaging**: Creation of a distroless container image (e.g., `gcr.io/distroless/static`).
5.  **Signing**: Every image is cryptographically signed using **Cosign** and the SNISID private key (stored in HSM).
6.  **Push**: Artifact is pushed to the regional OCI-compliant registry.

---

## 3. Optimization Strategies

- **Dependency Caching**: Distributed caching for `go/pkg/mod`, `~/.cache/pip`, and `.m2/repository` to reduce build times by 70%.
- **Layer Optimization**: Careful ordering of `Dockerfile` instructions to maximize layer reuse across microservices.
- **Parallel Execution**: Independent microservices are built in parallel using GitHub Actions job matrices.

---

## 4. Security Validation Pipelines

Every build must pass a **Security Audit** before the image is pushed:

- **SCA (Software Composition Analysis)**: Checks for known vulnerabilities in third-party dependencies (Trivy/Snyk).
- **Static Analysis**: Deep code inspection for security anti-patterns (Semgrep).
- **SBOM Generation**: Automatic creation of a **CycloneDX SBOM** for every image to track the provenance of all libraries.

---

## 5. Governance Framework

- **Reproducible Builds**: Ensures that the same source code always produces the identical binary bit-for-bit.
- **Immutable Tags**: Build tags are based on Git SHA and Semantic Versioning; `latest` tags are prohibited in production.
- **Cleanup Policy**: Automated pruning of build artifacts older than 90 days, unless they are tagged as "National Archive" for forensic compliance.

---

**PROMPT 267 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 268 — AUTOMATED TESTING PIPELINE.**
