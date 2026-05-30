# Inter-Agency Connectors Inventory

## Status of Integration Connectors

| Agency | Connector Name | Type | Status | Features |
| :--- | :--- | :--- | :--- | :--- |
| **ONI** | `oni-id-bridge` | REST/gRPC | ✅ Active | Identity verification, Biometric match |
| **DGI** | `dgi-tax-bridge` | REST | ✅ Active | Tax status, NIF validation |
| **DCPJ** | `dcpj-police-check` | gRPC | ⚠️ Pending | Criminal record check |
| **MJSP** | `mjsp-legal-sync` | Events | ✅ Active | Court order notifications |
| **Immigration**| `die-passport-api` | REST | ✅ Active | Travel document validation |
| **MSPP** | `mspp-health-link` | FHIR/REST | ✅ Active | Vaccination records, Patient ID |

## Technical Implementation
- Connectors are deployed as sidecars or standalone micro-gateways.
- They handle protocol translation if legacy systems are involved.
- Mandatory audit logs for every data exchange.
