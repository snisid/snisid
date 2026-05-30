package snisid.abac

default allow = false

# Allow if the user has the required role AND contextual conditions are met
allow {
    user_has_role("AGENT_ENROLEMENT")
    is_working_hours
    is_authorized_ip
}

allow {
    user_has_role("SUPERVISEUR_REGION")
    is_working_hours
    is_authorized_ip
}

# Admins and Auditors have 24/7 access but still require authorized networks/VPN
allow {
    user_has_role("ADMIN_SYSTEME")
    is_authorized_ip
}

allow {
    user_has_role("AUDITEUR_NATIONAL")
    is_authorized_ip
}

# --- Contextual Helpers ---

user_has_role(role_name) {
    input.jwt.realm_access.roles[_] == role_name
}

is_working_hours {
    # Extract current hour in local timezone (assume input provides local time context or calculate from epoch)
    time.clock([input.time_ns, "UTC"])[1][0] >= 8  # 08:00
    time.clock([input.time_ns, "UTC"])[1][0] <= 18 # 18:00
}

is_authorized_ip {
    # Check if the client IP is within the authorized regional offices or VPN subnets
    authorized_cidrs := [
        "10.200.0.0/16",  # Regional Offices
        "192.168.20.0/24" # Admin VPN
    ]
    net.cidr_contains(authorized_cidrs[_], input.client_ip)
}
