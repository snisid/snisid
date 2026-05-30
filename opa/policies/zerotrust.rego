package snisid.zerotrust
import future.keywords.if
import future.keywords.in
 
# Politique Zero Trust SNISID - Decret Presidentiel 2026-ZT-001
default allow = false
 
# Regle 1: Services d'urgence - JAMAIS bloques
allow if {
    input.context.emergency == true
    input.user.role in ["PNH_EMERGENCY", "MEDICAL_EMERGENCY", "CIVIL_PROTECTION"]
    input.user.clearance in ["SECRET", "TRES_SECRET"]
}
 
# Regle 2: Agents de terrain (ONI, PNH, DIE) - acces lecture citoyens
allow if {
    input.resource in ["citizen_identity_read", "civil_registry_read", "passport_status"]
    input.action == "read"
    input.user.role in data.snisid.roles_agents_terrain
    input.context.mfa == true
}
 
# Regle 3: Acces production - clearance SECRET requise
allow if {
    input.resource in data.snisid.production_resources
    input.user.clearance in ["SECRET", "TRES_SECRET"]
    input.context.mfa == true
    not input.user.clearance_expired
}
 
# Regle 4: mTLS inter-services (Istio Service Mesh)
allow if {
    input.type == "mtls_service_to_service"
    input.source_certificate.issuer == "SNISID-PKI-ROOT-2026"
    input.source_certificate.san in data.snisid.authorized_service_accounts
    not certificate_expired(input.source_certificate)
}
 
certificate_expired(cert) if {
    time.now_ns() > cert.not_after_ns
}
