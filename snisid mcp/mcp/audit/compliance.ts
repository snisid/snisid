export interface ComplianceControl {
  id: string;
  name: string;
  evidence: string;
  status: 'PASS' | 'FAIL' | 'MANUAL_REVIEW';
}

export function baselineComplianceControls(): ComplianceControl[] {
  return [
    { id: 'SNISID-ZT-001', name: 'Deny by default', evidence: 'RBAC and middleware', status: 'PASS' },
    { id: 'SNISID-AUD-001', name: 'Immutable audit trail', evidence: 'hash chained JSONL', status: 'PASS' },
    { id: 'SNISID-DB-001', name: 'No direct DB access from MCP tools', evidence: 'service API clients only', status: 'PASS' },
    { id: 'SNISID-MFA-001', name: 'MFA for sensitive permissions', evidence: 'security/auth.ts', status: 'PASS' }
  ];
}
