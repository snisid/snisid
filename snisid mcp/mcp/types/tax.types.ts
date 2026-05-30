export interface TaxpayerRecord {
  nif: string;
  nationalId?: string;
  businessName?: string;
  status: 'COMPLIANT' | 'NON_COMPLIANT' | 'UNDER_REVIEW';
}

export interface FinancialRiskResult {
  score: number;
  flags: string[];
  explanation: string;
}
