# tools

**Parent**: [examples](../AGENTS.md) | **Repository**: An open-source, code-first Go

> Contains example applications demonstrating the use of an Agent Development Kit (ADK) with Gemini Large Language Models. Implements both a command-line application for loading artifacts and a multi-tool AI application orchestrating various agents.

## Key Components

*No components detected*

## Folder Overview

## See Also

- [Parent overview →](../AGENTS.md) - Repository-level concepts and architecture
- [a2a/ →](../a2a/AGENTS.md) - Provides the main entry point for an Agent-to-Agent (A2A) application, orchestrating the setup and execution of an LLM-powered agent. Initializes all necessary components and launches an HTTP server to enable remote interaction with the agent.
- [mcp/ →](../mcp/AGENTS.md) - Orchestrates the initialization and execution of an AI agent application, configuring a Gemini model, dynamically selecting a Multi-modal Conversational Platform (MCP) transport, and integrating a specific toolset. Defines input/output data structures for agent interactions and provides mock tool implementations.
- [quickstart/ →](../quickstart/AGENTS.md) - Provides the primary entry point for an application that initializes and executes an LLM agent. Orchestrates the configuration of the agent with a Google Gemini model and integrates a Google Search tool, managing the application's startup and command-line processing.
