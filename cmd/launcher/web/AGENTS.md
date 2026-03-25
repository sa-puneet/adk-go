# web

**Parent**: [launcher](../AGENTS.md) | **Repository**: An open-source, code-first Go toolkit

> Provides the foundational infrastructure for launching and managing diverse web-facing components, including an Agent-to-Agent interface, a REST API server, and an embedded web user interface. Orchestrates the configuration, routing, and lifecycle of these sub-applications within a unified web server.

## Key Components

*No components detected*

## Folder Overview

## See Also

- [Parent overview →](../AGENTS.md) - Repository-level concepts and architecture
- [console/ →](../console/AGENTS.md) - Implements an interactive console interface for AI agent sessions, orchestrating user input, displaying agent responses, and managing session configuration and lifecycle. Provides a command-line entry point for continuous interaction with an AI agent.
- [full/ →](../full/AGENTS.md) - Implements a comprehensive launcher by aggregating various specialized launcher components. Orchestrates the creation of a single universal launcher from console, web, API, A2A, and web UI functionalities.
- [prod/ →](../prod/AGENTS.md) - Implements a factory function responsible for constructing a fully configured production `Launcher` instance. Provides a comprehensive system by orchestrating the integration of specialized API, A2A, web, and universal launcher components.
