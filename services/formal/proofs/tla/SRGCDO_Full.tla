--------------------------- MODULE SRGCDO_Full ---------------------------

EXTENDS Naturals, Sequences, FiniteSets

VARIABLES
    nodes,
    events,
    risk,
    policies

CONSTANTS
    THRESHOLD,
    MAX_NODES

=============================================================================

(* Safety: No data leakage from restricted nodes *)
Safety ==
    \A e \in events :
        (policies[source(e)] = "RESTRICT") => (not exported(e))

(* Liveness: System eventually converges to a safe state *)
Liveness ==
    <> (\A n \in nodes : risk[n] <= THRESHOLD)

(* Decentralization: No node controls all state transitions *)
NoCentralControl ==
    \A n \in nodes :
        ~ ( \A m \in nodes : \E act \in Actions : Controls(n, m, act) )

(* State Invariant *)
TypeOK ==
    /\ nodes \subseteq 1..MAX_NODES
    /\ risk \in [nodes -> Nat]
    /\ policies \in [nodes -> {"ALLOW", "RESTRICT"}]

(* System Actions *)
Next ==
    \E n \in nodes :
        \/ NewEvent(n)
        \/ UpdatePolicy(n)
        \/ \E m \in nodes \ {n} : Propagate(n, m)

-----------------------------------------------------------------------------
=============================================================================
