package snisid.engineering

default allow := false

# DECRET HUMAN CAPACITY (PHASE 16)
# Règle : L'accès aux environnements de production (Zero Trust / Core SNISID)
# est strictement réservé au personnel d'ingénierie possédant l'accréditation "SECRET".

allow {
    input.user.department == "SNISID_CORE_ENGINEERING"
    input.user.clearance == "SECRET"
    input.resource.environment == "PRODUCTION"
}

deny[msg] {
    input.resource.environment == "PRODUCTION"
    input.user.clearance != "SECRET"
    msg := "ACCES REFUSE (Phase 16) : Accréditation SECRET requise pour accéder à la production SNISID. Risque de sécurité nationale."
}

# Les ingénieurs avec une clearance standard (ex: CONFIDENTIAL) peuvent accéder au STAGING
allow {
    input.user.department == "SNISID_CORE_ENGINEERING"
    input.user.clearance == "CONFIDENTIAL"
    input.resource.environment == "STAGING"
}
