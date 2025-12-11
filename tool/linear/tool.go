// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package linear

import (
    "fmt"

    "google.golang.org/adk/agent"
    "google.golang.org/adk/tool"
    "google.golang.org/adk/tool/functiontool"
)

// Config provides configuration for the Linear ToolSet.
type Config struct {
    APIKey   string
    Endpoint string
}

// New creates a new Linear ToolSet.
func New(cfg Config) (tool.Toolset, error) {
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("API key is required")
    }
    var opts []ClientOption
    if cfg.Endpoint != "" {
        opts = append(opts, WithEndpoint(cfg.Endpoint))
    }
    client := NewClient(cfg.APIKey, opts...)
    return &linearToolset{
        client: client,
    }, nil
}

type linearToolset struct {
    client *Client
}

func (s *linearToolset) Name() string {
    return "linear"
}

func (s *linearToolset) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
    var tools []tool.Tool

    // Create Issue Tool
    createIssueTool, err := functiontool.New(functiontool.Config{
        Name:        "linear_create_issue",
        Description: "Creates a new issue in Linear.",
    }, func(ctx tool.Context, args CreateIssueInput) (*Issue, error) {
        return s.client.CreateIssue(ctx, args)
    })
    if err != nil {
        return nil, err
    }
    tools = append(tools, createIssueTool)

    // Update Issue Tool
    updateIssueTool, err := functiontool.New(functiontool.Config{
        Name:        "linear_update_issue",
        Description: "Updates an existing issue in Linear.",
    }, func(ctx tool.Context, args UpdateIssueInput) (*Issue, error) {
        return s.client.UpdateIssue(ctx, args)
    })
    if err != nil {
        return nil, err
    }
    tools = append(tools, updateIssueTool)

    // Get Issue Tool
    getIssueTool, err := functiontool.New(functiontool.Config{
        Name:        "linear_get_issue",
        Description: "Retrieves an issue from Linear by ID or identifier.",
    }, func(ctx tool.Context, args GetIssueArgs) (*Issue, error) {
        return s.client.GetIssue(ctx, args.ID)
    })
    if err != nil {
        return nil, err
    }
    tools = append(tools, getIssueTool)

    // List Issues Tool
    listIssuesTool, err := functiontool.New(functiontool.Config{
        Name:        "linear_list_issues",
        Description: "Lists issues from Linear, optionally filtered by team.",
    }, func(ctx tool.Context, args ListIssuesArgs) ([]Issue, error) {
        limit := args.Limit
        if limit <= 0 {
            limit = 50
        }
        return s.client.ListIssues(ctx, limit, args.TeamID)
    })
    if err != nil {
        return nil, err
    }
    tools = append(tools, listIssuesTool)

    // Get Projects Tool
    getProjectsTool, err := functiontool.New(functiontool.Config{
        Name:        "linear_get_projects",
        Description: "Lists projects from Linear.",
    }, func(ctx tool.Context, args GetProjectsArgs) ([]Project, error) {
        limit := args.Limit
        if limit <= 0 {
            limit = 50
        }
        return s.client.GetProjects(ctx, limit)
    })
    if err != nil {
        return nil, err
    }
    tools = append(tools, getProjectsTool)

    // Create Project Tool
    createProjectTool, err := functiontool.New(functiontool.Config{
        Name:        "linear_create_project",
        Description: "Creates a new project in Linear.",
    }, func(ctx tool.Context, args CreateProjectArgs) (*Project, error) {
        return s.client.CreateProject(ctx, args.Name, args.TeamIDs, args.Description)
    })
    if err != nil {
        return nil, err
    }
    tools = append(tools, createProjectTool)

    // Add Comment Tool
    addCommentTool, err := functiontool.New(functiontool.Config{
        Name:        "linear_add_comment",
        Description: "Adds a comment to a Linear issue.",
    }, func(ctx tool.Context, args AddCommentArgs) (*Comment, error) {
        return s.client.AddComment(ctx, args.IssueID, args.Body)
    })
    if err != nil {
        return nil, err
    }
    tools = append(tools, addCommentTool)

    return tools, nil
}
