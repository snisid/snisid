package snisid.zerotrust.access

import data.snisid.risk_scores
import data.snisid.device_inventory

# Default deny
default allow = false

# Continuous Verification Logic
allow {
    # 1. Identity Verification
    valid_identity
    
    # 2. Context / Risk Scoring (Must be low or medium risk)
    acceptable_risk
    
    # 3. Device Posture Verification
    healthy_device
}

valid_identity {
    input.request.auth.claims.iss == "https://iam.snisid.gov/realms/national"
    input.request.auth.claims.acr == "Level3" # Enforces MFA / FIDO2
}

acceptable_risk {
    # The user's dynamic risk score must be under 75
    # This score is continuously calculated by the SIEM ML jobs
    user_risk := risk_scores[input.request.auth.claims.sub]
    user_risk < 75
}

healthy_device {
    # The device presenting the certificate must be known to the MDM
    # and marked as compliant (e.g., patched, encrypted, no malware)
    device := device_inventory[input.request.headers["x-device-id"]]
    device.status == "compliant"
    device.os_updated == true
}
