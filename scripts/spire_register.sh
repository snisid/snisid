# SPIRE Registration Entries for SNISID Workloads

spire-server entry create \
    -parentID spiffe://cluster.local/ns/spire/sa/spire-agent \
    -spiffeID spiffe://cluster.local/ns/snisid/sa/api-gateway-sa \
    -selector k8s:ns:snisid \
    -selector k8s:sa:api-gateway-sa

spire-server entry create \
    -parentID spiffe://cluster.local/ns/spire/sa/spire-agent \
    -spiffeID spiffe://cluster.local/ns/snisid/sa/identity-api-sa \
    -selector k8s:ns:snisid \
    -selector k8s:sa:identity-api-sa

spire-server entry create \
    -parentID spiffe://cluster.local/ns/spire/sa/spire-agent \
    -spiffeID spiffe://cluster.local/ns/snisid/sa/fraud-engine-sa \
    -selector k8s:ns:snisid \
    -selector k8s:sa:fraud-engine-sa
