# AGENTS.md - /

> This is the root documentation folder for the Agent Development Kit (ADK) for Go, a framework for building AI agents. It contains contribution guidelines, project overview, and community standards for developers working with the adk-go library.

## Tech Stack

- **Go** unspecified - Primary programming language for the ADK framework
- **Markdown** N/A - Documentation format for README, CONTRIBUTING, and other guides

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `/README.md` (reference) | Project overview, badges, installation instructions, and quick start guide for adk-go | First time understanding the project scope, features, and how to get started |
| `/CONTRIBUTING.md` (reference) | Contribution workflow, CLA requirements, code review process, and community guidelines | Before making any contributions, submitting PRs, or engaging with the community |

## Architecture

```
Documentation Layer (Root)
├── README.md (Project Entry Point)
│   ├── Badges & Status
│   ├── Installation Guide
│   └── Quick Start Examples
├── CONTRIBUTING.md (Contribution Process)
│   ├── CLA Requirements
│   ├── Code Review Standards
│   ├── Issue Workflow
│   └── Community Guidelines
└── Supporting Docs
    ├── License Information
    └── Community Standards
```

## Patterns

### Multi-Level Documentation

Multiple CONTRIBUTING.md files at different levels (16, 55, 90 lines) suggest hierarchical or modular documentation structure

See `/CONTRIBUTING.md` for reference.

### Badge-Driven Status

README uses shields.io badges for license, documentation, CI/CD status, and community links

See `/README.md` for reference.

### Structured Contribution Workflow

Formal process including CLA signing, issue finding, code reviews, and community guidelines

See `/CONTRIBUTING.md` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Sign the Contributor License Agreement (CLA) before submitting any code contributions | Legal requirement for accepting patches and contributions to the project |
| Follow the documented contribution workflow: find issue → discuss → implement → submit PR → code review | Ensures organized collaboration and prevents duplicate work or misaligned contributions |
| Review community guidelines before engaging with the project | Maintains healthy community standards and respectful collaboration |
| Check existing issues before starting new work to avoid duplication | Efficient resource allocation and prevents wasted effort on duplicate solutions |
| Reference the Go package documentation at pkg.go.dev/google.golang.org/adk for API details | Official API documentation source for accurate implementation guidance |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never submit contributions without signing the CLA first | CRITICAL | Legal blocker - contributions cannot be accepted without CLA, wasting time for both contributor and reviewers |
| Never bypass the code review process | HIGH | Code reviews are mandatory part of the contribution workflow to maintain code quality and consistency |
| Never ignore community guidelines when interacting with the project | HIGH | Violating community standards can lead to exclusion from the project and damages collaborative environment |
| Never start work on an issue without checking if it's already assigned or in progress | MEDIUM | Prevents duplicate effort and potential conflicts between contributors |

### Ask First

- **Making significant architectural changes to the ADK framework** - Large changes should be discussed in issues first to ensure alignment with project direction and avoid wasted effort
- **Adding new major features or capabilities** - Feature additions should be proposed and discussed to validate need and design before implementation
- **Modifying contribution guidelines or documentation structure** - Process changes affect all contributors and should have community input
- **Creating new documentation files beyond standard README/CONTRIBUTING patterns** - Documentation structure should remain consistent and organized across the project

## Commands

### view-package-docs

Open the official Go package documentation for ADK API reference

```bash
open https://pkg.go.dev/google.golang.org/adk
```

### check-ci-status

View the nightly CI/CD build status and test results

```bash
open https://github.com/google/adk-go/actions/workflows/nightly.yml
```

### view-reddit-community

Access the Reddit community for discussions and support

```bash
open https://reddit.com/r/agentdevelopmentkit
```

## Testing

Testing strategy not defined in root documentation files - likely defined in deeper project structure

Test directory: `Not specified in root documentation`

## Common Tasks

### Start Contributing to ADK

1. Read /README.md to understand project scope and features
2. Review /CONTRIBUTING.md for contribution workflow and requirements
3. Sign the Contributor License Agreement (CLA)
4. Browse issues to find 'good first issue' or 'help wanted' labels
5. Comment on issue to express interest and discuss approach
6. Fork repository and create feature branch
7. Implement changes following project patterns
8. Submit PR and engage in code review process

### Understand ADK Architecture

1. Read /README.md for high-level overview and quick start
2. Visit pkg.go.dev/google.golang.org/adk for API documentation
3. Review learning units: Agent-Based Architecture Fundamentals, Tool System, LLM Agent Implementation
4. Examine example code in quick start section
5. Study Model Abstraction and Integration patterns

### Set Up Development Environment

1. Install Go (version specified in deeper project files)
2. Clone repository: git clone https://github.com/google/adk-go
3. Review /CONTRIBUTING.md for development setup requirements
4. Install dependencies (likely via go mod download)
5. Verify setup by running tests (commands in deeper structure)

### Report an Issue or Bug

1. Search existing issues to avoid duplicates
2. Review /CONTRIBUTING.md for issue reporting guidelines
3. Create new issue with clear description and reproduction steps
4. Add relevant labels (bug, enhancement, question, etc.)
5. Engage with maintainers in issue discussion

