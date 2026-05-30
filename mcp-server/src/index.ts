import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";

const server = new Server(
  {
    name: "snisid-context-server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

const SERVICES = [
  { name: "identity-api", description: "Core identity management and authentication" },
  { name: "frontend", description: "User interface for the SNISID system" },
  { name: "kafka", description: "Message broker for service communication" },
  { name: "mysql", description: "Primary relational database" },
];

server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "list_services",
        description: "List all microservices in the SNISID system",
        inputSchema: {
          type: "object",
          properties: {},
        },
      },
      {
        name: "get_system_architecture",
        description: "Get a high-level overview of the SNISID architecture",
        inputSchema: {
          type: "object",
          properties: {},
        },
      },
    ],
  };
});

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  switch (request.params.name) {
    case "list_services":
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify(SERVICES, null, 2),
          },
        ],
      };
    case "get_system_architecture":
      return {
        content: [
          {
            type: "text",
            text: "SNISID is a microservices-based national identity system. It uses Kafka for event-driven communication, MySQL for persistence, and a Go-based API layer with a React frontend.",
          },
        ],
      };
    default:
      throw new Error("Unknown tool");
  }
});

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("SNISID MCP Server running on stdio");
}

main().catch((error) => {
  console.error("Fatal error in main():", error);
  process.exit(1);
});
