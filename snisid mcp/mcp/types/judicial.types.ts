export interface JudicialCase {
  caseId: string;
  court: string;
  status: 'OPEN' | 'CLOSED' | 'SEALED' | 'APPEAL';
  classification: 'CONFIDENTIAL' | 'SECRET';
}

export interface WarrantRecord {
  warrantId: string;
  nationalId: string;
  status: 'ACTIVE' | 'EXECUTED' | 'CANCELLED';
  issuingAuthority: string;
}
