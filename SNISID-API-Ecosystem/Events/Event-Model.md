# SNISID Event-Driven Interoperability Model

## Architecture
The National Event Bus (Apache Kafka) facilitates real-time, decoupled communication between agencies.

## Core Topics & Events
| Topic | Event Type | Description | Producer | Consumers |
| :--- | :--- | :--- | :--- | :--- |
| `national.identity.events` | `citizen.created` | New citizen registration | ONI | DGI, MSPP, DIE |
| `national.identity.events` | `identity.updated` | Update of personal info | ONI | All |
| `national.justice.events` | `judicial.case.created` | New judicial case | MJSP | DCPJ, ONI |
| `national.civil.events` | `birth.registered` | New birth registration | État Civil | ONI, MSPP |
| `national.border.events` | `passport.issued` | Passport delivery | Immigration | ONI, DCPJ |

## Event Payload Structure (JSON)
```json
{
  "event_id": "uuid-v4",
  "event_type": "citizen.created",
  "timestamp": "2026-05-25T14:30:00Z",
  "source": "ONI-API",
  "version": "1.0",
  "data": {
    "niu": "1234-5678-9012",
    "last_name": "Jean-Pierre",
    "first_name": "Marie"
  },
  "metadata": {
    "trace_id": "correlation-id-abc",
    "security_context": "level-high"
  }
}
```

## Resilience
- **Persistence**: Events are stored for 30 days.
- **Replay**: Agencies can replay events in case of failure.
