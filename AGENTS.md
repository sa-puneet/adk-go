# AGENTS.md - /

> This is the root documentation folder for the Agent Development Kit (ADK) for Go, a framework for building AI agents. It contains contribution guidelines, project overview, and community standards for developers working with the adk-go project.

## Tech Stack

- **Go** unspecified - Primary programming language for the ADK framework
- **Markdown** N/A - Documentation format for README, CONTRIBUTING, and other project docs

## Key Files

| File | Purpose | Read When |
|------|---------|------------|
| `/README.md` (reference) | Main project overview, installation instructions, and quick start guide for ADK-Go | First time understanding the project, checking badges/status, or finding package documentation links |
| `/CONTRIBUTING.md` (reference) | Contribution workflow, CLA requirements, code review process, and community guidelines | Before making any code contributions, submitting PRs, or engaging with the community |

## Architecture

```
Documentation Structure:

┌─────────────────┐
│   README.md     │  Entry point: Project overview, badges, links
│   (54 lines)    │  → Points to pkg.go.dev for API docs
└────────┬────────┘  → Links to examples and tutorials
         │
         ▼
┌─────────────────┐
│ CONTRIBUTING.md │  Contribution process and guidelines
│  (90/55/16 ln)  │  → CLA signing requirements
└────────┬────────┘  → Code review standards
         │            → Issue workflow
         ▼
┌─────────────────┐
│  Community      │  External links to:
│  Resources      │  - Reddit: r/agentdevelopmentkit
└─────────────────┘  - GitHub Actions (nightly checks)
                     - Apache 2.0 License
```

## Patterns

### Multi-Version Documentation

Multiple CONTRIBUTING.md files exist (90, 55, 16 lines), suggesting versioned or context-specific contribution guides for different aspects of the project

See `/CONTRIBUTING.md` for reference.

### Badge-Driven Status

README uses shields.io badges to communicate license, documentation status, and CI/CD health at a glance

See `/README.md` for reference.

### External Documentation Reference

Core API documentation hosted on pkg.go.dev rather than in-repo, keeping docs synchronized with code

See `/README.md` for reference.

## Rules

### Always (Do These)

| Rule | Rationale |
|------|------------|
| Sign the Contributor License Agreement (CLA) before submitting any code contributions | Legal requirement for accepting patches and contributions to this Google-maintained project |
| Review community guidelines before engaging with the project | Ensures all contributors understand behavioral expectations and collaboration standards |
| Follow the established contribution workflow when submitting changes | Standardized process ensures code quality, proper review, and maintainable contributions |
| Reference pkg.go.dev for authoritative API documentation | Single source of truth for Go package documentation, auto-generated from code |
| Check GitHub Actions nightly build status before assuming codebase health | Nightly checks badge indicates current CI/CD status and potential breaking changes |

### Never (Don't Do These)

| Rule | Severity | Rationale |
|------|----------|------------|
| Never submit code contributions without signing the CLA first | CRITICAL | Contributions cannot be legally accepted without CLA, will block PR merge |
| Never bypass the code review process outlined in CONTRIBUTING.md | HIGH | Code reviews are mandatory for quality assurance and knowledge sharing |
| Never duplicate API documentation in README that belongs in Go package docs | MEDIUM | Creates documentation drift; pkg.go.dev is the authoritative source for API details |
| Never ignore community guidelines when interacting with maintainers or contributors | HIGH | Violations can result in contribution rejection or community access restrictions |

### Ask First

- **Making significant architectural changes to the ADK framework** - Large changes should be discussed in issues first to align with project direction and avoid wasted effort
- **Adding new external dependencies to the project** - Dependency additions impact all users and require maintainer approval for licensing and maintenance considerations
- **Modifying contribution guidelines or documentation structure** - Meta-documentation changes affect all contributors and should be coordinated with maintainers
- **Creating new top-level documentation files in root directory** - Root directory structure is intentionally minimal; new docs may belong in subdirectories or wiki

## Commands

### view-package-docs

Open the official Go package documentation for ADK in browser

```bash
open https://pkg.go.dev/google.golang.org/adk
```

### check-ci-status

View the latest nightly build status and CI/CD results

```bash
open https://github.com/google/adk-go/actions/workflows/nightly.yml
```

### join-community

Access the ADK community on Reddit for discussions and support

```bash
open https://reddit.com/r/agentdevelopmentkit
```

## Testing

Testing information not present in root documentation files. Refer to CONTRIBUTING.md for testing requirements during code review process.

Test directory: `Not specified in root documentation`

## Common Tasks

### Start Contributing to ADK-Go

1. Read /README.md to understand project scope and architecture
2. Review /CONTRIBUTING.md for contribution workflow and requirements
3. Sign the Contributor License Agreement (CLA) as specified in CONTRIBUTING.md
4. Review community guidelines linked in CONTRIBUTING.md
5. Find an issue to work on using the 'Finding Issues to Work On' section
6. Follow the contribution workflow for submitting your changes

### Understand ADK API and Usage

1. Read /README.md for high-level overview and quick start
2. Visit pkg.go.dev link from README for detailed API documentation
3. Review learning units mentioned in folder context (Agent-Based Architecture, Tool System, etc.)
4. Check examples and tutorials referenced in README
5. Join r/agentdevelopmentkit community for questions and discussions

### Verify Project Health Before Contributing

1. Check README.md badges for license and documentation status
2. Click nightly build badge to view recent CI/CD results
3. Review open issues on GitHub for known problems
4. Check if CLA signing process is current and accessible
5. Verify pkg.go.dev documentation is up-to-date with latest release

### Prepare for Code Review

1. Ensure CLA is signed (CRITICAL - see CONTRIBUTING.md)
2. Review code review guidelines in CONTRIBUTING.md
3. Verify changes follow contribution workflow
4. Check that community guidelines have been followed
5. Ensure PR description references relevant issues
6. Confirm all required checks pass before requesting review

