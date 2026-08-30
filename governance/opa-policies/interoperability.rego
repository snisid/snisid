package snisid.interoperability

import future.keywords.in

default allow := false

# DECRET D'INTEROPERABILITE (PHASE 4)
# Règle : Il est interdit pour un ministère de créer une base de données d'identité locale.
# Tout accès aux données d'identité doit passer par l'API Gateway Nationale.

allow {
    input.request.path == "/api/identity"
    input.request.method in ["GET", "POST"]
    input.source_network == "national_api_gateway"
}

deny[msg] {
    input.request.path == "/api/identity"
    input.source_network != "national_api_gateway"
    msg := "VIOLATION DU DECRET : Accès direct à la base d'identité refusé. Passage obligatoire par la National API Gateway."
}
