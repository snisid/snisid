# SNISID Offline-First National Architecture

## Overview
The SNISID (Système National d'Identification et de Services d'Identité Digitale) must operate as a sovereign, resilient national infrastructure capable of surviving without internet, fiber connectivity, or during major national crises. This architecture ensures operational continuity across Haiti.

## National Resilience Objectives
- Function during internet outages
- Function during power failures
- Function during natural disasters or political instability
- Maintain data integrity and security offline
- Enable delayed synchronization upon reconnection

## Architecture Layers

### 1. National Core (Coordination Layer)
- Central coordination only when connectivity available
- No real-time dependency for regional operations
- Conflict resolution hub

### 2. Regional Edge Nodes (Departmental Autonomy)
- One per Haitian department
- Full local runtime, cache, IAM, workflows
- Temporary autonomous operation (days to weeks)

### 3. Mobile Nodes (Field Operations)
- Rugged mobile enrollment units
- Solar-powered with satellite backup
- Biometric enrollment and local verification

### 4. Offline Stations (Isolated Zones)
- Static offline-capable stations in remote areas
- Local data storage and processing

### 5. Emergency Kits (Crisis Response)
- Portable disaster recovery kits
- Pre-configured for rapid deployment

## Key Principles
- **No real-time core dependency**: All critical functions work offline
- **Event-driven with buffering**: All actions logged as events for replay
- **Prioritized sync**: Critical data (identity, judicial) sync first
- **Human-in-the-loop for conflicts**: Automated resolution where safe, human review for critical cases
- **Energy independence**: Solar + battery + generator at all layers
- **Security parity offline**: Same cryptographic standards apply offline

## Synchronization Model
- Delayed, queued synchronization
- Conflict detection and resolution engine
- Secure transport (when available)
- Partial sync support

This architecture guarantees national continuity under all conditions.