package snisid.abac

default allow = false

# Allow if RBAC allows
allow {
    data.snisid.rbac.allow
}

# Allow if user belongs to the same agency as the resource they are trying to access
allow {
    input.action == "read"
    input.resource_type == "document"
    input.user.agency == input.resource_attributes.agency
}

# Deny if resource is classified and user does not have clearance
deny {
    input.resource_attributes.classification == "top_secret"
    input.user.clearance != "top_secret"
}
