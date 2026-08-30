package snisid.bio_adn.access

import future.keywords.if
import future.keywords.in

default allow := false
default deny := false

# ── ROLES ───────────────────────────────────────────────────────────────────

# bio.lab.technician  — Create/modify STR profiles for own lab only
# bio.lab.supervisor   — Validate profiles, upload to SDIS
# bio.sdis.operator    — Read all dept profiles, SDIS matching
# bio.ndis.analyst     — NDIS matching, read hits, no identity_links
# bio.dcpj.investigator — Read wanted/missing/gang, create warrants
# bio.dcpj.director    — Access identity_links, declassification
# bio.admin            — Technical administration only
# bio.auditor          — Read-only audit log

# ── STR Profiles ────────────────────────────────────────────────────────────

allow if {
	input.action == "create"
	input.resource == "bio_str_profiles"
	"bio.lab.technician" in input.user.roles
	input.user.lab_id == input.record.lab_id
}

allow if {
	input.action == "read"
	input.resource == "bio_str_profiles"
	"bio.lab.technician" in input.user.roles
	input.user.lab_id == input.record.lab_id
}

allow if {
	input.action == "read"
	input.resource == "bio_str_profiles"
	"bio.lab.supervisor" in input.user.roles
	input.user.lab_id == input.record.lab_id
}

allow if {
	input.action == "read"
	input.resource == "bio_str_profiles"
	"bio.sdis.operator" in input.user.roles
	input.user.sdis_code == input.record.sdis_code
}

allow if {
	input.action in {"read", "match"}
	input.resource == "bio_str_profiles"
	"bio.ndis.analyst" in input.user.roles
	input.context.case_number != ""
	input.context.purpose in {"criminal_investigation", "missing_person", "identification"}
}

# ── Identity Links (restreint DCPJ Director) ────────────────────────────────

allow if {
	input.action in {"read", "create"}
	input.resource == "bio_identity_links"
	"bio.dcpj.director" in input.user.roles
	input.context.court_order_ref != ""
	input.context.mfa_verified == true
}

allow if {
	input.action == "delete"
	input.resource == "bio_identity_links"
	"bio.dcpj.director" in input.user.roles
	input.context.court_order_ref != ""
	input.context.mfa_verified == true
	input.context.judge_approval_ref != ""
}

# ── Wanted / Missing / Gang / Sex Offender ──────────────────────────────────

allow if {
	input.action in {"read", "create", "update"}
	input.resource in {"per_wanted_persons", "per_missing_persons", "per_gang_members", "per_sex_offenders"}
	"bio.dcpj.investigator" in input.user.roles
}

allow if {
	input.action == "read"
	input.resource in {"per_wanted_persons", "per_missing_persons"}
	"bio.ndis.analyst" in input.user.roles
}

# ── NDIS Matching ───────────────────────────────────────────────────────────

allow if {
	input.action == "match"
	input.resource == "ndis_cross_dept"
	"bio.ndis.analyst" in input.user.roles
	input.context.case_number != ""
}

# ── Audit Log ───────────────────────────────────────────────────────────────

allow if {
	input.action == "read"
	input.resource == "bio_audit_log"
	"bio.auditor" in input.user.roles
}

deny if {
	input.action == "delete"
	input.resource == "bio_audit_log"
}

# ── Admin operations (no data access) ───────────────────────────────────────

allow if {
	input.action in {"configure", "monitor"}
	input.resource in {"system", "metrics", "topics"}
	"bio.admin" in input.user.roles
}

# ── Default deny for unmatched rules ────────────────────────────────────────

deny if {
	not allow
}
