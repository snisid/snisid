package snisid.authz

import input.attributes.request.http as http_request

default allow = false

# Allow identity creation only for users with 'officer' role from authorized agencies
allow {
    http_request.method == "POST"
    http_request.path == "/api/v1/identities"
    token.payload.role == "officer"
    is_authorized_agency(token.payload.agency_id)
}

# Helper to verify agency authorization
is_authorized_agency(agency_id) {
    # In production, this would query a data source or OPA bundle
    authorized_agencies := {"AGENCY-PRP", "AGENCY-CAP", "AGENCY-JAC"}
    authorized_agencies[agency_id]
}

token = {"payload": payload} {
    [_, payload, _] := io.jwt.decode(http_request.headers.authorization)
}
