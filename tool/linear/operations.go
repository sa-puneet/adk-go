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
    "context"
    "fmt"
)

// CreateIssue creates a new issue in Linear.
func (c *Client) CreateIssue(ctx context.Context, input CreateIssueInput) (*Issue, error) {
    query := `
        mutation IssueCreate($input: IssueCreateInput!) {
            issueCreate(input: $input) {
                success
                issue {
                    id
                    identifier
                    title
                    description
                    priority
                    url
                    createdAt
                    updatedAt
                    state {
                        id
                        name
                        type
                    }
                    assignee {
                        id
                        name
                    }
                }
            }
        }
    `

    variables := map[string]any{
        "input": input,
    }

    var resp struct {
        IssueCreate struct {
            Success bool   `json:"success"`
            Issue   *Issue `json:"issue"`
        } `json:"issueCreate"`
    }

    if err := c.do(ctx, query, variables, &resp); err != nil {
        return nil, err
    }

    if !resp.IssueCreate.Success {
        return nil, fmt.Errorf("failed to create issue")
    }

    return resp.IssueCreate.Issue, nil
}

// UpdateIssue updates an existing issue.
func (c *Client) UpdateIssue(ctx context.Context, input UpdateIssueInput) (*Issue, error) {
    query := `
        mutation IssueUpdate($id: String!, $input: IssueUpdateInput!) {
            issueUpdate(id: $id, input: $input) {
                success
                issue {
                    id
                    identifier
                    title
                    description
                    priority
                    url
                    createdAt
                    updatedAt
                    state {
                        id
                        name
                        type
                    }
                    assignee {
                        id
                        name
                    }
                }
            }
        }
    `

    // Extract ID and remove it from input map to avoid sending it in the input object if not needed,
    // but here we defined UpdateIssueInput struct which has ID field.
    // The GraphQL mutation expects `id` as a separate argument and `input` as an object.
    // We need to map our struct to the GraphQL input properly.

    // Let's create a separate struct for the GraphQL input to be clean, or use map.
    updateInput := map[string]any{}
    if input.Title != nil {
        updateInput["title"] = *input.Title
    }
    if input.Description != nil {
        updateInput["description"] = *input.Description
    }
    if input.Priority != nil {
        updateInput["priority"] = *input.Priority
    }
    if input.AssigneeID != nil {
        updateInput["assigneeId"] = *input.AssigneeID
    }
    if input.StateID != nil {
        updateInput["stateId"] = *input.StateID
    }

    variables := map[string]any{
        "id":    input.ID,
        "input": updateInput,
    }

    var resp struct {
        IssueUpdate struct {
            Success bool   `json:"success"`
            Issue   *Issue `json:"issue"`
        } `json:"issueUpdate"`
    }

    if err := c.do(ctx, query, variables, &resp); err != nil {
        return nil, err
    }

    if !resp.IssueUpdate.Success {
        return nil, fmt.Errorf("failed to update issue")
    }

    return resp.IssueUpdate.Issue, nil
}

// GetIssue retrieves an issue by ID or Identifier (e.g. LIN-123).
func (c *Client) GetIssue(ctx context.Context, id string) (*Issue, error) {
    query := `
        query Issue($id: String!) {
            issue(id: $id) {
                id
                identifier
                title
                description
                priority
                url
                createdAt
                updatedAt
                state {
                    id
                    name
                    type
                }
                assignee {
                    id
                    name
                }
            }
        }
    `

    variables := map[string]any{
        "id": id,
    }

    var resp struct {
        Issue *Issue `json:"issue"`
    }

    if err := c.do(ctx, query, variables, &resp); err != nil {
        return nil, err
    }

    if resp.Issue == nil {
        return nil, fmt.Errorf("issue not found: %s", id)
    }

    return resp.Issue, nil
}

// ListIssues lists issues with optional filtering.
func (c *Client) ListIssues(ctx context.Context, limit int, teamID *string) ([]Issue, error) {
    // Simple query for now. supporting filtering by team if provided.
    query := `
        query Issues($first: Int, $filter: IssueFilter) {
            issues(first: $first, filter: $filter) {
                nodes {
                    id
                    identifier
                    title
                    priority
                    url
                    createdAt
                    updatedAt
                    state {
                        id
                        name
                        type
                    }
                    assignee {
                        id
                        name
                    }
                }
                pageInfo {
                    hasNextPage
                    endCursor
                }
            }
        }
    `
    
    filter := map[string]any{}
    if teamID != nil {
        filter["team"] = map[string]any{
            "id": map[string]any{
                "eq": *teamID,
            },
        }
    }

    variables := map[string]any{
        "first":  limit,
        "filter": filter,
    }

    var resp struct {
        Issues IssueConnection `json:"issues"`
    }

    if err := c.do(ctx, query, variables, &resp); err != nil {
        return nil, err
    }

    return resp.Issues.Nodes, nil
}

// GetProjects lists projects.
func (c *Client) GetProjects(ctx context.Context, limit int) ([]Project, error) {
    query := `
        query Projects($first: Int) {
            projects(first: $first) {
                nodes {
                    id
                    name
                    description
                    state
                }
            }
        }
    `
    
    variables := map[string]any{
        "first": limit,
    }

    var resp struct {
        Projects struct {
            Nodes []Project `json:"nodes"`
        } `json:"projects"`
    }

    if err := c.do(ctx, query, variables, &resp); err != nil {
        return nil, err
    }

    return resp.Projects.Nodes, nil
}

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, name string, teamIds []string, description *string) (*Project, error) {
    query := `
        mutation ProjectCreate($input: ProjectCreateInput!) {
            projectCreate(input: $input) {
                success
                project {
                    id
                    name
                    description
                    state
                }
            }
        }
    `

    input := map[string]any{
        "name":    name,
        "teamIds": teamIds,
    }
    if description != nil {
        input["description"] = *description
    }

    variables := map[string]any{
        "input": input,
    }

    var resp struct {
        ProjectCreate struct {
            Success bool     `json:"success"`
            Project *Project `json:"project"`
        } `json:"projectCreate"`
    }

    if err := c.do(ctx, query, variables, &resp); err != nil {
        return nil, err
    }
    
    if !resp.ProjectCreate.Success {
        return nil, fmt.Errorf("failed to create project")
    }

    return resp.ProjectCreate.Project, nil
}

// AddComment adds a comment to an issue.
func (c *Client) AddComment(ctx context.Context, issueID string, body string) (*Comment, error) {
    query := `
        mutation CommentCreate($input: CommentCreateInput!) {
            commentCreate(input: $input) {
                success
                comment {
                    id
                    body
                    createdAt
                    user {
                        id
                        name
                    }
                }
            }
        }
    `

    input := map[string]any{
        "issueId": issueID,
        "body":    body,
    }

    variables := map[string]any{
        "input": input,
    }

    var resp struct {
        CommentCreate struct {
            Success bool     `json:"success"`
            Comment *Comment `json:"comment"`
        } `json:"commentCreate"`
    }

    if err := c.do(ctx, query, variables, &resp); err != nil {
        return nil, err
    }
    
    if !resp.CommentCreate.Success {
        return nil, fmt.Errorf("failed to add comment")
    }

    return resp.CommentCreate.Comment, nil
}
