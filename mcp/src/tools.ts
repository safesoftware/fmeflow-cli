import { spawnSync } from 'child_process';
import { existsSync } from 'fs';
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

// ─── Binary resolution (memoized) ─────────────────────────────────────────────

let cachedBinary: string | null | undefined = undefined;

function resolveBinary(): string | null {
  if (cachedBinary !== undefined) return cachedBinary;
  const check = spawnSync('fmeflow', ['--version'], { encoding: 'utf8', stdio: 'pipe' });
  if (!check.error && check.status === 0) return (cachedBinary = 'fmeflow');
  const repoBinary = join(__dirname, '..', '..', 'fmeflow');
  if (existsSync(repoBinary)) return (cachedBinary = repoBinary);
  return (cachedBinary = null);
}

// ─── Required-arg validation ───────────────────────────────────────────────────

function required(args: Record<string, unknown>, ...keys: string[]): string | null {
  for (const key of keys) {
    if (args[key] === undefined || args[key] === null || args[key] === '') {
      return `Missing required argument: ${key}`;
    }
  }
  return null;
}

// ─── Tool registry ─────────────────────────────────────────────────────────────

type ToolEntry = {
  definition: { name: string; description: string; inputSchema: object };
  requiredArgs: string[];
  buildArgs: (args: Record<string, unknown>) => string[];
};

