package snisid.zerotrust

default allow = false

allow {
    input.authenticated == true
    input.clearance_level >= 3
}

allow {
    input.role == "admin"
}
