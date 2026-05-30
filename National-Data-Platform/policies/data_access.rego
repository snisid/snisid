package data.access

default allow = false
allow {
    input.user.role == "auditor"
}

