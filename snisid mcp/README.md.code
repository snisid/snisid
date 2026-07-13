# SNISID MCP Ecosystem — Infrastructure IA souveraine d’Haïti 🇭🇹

Ce dépôt est un socle **TypeScript Enterprise / MCP / Zero Trust** pour le Système National d’Identification et de Sécurité Intelligente (SNISID).

> Positionnement : ce code est un **starter gouvernement-grade** sécurisé. Les intégrations vers ONI, DGI, PNH, MJSP, Immigration, ANH, Passeport, Biométrie, Éducation et Intelligence passent exclusivement par des API gouvernementales contrôlées. Les tools MCP n’accèdent jamais directement aux bases de données.

## Chaîne cible

```text
AI Agents
  ↓
SNISID MCP
  ↓
API Gateway
  ↓
Government APIs
  ↓
Microservices
  ↓
Databases
```

## Dossiers racine

| Dossier | Rôle |
|---|---|
| `core/` | Principes de domaine, contrats transverses, politiques nationales. |
| `services/` | Documentation et connecteurs métiers hors MCP. |
| `database/` | Migrations/contrats DB des microservices, pas utilisés directement par MCP. |
| `gateway/` | API Gateway, politiques d’entrée, OpenAPI, routage, mTLS. |
| `security/` | Politiques institutionnelles et runbooks sécurité globaux. |
| `ai/` | Abstraction providers IA, orchestration, routage, failover. |
| `mcp/` | Serveur MCP souverain, tools, sécurité, audit, prompts, resources, tests. |

## Démarrage

```bash
cd SNISID
cp .env.example .env
npm install
npm run typecheck
npm run test
npm run mcp:stdio
# ou
npm run mcp:http
```

## Garde-fous obligatoires

- JWT, MFA, RBAC, Device Trust et Zero Trust sur chaque appel sensible.
- Audit immuable chiffré/rédigé compatible SIEM.
- Aucun secret hardcodé.
- Aucun accès DB direct depuis MCP.
- Human-in-the-loop pour décisions à impact légal ou sécuritaire.
