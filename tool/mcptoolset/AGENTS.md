# mcptoolset

**Parent**: [tool](../AGENTS.md) | **Repository**: An open-source Go toolkit for building, evaluating, and

> Provides a comprehensive framework for defining, configuring, and managing sets of tools specifically designed for Large Language Model interactions. Implements integration with a Managed Control Plane (MCP) to discover, abstract, and execute these tools, ensuring consistent operation and lifecycle management.

## Key Components

*No components detected*

## Folder Overview

## See Also

- [Parent overview →](../AGENTS.md) - Repository-level concepts and architecture
- [agenttool/ →](../agenttool/AGENTS.md) - Implements an `agentTool` wrapper, enabling the composition of agents by allowing one agent to invoke another as a tool within larger language model workflows. Defines a `Config` structure for tool-specific options and provides a constructor for creating `agentTool` instances.
- [exitlooptool/ →](../exitlooptool/AGENTS.md) - Provides a specialized tool designed to signal an agent's immediate exit from a processing loop. Implements a factory function for tool creation and includes comprehensive tests to ensure correct termination behavior.
- [functiontool/ →](../functiontool/AGENTS.md) - Implements a comprehensive framework for integrating Go functions as callable tools with generative AI models, ensuring type-safe execution and robust schema management. Provides mechanisms for defining, registering, and testing these tools, including support for asynchronous operations.
