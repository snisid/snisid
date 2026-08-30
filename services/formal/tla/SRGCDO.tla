------------------------------ MODULE SRGCDO ------------------------------

EXTENDS Naturals, Sequences

VARIABLES
    nodes,
    events,
    risk,
    policy

CONSTANTS
    THRESHOLD

=============================================================================

(* State Initialization *)
Init ==
    /\ nodes = {"HT", "CA", "US"}
    /\ events = << >>
    /\ risk = [n \in nodes |-> 0]
    /\ policy = [n \in nodes |-> "ALLOW"]

(* Action: Record New Security Event *)
NewEvent(e, n) ==
    /\ events' = Append(events, e)
    /\ risk' = [risk EXCEPT ![n] = risk[n] + e.risk]
    /\ UNCHANGED policy

(* Action: Update Local Security Policy based on Risk *)
UpdatePolicy(n) ==
    /\ IF risk[n] > THRESHOLD
       THEN policy' = [policy EXCEPT ![n] = "RESTRICT"]
       ELSE UNCHANGED policy
    /\ UNCHANGED <<events, risk>>

(* Action: Propagate Event across Federation nodes *)
Propagate(e, src, dst) ==
    /\ policy[src] = "ALLOW"
    /\ events' = Append(events, e)
    /\ UNCHANGED <<risk, policy>>

(* Next State Transition *)
Next ==
    \E n \in nodes :
        \/ \E e \in [risk: 1..10] : NewEvent(e, n)
        \/ UpdatePolicy(n)
        \/ \E m \in nodes \ {n} : \E e \in [risk: 1..10] : Propagate(e, n, m)

(* Safety Invariant: Risk must not exceed THRESHOLD without restriction *)
Invariant ==
    \A n \in nodes :
        risk[n] <= THRESHOLD => policy[n] = "ALLOW"

-----------------------------------------------------------------------------
=============================================================================