const TOOL_REGISTRY: Record<string, ToolEntry> = {
  fmeflow_healthcheck: {
    definition: {
      name: 'fmeflow_healthcheck',
      description: 'Check the health status of an FME Flow instance. Does not require authentication.',
      inputSchema: {
        type: 'object',
        properties: {
          url: {
            type: 'string',
            description: 'Base URL of the FME Flow instance (e.g. https://my-fmeflow.internal). Optional if already logged in.',
          },
          ready: {
            type: 'boolean',
            description: 'If true, checks whether FME Flow is ready to process jobs, not just alive.',
          },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['healthcheck'];
      if (args.url) flags.push('--url', String(args.url));
      if (args.ready) flags.push('--ready');
      return flags;
    },
  },

  fmeflow_info: {
    definition: {
      name: 'fmeflow_info',
      description: 'Get version and build information about the connected FME Flow instance.',
      inputSchema: { type: 'object', properties: {} },
    },
    requiredArgs: [],
    buildArgs: () => ['info'],
  },

  fmeflow_repositories_list: {
    definition: {
      name: 'fmeflow_repositories_list',
      description: 'List repositories on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Get details for a single named repository.',
          },
          owner: {
            type: 'string',
            description: 'Filter repositories by owner username.',
          },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['repositories'];
      if (args.name) flags.push('--name', String(args.name));
      if (args.owner) flags.push('--owner', String(args.owner));
      return flags;
    },
  },

  fmeflow_workspaces_list: {
    definition: {
      name: 'fmeflow_workspaces_list',
      description: 'List workspaces on FME Flow. Optionally filter by repository or retrieve details for a specific workspace.',
      inputSchema: {
        type: 'object',
        properties: {
          repository: {
            type: 'string',
            description: 'Repository name to filter workspaces.',
          },
          name: {
            type: 'string',
            description: 'Specific workspace filename to retrieve details for. Requires repository.',
          },
          filterString: {
            type: 'string',
            description: 'Filter workspaces by matching name or title substring.',
          },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['workspaces'];
      if (args.repository) flags.push('--repository', String(args.repository));
      if (args.name) flags.push('--name', String(args.name));
      if (args.filterString) flags.push('--filter-string', String(args.filterString));
      return flags;
    },
  },

  fmeflow_run: {
    definition: {
      name: 'fmeflow_run',
      description: 'Submit a workspace job on FME Flow. Use wait=true to block until the job completes and get the result immediately.',
      inputSchema: {
        type: 'object',
        properties: {
          repository: {
            type: 'string',
            description: 'Repository containing the workspace.',
          },
          workspace: {
            type: 'string',
            description: 'Workspace filename (e.g. myWorkspace.fmw).',
          },
          wait: {
            type: 'boolean',
            description: 'Wait for the job to finish and return the result. If false, returns only the job ID.',
          },
          publishedParameters: {
            type: 'object',
            description: 'Key-value pairs of published parameters to pass to the workspace.',
            additionalProperties: { type: 'string' },
          },
          queue: {
            type: 'string',
            description: 'Queue to submit the job to.',
          },
        },
        required: ['repository', 'workspace'],
      },
    },
    requiredArgs: ['repository', 'workspace'],
    buildArgs: (args) => {
      const flags = [
        'run',
        '--repository', String(args.repository),
        '--workspace', String(args.workspace),
      ];
      if (args.wait) flags.push('--wait');
      if (args.queue) flags.push('--queue', String(args.queue));
      if (args.publishedParameters && typeof args.publishedParameters === 'object') {
        for (const [k, v] of Object.entries(args.publishedParameters as Record<string, string>)) {
          flags.push('--published-parameter', `${k}=${v}`);
        }
      }
      return flags;
    },
  },

  fmeflow_jobs: {
    definition: {
      name: 'fmeflow_jobs',
      description: 'List jobs or get the status of a specific job on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          id: {
            type: 'number',
            description: 'Specific job ID to retrieve. Cannot be combined with status/repository/workspace filters.',
          },
          repository: {
            type: 'string',
            description: 'Filter jobs by repository.',
          },
          workspace: {
            type: 'string',
            description: 'Filter jobs by workspace. Requires repository.',
          },
          status: {
            type: 'string',
            enum: ['all', 'active', 'running', 'queued', 'completed', 'success', 'failure', 'cancelled'],
            description: 'Filter jobs by status. Defaults to all.',
          },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['jobs'];
      if (args.id !== undefined) {
        flags.push('--id', String(args.id));
      } else {
        const statusFlag: Record<string, string> = {
          all: '--all', active: '--active', running: '--running', queued: '--queued',
          completed: '--completed', success: '--success', failure: '--failure', cancelled: '--cancelled',
        };
        flags.push(statusFlag[String(args.status ?? 'all')] ?? '--all');
        if (args.repository) flags.push('--repository', String(args.repository));
        if (args.workspace) flags.push('--workspace', String(args.workspace));
      }
      return flags;
    },
  },

  fmeflow_job_cancel: {
    definition: {
      name: 'fmeflow_job_cancel',
      description: 'Cancel a running or queued job on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          id: {
            type: 'number',
            description: 'The job ID to cancel.',
          },
        },
        required: ['id'],
      },
    },
    requiredArgs: ['id'],
    buildArgs: (args) => ['cancel', '--id', String(args.id)],
  },

  fmeflow_connections_list: {
    definition: {
      name: 'fmeflow_connections_list',
      description: 'List connections defined on FME Flow (database, web, cloud storage, etc.).',
      inputSchema: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Get a single connection by name.',
          },
          type: {
            type: 'string',
            description: 'Filter by connection type (e.g. "Google Drive").',
          },
          category: {
            type: 'string',
            description: 'Filter by category: database, web, or oauthConnection.',
          },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['connections'];
      if (args.name) flags.push('--name', String(args.name));
      if (args.type) flags.push('--type', String(args.type));
      if (args.category) flags.push('--category', String(args.category));
      return flags;
    },
  },

  fmeflow_deployment_parameters_list: {
    definition: {
      name: 'fmeflow_deployment_parameters_list',
      description: 'List deployment parameters defined on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Get a single deployment parameter by name.',
          },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['deploymentparameters'];
      if (args.name) flags.push('--name', String(args.name));
      return flags;
    },
  },

  fmeflow_engines: {
    definition: {
      name: 'fmeflow_engines',
      description: 'List FME engines registered with FME Flow.',
      inputSchema: { type: 'object', properties: {} },
    },
    requiredArgs: [],
    buildArgs: () => ['engines'],
  },

  fmeflow_backup: {
    definition: {
      name: 'fmeflow_backup',
      description: 'Back up the FME Flow configuration to a local file or to a shared resource on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          file: { type: 'string', description: 'Local path to download the backup file to. Defaults to ServerConfigPackage.fsconfig.' },
          resource: { type: 'boolean', description: 'If true, back up to a shared resource on FME Flow instead of downloading locally.' },
          resourceName: { type: 'string', description: 'Shared resource name to save the backup to. Defaults to FME_SHAREDRESOURCE_BACKUP.' },
          exportPackage: { type: 'string', description: 'Path and name of the export package when backing up to a shared resource.' },
          failureTopic: { type: 'string', description: 'Topic to notify on backup failure.' },
          successTopic: { type: 'string', description: 'Topic to notify on backup success.' },
          suppressFileRename: { type: 'boolean', description: 'If true, do not automatically append .fsconfig to the file name.' },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['backup'];
      if (args.file) flags.push('--file', String(args.file));
      if (args.resource) flags.push('--resource');
      if (args.resourceName) flags.push('--resource-name', String(args.resourceName));
      if (args.exportPackage) flags.push('--export-package', String(args.exportPackage));
      if (args.failureTopic) flags.push('--failure-topic', String(args.failureTopic));
      if (args.successTopic) flags.push('--success-topic', String(args.successTopic));
      if (args.suppressFileRename) flags.push('--suppress-file-rename');
      return flags;
    },
  },

  fmeflow_restore: {
    definition: {
      name: 'fmeflow_restore',
      description: 'Restore the FME Flow configuration from a backup package.',
      inputSchema: {
        type: 'object',
        properties: {
          file: { type: 'string', description: 'Path to the backup file. Local path or path within a shared resource when --resource is set.' },
          resource: { type: 'boolean', description: 'If true, restore from a shared resource instead of a local file.' },
          resourceName: { type: 'string', description: 'Shared resource containing the backup package. Defaults to FME_SHAREDRESOURCE_BACKUP.' },
          importMode: { type: 'string', enum: ['INSERT', 'UPDATE'], description: 'How to handle existing items. Defaults to INSERT.' },
          projectsImportMode: { type: 'string', enum: ['INSERT', 'UPDATE'], description: 'Import mode specifically for projects.' },
          pauseNotifications: { type: 'boolean', description: 'Pause notifications during restore. Defaults to true.' },
          overwrite: { type: 'boolean', description: 'Overwrite existing items.' },
          failureTopic: { type: 'string', description: 'Topic to notify on restore failure.' },
          successTopic: { type: 'string', description: 'Topic to notify on restore success.' },
        },
        required: ['file'],
      },
    },
    requiredArgs: ['file'],
    buildArgs: (args) => {
      const flags = ['restore', '--file', String(args.file)];
      if (args.resource) flags.push('--resource');
      if (args.resourceName) flags.push('--resource-name', String(args.resourceName));
      if (args.importMode) flags.push('--import-mode', String(args.importMode));
      if (args.projectsImportMode) flags.push('--projects-import-mode', String(args.projectsImportMode));
      if (args.pauseNotifications === false) flags.push('--pause-notifications=false');
      if (args.overwrite) flags.push('--overwrite');
      if (args.failureTopic) flags.push('--failure-topic', String(args.failureTopic));
      if (args.successTopic) flags.push('--success-topic', String(args.successTopic));
      return flags;
    },
  },

  fmeflow_repositories_create: {
    definition: {
      name: 'fmeflow_repositories_create',
      description: 'Create a new repository on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the repository to create.' },
          description: { type: 'string', description: 'Description for the new repository.' },
        },
        required: ['name'],
      },
    },
    requiredArgs: ['name'],
    buildArgs: (args) => {
      const flags = ['repositories', 'create', '--name', String(args.name)];
      if (args.description) flags.push('--description', String(args.description));
      return flags;
    },
  },

  fmeflow_repositories_delete: {
    definition: {
      name: 'fmeflow_repositories_delete',
      description: 'Delete a repository from FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the repository to delete.' },
        },
        required: ['name'],
      },
    },
    requiredArgs: ['name'],
    buildArgs: (args) => ['repositories', 'delete', '--name', String(args.name), '--no-prompt'],
  },

  fmeflow_connections_create: {
    definition: {
      name: 'fmeflow_connections_create',
      description: 'Create a new connection on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the connection to create.' },
          category: { type: 'string', enum: ['basic', 'database', 'token', 'oauthV1', 'oauthV2'], description: 'Category of the connection.' },
          type: { type: 'string', description: 'Connection type (e.g. "PostgreSQL", "Google Drive").' },
          authenticationMethod: { type: 'string', description: 'Authentication method for the connection.' },
          username: { type: 'string', description: 'Username for the connection.' },
          password: { type: 'string', description: 'Password for the connection.' },
          parameters: {
            type: 'object',
            description: 'Additional connection parameters as key-value pairs.',
            additionalProperties: { type: 'string' },
          },
        },
        required: ['name', 'category'],
      },
    },
    requiredArgs: ['name', 'category'],
    buildArgs: (args) => {
      const flags = ['connections', 'create', '--name', String(args.name), '--category', String(args.category)];
      if (args.type) flags.push('--type', String(args.type));
      if (args.authenticationMethod) flags.push('--authentication-method', String(args.authenticationMethod));
      if (args.username) flags.push('--username', String(args.username));
      if (args.password) flags.push('--password', String(args.password));
      if (args.parameters && typeof args.parameters === 'object') {
        for (const [k, v] of Object.entries(args.parameters as Record<string, string>)) {
          flags.push('--parameter', `${k}=${v}`);
        }
      }
      return flags;
    },
  },

  fmeflow_connections_delete: {
    definition: {
      name: 'fmeflow_connections_delete',
      description: 'Delete a connection from FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the connection to delete.' },
        },
        required: ['name'],
      },
    },
    requiredArgs: ['name'],
    buildArgs: (args) => ['connections', 'delete', '--name', String(args.name), '--no-prompt'],
  },

  fmeflow_connections_update: {
    definition: {
      name: 'fmeflow_connections_update',
      description: 'Update an existing connection on FME Flow. Only fields that need to change must be specified.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the connection to update.' },
          authenticationMethod: { type: 'string', description: 'New authentication method.' },
          username: { type: 'string', description: 'New username.' },
          password: { type: 'string', description: 'New password.' },
          parameters: {
            type: 'object',
            description: 'Connection parameters to update as key-value pairs.',
            additionalProperties: { type: 'string' },
          },
        },
        required: ['name'],
      },
    },
    requiredArgs: ['name'],
    buildArgs: (args) => {
      const flags = ['connections', 'update', '--name', String(args.name)];
      if (args.authenticationMethod) flags.push('--authentication-method', String(args.authenticationMethod));
      if (args.username) flags.push('--username', String(args.username));
      if (args.password) flags.push('--password', String(args.password));
      if (args.parameters && typeof args.parameters === 'object') {
        for (const [k, v] of Object.entries(args.parameters as Record<string, string>)) {
          flags.push('--parameter', `${k}=${v}`);
        }
      }
      return flags;
    },
  },

  fmeflow_deployment_parameters_create: {
    definition: {
      name: 'fmeflow_deployment_parameters_create',
      description: 'Create a new deployment parameter on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the deployment parameter.' },
          type: { type: 'string', enum: ['text', 'database', 'web'], description: 'Type of the deployment parameter. Defaults to text.' },
          value: { type: 'string', description: 'Value of the deployment parameter.' },
          databaseType: { type: 'string', description: 'Database type (only for type=database).' },
          includedServices: { type: 'array', items: { type: 'string' }, description: 'Services to include (only for type=web).' },
          excludedServices: { type: 'array', items: { type: 'string' }, description: 'Services to exclude (only for type=web).' },
        },
        required: ['name'],
      },
    },
    requiredArgs: ['name'],
    buildArgs: (args) => {
      const flags = ['deploymentparameters', 'create', '--name', String(args.name)];
      if (args.type) flags.push('--type', String(args.type));
      if (args.value !== undefined) flags.push('--value', String(args.value));
      if (args.databaseType) flags.push('--database-type', String(args.databaseType));
      if (Array.isArray(args.includedServices)) {
        for (const s of args.includedServices) flags.push('--included-service', String(s));
      }
      if (Array.isArray(args.excludedServices)) {
        for (const s of args.excludedServices) flags.push('--excluded-service', String(s));
      }
      return flags;
    },
  },

  fmeflow_deployment_parameters_delete: {
    definition: {
      name: 'fmeflow_deployment_parameters_delete',
      description: 'Delete a deployment parameter from FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the deployment parameter to delete.' },
        },
        required: ['name'],
      },
    },
    requiredArgs: ['name'],
    buildArgs: (args) => ['deploymentparameters', 'delete', '--name', String(args.name), '--no-prompt'],
  },

  fmeflow_deployment_parameters_update: {
    definition: {
      name: 'fmeflow_deployment_parameters_update',
      description: 'Update an existing deployment parameter on FME Flow.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the deployment parameter to update.' },
          value: { type: 'string', description: 'New value for the deployment parameter.' },
          type: { type: 'string', enum: ['text', 'database', 'web'], description: 'Type of the deployment parameter.' },
          databaseType: { type: 'string', description: 'Database type (only for type=database).' },
          includedServices: { type: 'array', items: { type: 'string' }, description: 'Services to include (only for type=web).' },
          excludedServices: { type: 'array', items: { type: 'string' }, description: 'Services to exclude (only for type=web).' },
        },
        required: ['name', 'value'],
      },
    },
    requiredArgs: ['name', 'value'],
    buildArgs: (args) => {
      const flags = ['deploymentparameters', 'update', '--name', String(args.name), '--value', String(args.value)];
      if (args.type) flags.push('--type', String(args.type));
      if (args.databaseType) flags.push('--database-type', String(args.databaseType));
      if (Array.isArray(args.includedServices)) {
        for (const s of args.includedServices) flags.push('--included-service', String(s));
      }
      if (Array.isArray(args.excludedServices)) {
        for (const s of args.excludedServices) flags.push('--excluded-service', String(s));
      }
      return flags;
    },
  },

  fmeflow_projects_list: {
    definition: {
      name: 'fmeflow_projects_list',
      description: 'List projects on FME Flow, or retrieve details for a single project by name or ID.',
      inputSchema: {
        type: 'object',
        properties: {
          owner: { type: 'string', description: 'Filter projects by owner username.' },
          name: { type: 'string', description: 'Return a single project with this name.' },
          id: { type: 'string', description: 'Return a single project with this ID (v4 only).' },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['projects'];
      if (args.owner) flags.push('--owner', String(args.owner));
      if (args.name) flags.push('--name', String(args.name));
      if (args.id) flags.push('--id', String(args.id));
      return flags;
    },
  },

  fmeflow_projects_upload: {
    definition: {
      name: 'fmeflow_projects_upload',
      description: 'Import an FME Flow project from a local package file.',
      inputSchema: {
        type: 'object',
        properties: {
          file: { type: 'string', description: 'Local path to the .fsproject package file to import.' },
          overwrite: { type: 'boolean', description: 'Overwrite existing items. Defaults to true.' },
          pauseNotifications: { type: 'boolean', description: 'Disable notifications during import. Defaults to true.' },
          disableProjectItems: { type: 'boolean', description: 'Disable all imported project items after import.' },
          selectedItems: { type: 'string', description: 'Items to import: "all", "none", or a comma-separated list. Defaults to all.' },
          wait: { type: 'boolean', description: 'Wait for the import to complete. Defaults to true.' },
          failureTopic: { type: 'string', description: 'Topic to notify on import failure.' },
          successTopic: { type: 'string', description: 'Topic to notify on import success.' },
        },
        required: ['file'],
      },
    },
    requiredArgs: ['file'],
    buildArgs: (args) => {
      const flags = ['projects', 'upload', '--file', String(args.file)];
      if (args.overwrite === false) flags.push('--overwrite=false');
      if (args.pauseNotifications === false) flags.push('--pause-notifications=false');
      if (args.disableProjectItems) flags.push('--disable-project-items');
      if (args.selectedItems) flags.push('--selected-items', String(args.selectedItems));
      if (args.wait === false) flags.push('--wait=false');
      if (args.failureTopic) flags.push('--failure-topic', String(args.failureTopic));
      if (args.successTopic) flags.push('--success-topic', String(args.successTopic));
      return flags;
    },
  },

  fmeflow_projects_download: {
    definition: {
      name: 'fmeflow_projects_download',
      description: 'Download an FME Flow project to a local package file.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the project to download.' },
          id: { type: 'string', description: 'ID of the project to download (v4 only).' },
          file: { type: 'string', description: 'Local path to save the downloaded package. Defaults to ProjectPackage.fsproject.' },
          excludeSensitiveInfo: { type: 'boolean', description: 'Exclude sensitive information from the package.' },
          suppressFileRename: { type: 'boolean', description: 'Do not automatically append .fsproject to the file name.' },
          excludeAllSelectableItems: { type: 'boolean', description: 'Exclude all selectable items from the package.' },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['projects', 'download'];
      if (args.name) flags.push('--name', String(args.name));
      if (args.id) flags.push('--id', String(args.id));
      if (args.file) flags.push('--file', String(args.file));
      if (args.excludeSensitiveInfo) flags.push('--exclude-sensitive-info');
      if (args.suppressFileRename) flags.push('--suppress-file-rename');
      if (args.excludeAllSelectableItems) flags.push('--exclude-all-selectable-items');
      return flags;
    },
  },

  fmeflow_projects_items: {
    definition: {
      name: 'fmeflow_projects_items',
      description: 'List items contained in a specific FME Flow project.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the project to list items for.' },
          id: { type: 'string', description: 'ID of the project to list items for.' },
          type: { type: 'array', items: { type: 'string' }, description: 'Filter by item type (e.g. "workspace", "connection"). Repeatable.' },
          includeDependencies: { type: 'boolean', description: 'Include dependencies in the output. Defaults to true.' },
          filterString: { type: 'string', description: 'Filter items by matching substring.' },
          filterProperty: { type: 'array', items: { type: 'string' }, description: 'Properties to filter by (e.g. "name", "owner").' },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['projects', 'items'];
      if (args.name) flags.push('--name', String(args.name));
      if (args.id) flags.push('--id', String(args.id));
      if (Array.isArray(args.type)) {
        for (const t of args.type) flags.push('--type', String(t));
      }
      if (args.includeDependencies === false) flags.push('--include-dependencies=false');
      if (args.filterString) flags.push('--filter-string', String(args.filterString));
      if (Array.isArray(args.filterProperty)) {
        for (const p of args.filterProperty) flags.push('--filter-property', String(p));
      }
      return flags;
    },
  },

  fmeflow_projects_delete: {
    definition: {
      name: 'fmeflow_projects_delete',
      description: 'Delete an FME Flow project.',
      inputSchema: {
        type: 'object',
        properties: {
          name: { type: 'string', description: 'Name of the project to delete.' },
          id: { type: 'string', description: 'ID of the project to delete.' },
          all: { type: 'boolean', description: 'Delete the project and all its contents.' },
          dependencies: { type: 'boolean', description: 'Delete the project, its contents, and all dependencies.' },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['projects', 'delete', '--no-prompt'];
      if (args.name) flags.push('--name', String(args.name));
      if (args.id) flags.push('--id', String(args.id));
      if (args.all) flags.push('--all');
      if (args.dependencies) flags.push('--dependencies');
      return flags;
    },
  },

  fmeflow_license_machinekey: {
    definition: {
      name: 'fmeflow_license_machinekey',
      description: 'Retrieve the machine key of the machine running FME Flow.',
      inputSchema: { type: 'object', properties: {} },
    },
    requiredArgs: [],
    buildArgs: () => ['license', 'machinekey'],
  },

  fmeflow_migration_tasks: {
    definition: {
      name: 'fmeflow_migration_tasks',
      description: 'List migration task records on FME Flow, or retrieve a specific task by ID.',
      inputSchema: {
        type: 'object',
        properties: {
          id: { type: 'number', description: 'ID of a specific migration task to retrieve.' },
        },
      },
    },
    requiredArgs: [],
    buildArgs: (args) => {
      const flags = ['migration', 'tasks'];
      if (args.id !== undefined) flags.push('--id', String(args.id));
      return flags;
    },
  },
};

// ─── Public exports ────────────────────────────────────────────────────────────

export const TOOL_DEFINITIONS = Object.values(TOOL_REGISTRY).map(e => e.definition);

export function buildArgs(name: string, args: Record<string, unknown>): string[] | null {
  return TOOL_REGISTRY[name]?.buildArgs(args) ?? null;
}

export function validateArgs(name: string, args: Record<string, unknown>): string | null {
  const entry = TOOL_REGISTRY[name];
  if (!entry) return null;
  return required(args, ...entry.requiredArgs);
}

export function runFmeflow(args: string[]): { content: string; isError: boolean } {
  const binary = resolveBinary();
  if (!binary) {
    return {
      content: 'fmeflow binary not found on PATH or at repo root. Run: go build -o fmeflow',
      isError: true,
    };
  }

  const result = spawnSync(binary, [...args, '--json'], { encoding: 'utf8' });
  if (result.error) {
    return { content: result.error.message, isError: true };
  }

  const raw = (result.stdout || result.stderr || '').trim();
  const isError = result.status !== 0;

  // Strip HTML tags from description fields to reduce noise for the model
  const content = raw.replace(/"description"\s*:\s*"((?:[^"\\]|\\.)*)"/g, (_, val) => {
    const stripped = val
      .replace(/\\n/g, ' ')
      .replace(/<[^>]+>/g, '')
      .replace(/&amp;/g, '&')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&quot;/g, '"')
      .replace(/\s+/g, ' ')
      .trim();
    return `"description": "${stripped.replace(/"/g, '\\"')}"`;
  });

  return { content, isError };
}
