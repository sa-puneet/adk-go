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

import "time"

// Issue represents a Linear issue.
type Issue struct {
    ID          string    `json:"id"`
    Identifier  string    `json:"identifier"` // e.g. LIN-123
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Priority    float64   `json:"priority,omitempty"`
    Status      string    `json:"status,omitempty"` // simplified from State.name
    State       *WorkflowState `json:"state,omitempty"`
    Assignee    *User     `json:"assignee,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    URL         string    `json:"url"`
}

// User represents a Linear user.
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// WorkflowState represents a workflow state (e.g. Backlog, Todo, In Progress, Done).
type WorkflowState struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"` // e.g. backlog, unstarted, started, completed, canceled
}

// Project represents a Linear project.
type Project struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    State       string `json:"state"`
}

// Team represents a Linear team.
type Team struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Key  string `json:"key"`
}

// Comment represents a comment on an issue.
type Comment struct {
    ID        string    `json:"id"`
    Body      string    `json:"body"`
    User      *User     `json:"user,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
}

// PageInfo for pagination
type PageInfo struct {
    HasNextPage bool   `json:"hasNextPage"`
    EndCursor   string `json:"endCursor"`
}

// IssueConnection for listing issues
type IssueConnection struct {
    Nodes    []Issue  `json:"nodes"`
    PageInfo PageInfo `json:"pageInfo"`
}

// CreateIssueInput defines arguments for creating an issue.
type CreateIssueInput struct {
    TeamID      string  `json:"teamId"`
    Title       string  `json:"title"`
    Description *string `json:"description,omitempty"`
    Priority    *int    `json:"priority,omitempty"`
    AssigneeID  *string `json:"assigneeId,omitempty"`
    StateID     *string `json:"stateId,omitempty"`
}

// UpdateIssueInput defines arguments for updating an issue.
type UpdateIssueInput struct {
    ID          string  `json:"id"`
    Title       *string `json:"title,omitempty"`
    Description *string `json:"description,omitempty"`
    Priority    *int    `json:"priority,omitempty"`
    AssigneeID  *string `json:"assigneeId,omitempty"`
    StateID     *string `json:"stateId,omitempty"`
}

// IssuesFilter defines filters for listing issues.
type IssuesFilter struct {
    First *int `json:"first,omitempty"`
    // We can add more filters if needed (e.g., team ID, assignee ID)
}

// GetIssueArgs defines arguments for GetIssue tool.
type GetIssueArgs struct {
    ID string `json:"id"`
}

// ListIssuesArgs defines arguments for ListIssues tool.
type ListIssuesArgs struct {
    Limit  int     `json:"limit,omitempty"`
    TeamID *string `json:"teamId,omitempty"`
}

// GetProjectsArgs defines arguments for GetProjects tool.
type GetProjectsArgs struct {
    Limit int `json:"limit,omitempty"`
}

// CreateProjectArgs defines arguments for CreateProject tool.
type CreateProjectArgs struct {
    Name        string   `json:"name"`
    TeamIDs     []string `json:"teamIds"`
    Description *string  `json:"description,omitempty"`
}

// AddCommentArgs defines arguments for AddComment tool.
type AddCommentArgs struct {
    IssueID string `json:"issueId"`
    Body    string `json:"body"`
}
