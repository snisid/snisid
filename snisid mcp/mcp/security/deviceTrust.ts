export interface DeviceTrustResult {
  trusted: boolean;
  riskScore: number;
  reasons: string[];
}

const revokedDevices = new Set<string>();

export function revokeDevice(deviceId: string): void {
  revokedDevices.add(deviceId);
}

export function assessDeviceTrust(deviceId: string, userAgent?: string): DeviceTrustResult {
  const reasons: string[] = [];
  let riskScore = 0;
  if (!deviceId || deviceId.length < 8) {
    riskScore += 50;
    reasons.push('weak_or_missing_device_id');
  }
  if (revokedDevices.has(deviceId)) {
    riskScore += 100;
    reasons.push('revoked_device');
  }
  if (userAgent && /curl|bot|scanner/i.test(userAgent)) {
    riskScore += 20;
    reasons.push('suspicious_user_agent');
  }
  return { trusted: riskScore < 70, riskScore, reasons };
}
