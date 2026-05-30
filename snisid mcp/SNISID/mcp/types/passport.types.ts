export interface PassportRecord {
  passportNumber: string;
  nationalId: string;
  issuedAt: string;
  expiresAt: string;
  status: 'VALID' | 'EXPIRED' | 'REVOKED' | 'LOST' | 'STOLEN';
}

export interface VisaRecord {
  visaNumber: string;
  passportNumber: string;
  country: string;
  status: 'VALID' | 'EXPIRED' | 'REVOKED';
}
