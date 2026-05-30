package snisid.authz

default allow = false

allow {
	input.user.role == "admin"
}

allow {
	input.user.role == "investigator"
	input.action == "read"
}

allow {
	input.user.role == "analyst"
	input.resource == "fraud"
	input.action == "analyze"
}
