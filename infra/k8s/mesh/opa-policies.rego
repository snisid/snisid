package snisid.mesh.authz

import future.keywords.if

# Default: DENY ALL NETWORK TRAFFIC
default allow = false

# Rule: Services must have a verified SPIFFE identity
allow if {
	input.attributes.source.principal != ""
	startswith(input.attributes.source.principal, "spiffe://snisid.gov/")
}

# Rule: Identity service only accessible by Nexus Orchestrator and API Gateway
allow if {
	input.attributes.destination.principal == "spiffe://snisid.gov/ns/identity/sa/identity-service"
	input.attributes.source.principal in [
		"spiffe://snisid.gov/ns/nexus/sa/nexus-orchestrator",
		"spiffe://snisid.gov/ns/gateway/sa/api-gateway"
	]
}

# Rule: mTLS is MANDATORY for all cross-service communication
allow if {
	input.attributes.connection.mtls == true
}

# Audit: Log all denied requests for SIEM integration
deny_audit[msg] if {
	not allow
	msg := sprintf("UNAUTHORIZED_MESH_ACCESS: source=%s destination=%s", [
		input.attributes.source.principal,
		input.attributes.destination.principal
	])
}
