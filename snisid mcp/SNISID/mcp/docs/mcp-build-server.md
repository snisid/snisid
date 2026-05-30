# MCP Build Server / SDK officiel

Ce projet suit le modèle officiel MCP : créer un `McpServer`, enregistrer tools/resources/prompts, choisir un transport STDIO ou Streamable HTTP, puis connecter le transport.

Scripts :

- `npm run mcp:stdio` : serveur local STDIO pour Claude/Cursor/VSCode.
- `npm run mcp:http` : serveur Streamable HTTP derrière Gateway.
- `npm run mcp:build` : validation typecheck + tests.
- `npm run mcp:inspect` : inspection MCP locale.

Les tools sont enregistrés centralement dans `mcp/server/registry.ts` pour éviter le tool poisoning et garantir un inventaire auditable.
