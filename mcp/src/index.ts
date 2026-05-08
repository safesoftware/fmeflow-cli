#!/usr/bin/env node
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from '@modelcontextprotocol/sdk/types.js';
import { fileURLToPath } from 'url';
import { buildArgs, runFmeflow, validateArgs, TOOL_DEFINITIONS } from './tools.js';

const server = new Server(
  { name: 'fmeflow-mcp', version: '0.1.0' },
  { capabilities: { tools: {} } },
);

server.setRequestHandler(ListToolsRequestSchema, async () => ({
  tools: TOOL_DEFINITIONS,
}));

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args = {} } = request.params;

  const validationError = validateArgs(name, args);
  if (validationError) {
    return { content: [{ type: 'text', text: validationError }], isError: true };
  }

  const cliArgs = buildArgs(name, args);
  if (cliArgs === null) {
    return {
      content: [{ type: 'text', text: `Unknown tool: ${name}` }],
      isError: true,
    };
  }

  const { content, isError } = runFmeflow(cliArgs);
  return { content: [{ type: 'text', text: content }], isError };
});

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}
