package snisid.banking

default allow := false

# DECRET D'OUVERTURE BANCAIRE (PHASE 10)
# Règle : Autorise la monétisation des requêtes KYC pour les banques privées.

allow {
    input.request.path == "/api/kyc/verify"
    input.user.type == "PRIVATE_BANK"
    input.user.has_active_contract == true
    input.user.quota_remaining > 0
}

deny[msg] {
    input.request.path == "/api/kyc/verify"
    input.user.type == "PRIVATE_BANK"
    input.user.has_active_contract == false
    msg := "VIOLATION DU DECRET : Contrat de monétisation bancaire invalide ou expiré."
}

deny[msg] {
    input.request.path == "/api/kyc/verify"
    input.user.type == "PRIVATE_BANK"
    input.user.quota_remaining <= 0
    msg := "VIOLATION COMMERCIALE : Quota API dépassé."
}
