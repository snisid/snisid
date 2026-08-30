export interface CitizenProfile {
  nationalId: string;
  firstName: string;
  lastName: string;
  dateOfBirth: string;
  placeOfBirth?: string;
  nationality: string;
  status: 'ACTIVE' | 'SUSPENDED' | 'DECEASED' | 'UNDER_REVIEW';
}

export interface IdentityVerificationResult {
  verified: boolean;
  confidence: number;
  nationalId?: string;
  riskFlags: string[];
  dataMinimized: true;
}

export interface BiometricMatchResult {
  match: boolean;
  confidence: number;
  modality: 'FACE' | 'FINGERPRINT' | 'IRIS' | 'MULTI';
  referenceId?: string;
}
