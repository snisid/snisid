package snisid.zerotrust_test
import future.keywords.if
 
test_urgence_pnh_toujours_autorisee if {
    data.snisid.zerotrust.allow with input as {
        "type": "api_request",
        "user": {"role": "PNH_EMERGENCY", "clearance": "SECRET", "clearance_expired": false},
        "resource": "citizen_biometric",
        "action": "read",
        "context": {"mfa": true, "emergency": true}
    }
}
 
test_agent_oni_acces_lecture_autorise if {
    data.snisid.zerotrust.allow with input as {
        "type": "api_request",
        "user": {"role": "ONI_AGENT", "clearance": "CONFIDENTIEL", "clearance_expired": false},
        "resource": "citizen_identity_read",
        "action": "read",
        "context": {"mfa": true, "emergency": false}
    } with data.snisid.roles_agents_terrain as ["ONI_AGENT", "PNH_FIELD", "DIE_AGENT"]
}
 
test_dev_sans_clearance_production_refuse if {
    not data.snisid.zerotrust.allow with input as {
        "type": "api_request",
        "user": {"role": "DEVELOPER", "clearance": "NONE", "clearance_expired": false},
        "resource": "production_database",
        "action": "write",
        "context": {"mfa": true}
    }
}
