package snisid.interoperability

test_allow_gateway {
    allow with input as {
        "request": {"path": "/api/identity", "method": "GET"},
        "source_network": "national_api_gateway"
    }
}

test_deny_direct_access {
    not allow with input as {
        "request": {"path": "/api/identity", "method": "GET"},
        "source_network": "ministry_of_health"
    }
}
