# SNISID API Audit Fabric

## Objective
Provide a tamper-proof trail of all data exchanges between government agencies for accountability and forensic analysis.

## What is Audited?
- **All API Calls**: Caller ID, timestamp, endpoint, status code.
- **Data Access**: Which specific record was accessed (NIU, NIF).
- **Security Events**: Login/Logout, token generation, policy changes.
- **Failures**: Detailed stack traces for system errors.

## Storage and Integrity
- **Immutable Logs**: Audit logs are stored in a write-once-read-many (WORM) storage.
- **Digital Signatures**: Each audit entry is cryptographically signed.
- **Retention**: Minimum 10 years for judicial and identity-related audits.

## Search and Analysis
- Audit logs are indexed for rapid search by authorized personnel.
- AI-driven anomaly detection to identify suspicious patterns of data access.
