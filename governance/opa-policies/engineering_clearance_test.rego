package snisid.engineering

test_allow_secret_production {
    allow with input as {
        "user": {"department": "SNISID_CORE_ENGINEERING", "clearance": "SECRET"},
        "resource": {"environment": "PRODUCTION"}
    }
}

test_deny_unauthorized_production {
    not allow with input as {
        "user": {"department": "SNISID_CORE_ENGINEERING", "clearance": "CONFIDENTIAL"},
        "resource": {"environment": "PRODUCTION"}
    }
}
