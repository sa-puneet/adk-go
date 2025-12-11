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
    "bytes"
    "context"
    "io"
    "net/http"
    "testing"

    "google.golang.org/genai"

    "google.golang.org/adk/agent"
    "google.golang.org/adk/internal/toolinternal"
    "google.golang.org/adk/memory"
    "google.golang.org/adk/session"
    "google.golang.org/adk/tool"
)

type mockTransport struct {
    roundTripFunc func(*http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.roundTripFunc(req)
}

type mockToolContext struct {
    context.Context
}

func (m *mockToolContext) FunctionCallID() string                                                       { return "test-call-id" }
func (m *mockToolContext) Actions() *session.EventActions                                               { return nil }
func (m *mockToolContext) SearchMemory(ctx context.Context, query string) (*memory.SearchResponse, error) { return nil, nil }
func (m *mockToolContext) AgentName() string                                                            { return "test-agent" }
func (m *mockToolContext) ReadonlyState() session.ReadonlyState                                         { return nil }
func (m *mockToolContext) State() session.State                                                         { return nil }
func (m *mockToolContext) Artifacts() agent.Artifacts                                                   { return nil }
func (m *mockToolContext) InvocationID() string                                                         { return "test-invocation-id" }
func (m *mockToolContext) UserContent() *genai.Content                                                  { return nil }
func (m *mockToolContext) AppName() string                                                              { return "test-app" }
func (m *mockToolContext) Branch() string                                                               { return "test-branch" }
func (m *mockToolContext) SessionID() string                                                            { return "test-session-id" }
func (m *mockToolContext) UserID() string                                                               { return "test-user-id" }

func TestNew(t *testing.T) {
    cfg := Config{
        APIKey: "test-api-key",
    }
    ts, err := New(cfg)
    if err != nil {
        t.Fatalf("New() error = %v", err)
    }

    if ts == nil {
        t.Fatal("New() returned nil")
    }

    if ts.Name() != "linear" {
        t.Errorf("Name() = %v, want %v", ts.Name(), "linear")
    }
}

func TestTools(t *testing.T) {
    cfg := Config{
        APIKey: "test-api-key",
    }
    ts, err := New(cfg)
    if err != nil {
        t.Fatalf("New() error = %v", err)
    }

    tools, err := ts.Tools(agent.ReadonlyContext(nil))
    if err != nil {
        t.Fatalf("Tools() error = %v", err)
    }

    expectedTools := map[string]bool{
        "linear_create_issue":   false,
        "linear_update_issue":   false,
        "linear_get_issue":      false,
        "linear_list_issues":    false,
        "linear_get_projects":   false,
        "linear_create_project": false,
        "linear_add_comment":    false,
    }

    if len(tools) != len(expectedTools) {
        t.Errorf("Tools() returned %d tools, want %d", len(tools), len(expectedTools))
    }

    for _, tool := range tools {
        if _, ok := expectedTools[tool.Name()]; !ok {
            t.Errorf("Unexpected tool: %s", tool.Name())
        }
        expectedTools[tool.Name()] = true
    }

    for name, found := range expectedTools {
        if !found {
            t.Errorf("Tool not found: %s", name)
        }
    }
}

func TestCreateIssueTool(t *testing.T) {
    // Mock HTTP client
    mockClient := &http.Client{
        Transport: &mockTransport{
            roundTripFunc: func(req *http.Request) (*http.Response, error) {
                // Verify request
                if req.Method != "POST" {
                    t.Errorf("Method = %v, want POST", req.Method)
                }
                if req.Header.Get("Authorization") != "test-api-key" {
                    t.Errorf("Authorization header = %v, want test-api-key", req.Header.Get("Authorization"))
                }
                // Return mock response
                return &http.Response{
                    StatusCode: 200,
                    Body:       http.NoBody, // We need to mock the body, see below
                    Header:     make(http.Header),
                }, nil
            },
        },
    }

    // Update mockClient to return proper JSON
    mockClient.Transport.(*mockTransport).roundTripFunc = func(req *http.Request) (*http.Response, error) {
        jsonResp := `{
            "data": {
                "issueCreate": {
                    "success": true,
                    "issue": {
                        "id": "issue-123",
                        "identifier": "LIN-123",
                        "title": "Test Issue",
                        "url": "https://linear.app/issue/LIN-123"
                    }
                }
            }
        }`
        return &http.Response{
            StatusCode: 200,
            Body:       io.NopCloser(bytes.NewBufferString(jsonResp)),
            Header:     make(http.Header),
        }, nil
    }

    client := NewClient("test-api-key", WithHTTPClient(mockClient))
    ts := &linearToolset{client: client}

    tools, err := ts.Tools(nil)
    if err != nil {
        t.Fatalf("Tools() error = %v", err)
    }

    var createTool tool.Tool
    for _, tl := range tools {
        if tl.Name() == "linear_create_issue" {
            createTool = tl
            break
        }
    }

    if createTool == nil {
        t.Fatal("linear_create_issue tool not found")
    }

    // Call the tool
    args := map[string]any{
        "teamId": "team-123",
        "title":  "Test Issue",
    }

    fnTool, ok := createTool.(toolinternal.FunctionTool)
    if !ok {
        t.Fatal("tool does not implement toolinternal.FunctionTool")
    }
    
    result, err := fnTool.Run(&mockToolContext{Context: context.Background()}, args)
    if err != nil {
        t.Fatalf("Run() error = %v", err)
    }

    // Check result
    resMap := result

    // If conversion worked, it should be the fields of Issue.
    // But functiontool logic says:
    // If output schema is present (it is inferred), it validates.
    // If it fails to convert to map (e.g. primitive), it wraps in {"result": val}.
    // Struct should convert to map.

    if resMap["id"] != "issue-123" {
        t.Errorf("Result id = %v, want issue-123", resMap["id"])
    }
    if resMap["identifier"] != "LIN-123" {
        t.Errorf("Result identifier = %v, want LIN-123", resMap["identifier"])
    }
}
