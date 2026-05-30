package snisid.rbac

default allow = false

# Allow if the user has the 'admin' role
allow {
    "admin" == input.user.roles[_]
}

# Allow if there's an explicit grant for the role
allow {
    some i
    role := input.user.roles[i]
    grant := data.role_grants[role][_]
    grant.action == input.action
    grant.resource == input.resource
}
