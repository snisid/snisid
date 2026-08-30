package nexus.authz

import future.keywords.if
import future.keywords.in

default allow = false

# Allow administrators to perform any action on any resource
allow if {
    input.user.role == "admin"
}

# Allow Kai agents to stream signals and execute tasks
allow if {
    input.agent.id != ""
    input.action in ["StreamAgentSignals", "ExecuteTask"]
}

# Allow Vera (Strategic Engine) to submit tasks and check workflow status
allow if {
    input.user.role == "strategic_engine"
    input.action in ["SubmitTask", "GetWorkflowStatus"]
}

# Restrict task submission by type
allow if {
    input.user.role == "operator"
    input.action == "SubmitTask"
    input.task.type in ["RECON", "ANALYSIS"]
}
