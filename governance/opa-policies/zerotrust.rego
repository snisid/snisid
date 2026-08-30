package snisid.zerotrust

import future.keywords.in

default allow := false

# DECRET ZERO TRUST (PHASE 6)
# Règle : Le CISO National a l'autorité de déconnecter n'importe quel ministère (Drop Traffic)
# si ses systèmes ne sont pas patchés ou si IAM est compromis.

allow {
    input.user.role == "NATIONAL_CISO"
    input.action == "ISOLATE_NETWORK"
}

deny[msg] {
    input.user.role != "NATIONAL_CISO"
    input.action == "ISOLATE_NETWORK"
    msg := "VIOLATION DU DECRET : Seul le CISO National a l'autorité de déclencher un isolement Zero Trust inter-ministériel."
}

# Bloquer le trafic d'un ministère isolé
deny[msg] {
    input.source_ministry in data.isolated_ministries
    msg := "VIOLATION DE SECURITE : Ce ministère a été isolé par le CISO National (Décret Phase 6). Trafic bloqué."
}
