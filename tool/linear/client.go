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
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const (
    defaultEndpoint = "https://api.linear.app/graphql"
)

// Client handles communication with the Linear API.
type Client struct {
    apiKey     string
    endpoint   string
    httpClient *http.Client
}

// ClientOption configures the Linear client.
type ClientOption func(*Client)

// WithEndpoint sets a custom GraphQL endpoint.
func WithEndpoint(endpoint string) ClientOption {
    return func(c *Client) {
        c.endpoint = endpoint
    }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
    return func(c *Client) {
        c.httpClient = client
    }
}

// NewClient creates a new Linear API client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
    c := &Client{
        apiKey:     apiKey,
        endpoint:   defaultEndpoint,
        httpClient: http.DefaultClient,
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

type graphQLRequest struct {
    Query     string         `json:"query"`
    Variables map[string]any `json:"variables,omitempty"`
}

type graphQLResponse struct {
    Data   json.RawMessage `json:"data"`
    Errors []graphQLError  `json:"errors,omitempty"`
}

type graphQLError struct {
    Message string `json:"message"`
    // We can add locations and path if needed, but message is usually enough for basic error handling
}

func (c *Client) do(ctx context.Context, query string, variables map[string]any, responseData any) error {
    reqBody, err := json.Marshal(graphQLRequest{
        Query:     query,
        Variables: variables,
    })
    if err != nil {
        return fmt.Errorf("failed to marshal request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(reqBody))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    var apiResp graphQLResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }

    if len(apiResp.Errors) > 0 {
        return fmt.Errorf("graphql error: %s", apiResp.Errors[0].Message)
    }

    if err := json.Unmarshal(apiResp.Data, responseData); err != nil {
        return fmt.Errorf("failed to unmarshal data: %w", err)
    }

    return nil
}
