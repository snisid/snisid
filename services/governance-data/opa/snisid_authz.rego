package snisid.authz

import future.keywords.if
import future.keywords.in

# Default: DENY ALL
default allow = false

# Rule: Investigators can read citizen data IF they have a case and high risk
allow if {
	input.role == "investigator"
	input.action == "read"
	input.context.case_id != ""
	input.context.risk_score > 0.7
	input.context.justification != ""
}

# Rule: Auditors can read everything for audit purposes
allow if {
	input.role == "auditor"
	input.action == "read"
	input.context.justification != ""
}

# Rule: System-to-System access for Risk Engine
allow if {
	input.subject == "risk-engine"
	input.action == "enrich"
	input.mtls_verified == true
}

# Enforce field-level security
allowed_fields[field] {
	some field in ["name", "risk_score", "identity_id"]
}

allowed_fields["*"] if {
	input.role == "admin"
}

# Deny reasons
deny[msg] {
	input.context.justification == ""
	msg := "JUSTIFICATION_REQUIRED"
}

deny[msg] {
	input.role == "investigator"
	input.context.case_id == ""
	msg := "CASE_ID_REQUIRED_FOR_INVESTIGATION"
}
