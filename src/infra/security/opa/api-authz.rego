package envoy.authz

import input.attributes.request.http as http_request

default allow = false

# Helper to extract the JWT payload from the Authorization header
token = payload {
    [_, jwt] := split(http_request.headers["authorization"], " ")
    [_, payload, _] := io.jwt.decode(jwt)
}

# Rule: Allow basic profile read if the user has 'profile.read' scope
allow {
    http_request.method == "GET"
    startswith(http_request.path, "/v1/identity/profile")
    "profile.read" == token.scopes[_]
}

# Rule: Strict Data Minimization - Only specific judicial/security agencies can read biometrics
allow {
    http_request.method == "GET"
    startswith(http_request.path, "/v1/identity/biometrics")
    "biometrics.read" == token.scopes[_]
    
    # Attribute-Based Access Control (ABAC): Only DGI (Intelligence) or DCPJ (Police) can access
    allowed_agencies := {"DGI", "DCPJ"}
    token.agency_id == allowed_agencies[_]
}

# Rule: Consent Management - Only authorized operators can revoke consent
allow {
    http_request.method == "POST"
    startswith(http_request.path, "/v1/consent/revoke")
    "consent.admin" == token.scopes[_]
}
