package snisid.zerotrust

test_admin_allowed {
    allow with input as {"role": "admin"}
}

test_authenticated_high_clearance_allowed {
    allow with input as {"authenticated": true, "clearance_level": 3}
}

test_unauthenticated_denied {
    not allow with input as {"authenticated": false, "clearance_level": 1}
}
