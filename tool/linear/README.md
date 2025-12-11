# Linear Tool Integration

This package provides a toolset for integrating with the [Linear](https://linear.app/) API. It allows agents to perform common operations such as creating issues, updating status, listing issues, and managing projects.

## Configuration

To use the Linear toolset, you need a Linear API Key. You can generate one in your Linear workspace settings under "API".

Initialize the toolset with the API key:

```go
import (
    "google.golang.org/adk/tool/linear"
)

cfg := linear.Config{
    APIKey: "your-linear-api-key",
}
linearTools, err := linear.New(cfg)
if err != nil {
    // handle error
}
```

## Available Tools

The toolset exposes the following tools to the agent:

- `linear_create_issue`: Creates a new issue.
- `linear_update_issue`: Updates an existing issue (title, description, priority, status, assignee).
- `linear_get_issue`: Retrieves details of an issue by ID or identifier (e.g., LIN-123).
- `linear_list_issues`: Lists issues, optionally filtered by team.
- `linear_get_projects`: Lists projects.
- `linear_create_project`: Creates a new project.
- `linear_add_comment`: Adds a comment to an issue.

## Usage Example

Add the toolset to your agent configuration:

```go
agentCfg := llmagent.Config{
    // ... other config ...
    Toolsets: []tool.Toolset{
        linearTools,
    },
}
```

The agent can then use these tools to interact with Linear. For example:
"Create a high priority issue in team T-123 titled 'Fix critical bug' assigned to user U-456."
