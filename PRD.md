# Product Requirements Document: envtoconf v2.0

**Status**: Draft
**Last Updated**: 2025-11-27
**Author**: Benjamin Rizkowsky
**Version**: 2.0

---

## Executive Summary

envtoconf is being modernized to remain the simplest, most secure way to generate configuration files from templates in containerized environments. While tools like [gomplate](https://github.com/hairyhenderson/gomplate) offer extensive datasource integrations, envtoconf will focus on doing one thing exceptionally well: secure, simple template rendering with first-class secret management.

### Vision Statement

**"The secure, lightweight config generator that developers trust for containers and cloud-native applications."**

### Key Differentiators
- **Simplicity First**: No complex datasource integrations - focus on env vars and secrets
- **Security Built-in**: Native keychain/vault integration inspired by [envchain](https://github.com/sorah/envchain)
- **Container Optimized**: Minimal footprint, fast startup, designed for Docker/K8s
- **Developer Friendly**: Great error messages, validation, and DX

---

## Problem Statement

### Current Challenges

1. **Outdated Technology Stack**
   - Using deprecated `dep` (Gopkg.toml) instead of Go modules
   - Dependencies from 2017 (Sprig 2.12.0, Logrus 1.0.2)
   - Vendored dependencies bloat the repository
   - Not usable as a library (only CLI)

2. **Security Gaps**
   - No secure secret storage - users must expose secrets as env vars
   - Secrets visible in process listings and container inspect
   - No integration with modern secret management tools

3. **Limited Functionality**
   - Single template processing only
   - No configuration file support
   - Poor error messages without context
   - No validation of generated configs

4. **Developer Experience**
   - Fatal errors everywhere (not testable)
   - No dry-run or preview mode
   - No shell completion
   - Limited observability

---

## Goals & Non-Goals

### Goals

**Phase 1: Modernization** (v2.0)
- ✅ Migrate to Go modules
- ✅ Update all dependencies to latest stable versions
- ✅ Remove vendoring
- ✅ Add proper error handling (library-friendly)
- ✅ Comprehensive test coverage (>80%)
- ✅ CI/CD with GitHub Actions
- ✅ Modern CLI framework (Cobra)
- ✅ Structured logging with levels

**Phase 2: Core Features** (v2.1)
- ✅ Multiple template processing
- ✅ Configuration file support (.envtoconf.yaml)
- ✅ Stdin/stdout support for piping
- ✅ Better error messages with line numbers and context
- ✅ Dry-run mode
- ✅ Shell completion (bash, zsh, fish)

**Phase 3: Secret Management** (v2.2)
- ✅ Namespace-based secret storage
- ✅ OS keychain integration (macOS, Linux)
- ✅ Interactive secret management
- ✅ Hybrid mode: env vars + keychain secrets
- ✅ Secret backend abstraction

**Phase 4: Validation & Polish** (v2.3)
- ✅ JSON/YAML output validation
- ✅ Variable usage analysis
- ✅ Watch mode for development
- ✅ Performance optimizations

### Non-Goals

- ❌ **Remote datasources** (HTTP, S3, Consul, Vault polling) - Keep it simple; users can fetch data themselves
- ❌ **Complex DSL** - Stick to Go templates + Sprig
- ❌ **Plugin system** - Avoid complexity
- ❌ **GUI or web interface** - CLI only
- ❌ **Backwards compatibility** with v1.x flags (will provide migration guide)

---

## User Personas

### 1. **Container Dev (Primary)**
**Maya - Platform Engineer**
- Builds Docker images for microservices
- Needs config generation at container startup
- Values simplicity and small image size
- Uses K8s secrets but wants better local dev experience

### 2. **DevOps Engineer (Primary)**
**Alex - SRE**
- Manages deployment pipelines
- Needs secure credential injection
- Uses multiple environments (dev, staging, prod)
- Values observability and error handling

### 3. **Security-Conscious Developer (Secondary)**
**Sam - Security Engineer**
- Audits applications for secret exposure
- Needs credentials out of env vars and logs
- Wants integration with enterprise secret management
- Values compliance and audit trails

---

## Competitive Analysis

### vs Gomplate
| Feature | envtoconf v2 | gomplate |
|---------|--------------|----------|
| Template Engine | Go + Sprig | Go + custom (200+ funcs) |
| Datasources | Env vars, Keychain, Stdin | 15+ sources (HTTP, Vault, AWS, etc.) |
| Secret Storage | Native keychain | External (Vault) |
| Binary Size | <5MB | ~20MB |
| Use Case | Containers, simple configs | Complex multi-source configs |
| Learning Curve | Low | Medium-High |

**Strategy**: Stay focused on container/secret use case; refer power users to gomplate for complex needs.

### vs Envchain
| Feature | envtoconf v2 | envchain |
|---------|--------------|----------|
| Language | Go | C |
| Template Support | Yes (core feature) | No |
| Secret Storage | Keychain + future backends | Keychain only |
| Cross-platform | Linux, macOS, Windows | Linux, macOS |
| Use Case | Config generation | Command wrapping |

**Strategy**: Combine envchain's security model with template rendering - best of both worlds.

---

## Feature Specifications

### F1: Modernized Architecture

**Priority**: P0 (Must Have)
**Phase**: 1 (v2.0)

#### Requirements
- Migrate to Go modules (go.mod/go.sum)
- Update dependencies:
  - Sprig: latest v3.x
  - Logrus → structured logging (slog or zerolog)
  - Kingpin → Cobra
- Remove vendor directory
- Refactor to library + CLI:
  ```
  pkg/
    template/    # Core template engine
    secret/      # Secret management
    config/      # Config file parsing
  cmd/
    envtoconf/   # CLI entrypoint
  ```
- Replace `log.Fatal` with proper error returns
- Add context.Context support throughout

#### Success Criteria
- [ ] `go.mod` with Go 1.21+
- [ ] All tests pass with updated dependencies
- [ ] Can be imported as library: `import "github.com/benoahriz/envtoconf/pkg/template"`
- [ ] Zero `log.Fatal` in library code
- [ ] CI passing on GitHub Actions (test, lint, build)

---

### F2: Multiple Template Processing

**Priority**: P0 (Must Have)
**Phase**: 2 (v2.1)

#### Requirements
- Support multiple input/output pairs in single invocation
- Directory-based processing: `--input-dir` and `--output-dir`
- Glob pattern support: `--template "configs/*.tpl"`
- Preserve directory structure in output

#### CLI Examples
```bash
# Multiple explicit files
envtoconf --template app.tpl --outfile app.conf \
          --template db.tpl --outfile db.conf

# Directory processing
envtoconf --input-dir ./templates --output-dir ./configs

# Glob patterns
envtoconf --template "configs/*.yaml.tpl" --output-dir ./rendered
```

#### Success Criteria
- [ ] Process 100 templates in <1s
- [ ] Clear error messages showing which template failed
- [ ] Atomic writes (don't corrupt output on failure)

---

### F3: Configuration File Support

**Priority**: P1 (Should Have)
**Phase**: 2 (v2.1)

#### Requirements
- Support `.envtoconf.yaml`, `.envtoconf.yml`, `.envtoconf.json`
- Cascade: CLI flags > env vars > config file > defaults
- Config file schema:

```yaml
# .envtoconf.yaml
templates:
  - input: app.tpl
    output: /etc/app/config.conf
  - input: nginx.tpl
    output: /etc/nginx/nginx.conf

# Secret namespaces to load
secrets:
  namespaces:
    - production-db
    - api-keys

# Template options
options:
  strict: true          # Fail on missing variables
  validate: yaml        # Validate output as YAML
  backup: true          # Create .bak before overwrite

# Logging
log:
  level: info
  format: json
```

#### Success Criteria
- [ ] Auto-discover config in current dir or `$HOME/.config/envtoconf/`
- [ ] Config file validation with helpful errors
- [ ] `envtoconf config validate` command

---

### F4: Secret Management (Keychain Integration)

**Priority**: P0 (Must Have) - **This is the differentiator!**
**Phase**: 3 (v2.2)

#### Requirements

**Secret Storage Backends**:
1. **macOS**: Keychain via [go-keychain](https://github.com/keybase/go-keychain)
2. **Linux**: Secret Service API via [99designs/keyring](https://pkg.go.dev/github.com/99designs/keyring)
3. **Windows**: Windows Credential Manager
4. **Future**: HashiCorp Vault, AWS Secrets Manager (pluggable)

**Namespace-Based Organization** (like envchain):
```bash
# Set secrets in a namespace
envtoconf secret set production-db DB_PASSWORD DB_USER

# List namespaces
envtoconf secret list

# Show secrets in namespace (masked)
envtoconf secret show production-db

# Delete namespace
envtoconf secret delete production-db

# Render template with secrets from namespace
envtoconf render --template app.tpl --secrets production-db --outfile app.conf
```

**Template Functions**:
```go
// Existing
{{ env "PUBLIC_VAR" }}

// New secret functions
{{ secret "production-db" "DB_PASSWORD" }}
{{ secretDefault "staging-db" "DB_PORT" "5432" }}

// Hybrid approach - try secret, fallback to env
{{ secretOrEnv "api-keys" "API_KEY" }}
```

**Security Features**:
- Secrets never logged (even in verbose mode)
- Secrets never in error messages (show "[REDACTED]")
- Secrets cleared from memory after use
- Audit log option: record secret access (timestamp, namespace, key)

#### Success Criteria
- [ ] Secrets stored in OS keychain, not filesystem
- [ ] Interactive prompts for secret entry (with confirmation)
- [ ] Secrets not visible in `ps aux` or container inspect
- [ ] Support for secret rotation (update existing secret)
- [ ] Import/export for team sharing (encrypted)

---

### F5: Enhanced Error Handling

**Priority**: P1 (Should Have)
**Phase**: 2 (v2.1)

#### Requirements
- Template parse errors with line/column numbers
- Missing variable errors with suggestions
- Validation errors with context
- Color-coded terminal output

#### Examples
```
Error: template parse failed
  File: /etc/app/config.tpl:15:3
  Line: database_host={{ env "DB_HOST" }
  Error: unclosed action

Suggestion: Missing closing '}}' on line 15
```

```
Error: undefined variable
  Template: app.tpl:23
  Variable: DB_PASWORD

Did you mean?
  - DB_PASSWORD (in namespace: production-db)
  - DB_PORT (env var)
```

#### Success Criteria
- [ ] All errors include actionable next steps
- [ ] Color output respects `NO_COLOR` env var
- [ ] JSON error output for CI: `--error-format=json`

---

### F6: Validation & Dry-Run

**Priority**: P1 (Should Have)
**Phase**: 4 (v2.3)

#### Requirements

**Dry-Run Mode**:
```bash
envtoconf --template app.tpl --dry-run
# Outputs rendered config to stdout without writing file
# Shows what would be written with file paths
```

**Output Validation**:
```bash
envtoconf --template app.yaml.tpl --validate yaml --outfile app.yaml
# Validates output is valid YAML before writing
# Supports: yaml, json, toml, xml
```

**Variable Analysis**:
```bash
envtoconf analyze --template app.tpl
# Lists all variables referenced
# Shows: variable name, source (env/secret/undefined), value status
```

Example output:
```
Variables in app.tpl:
  ✓ DB_HOST           env var         (set: "localhost")
  ✓ DB_PASSWORD       secret          (namespace: prod-db)
  ✗ API_ENDPOINT      undefined       (not set)
  ✓ APP_NAME          env var         (set: "myapp")
```

#### Success Criteria
- [ ] Validation catches 100% of malformed YAML/JSON
- [ ] Analyze command helps debug missing vars
- [ ] Dry-run mode useful for testing templates locally

---

### F7: Developer Experience Improvements

**Priority**: P2 (Nice to Have)
**Phase**: 4 (v2.3)

#### Requirements

**Shell Completion**:
```bash
# Install completion
envtoconf completion bash > /etc/bash_completion.d/envtoconf

# Auto-complete flags, files, namespaces
envtoconf --template <TAB>
envtoconf secret show <TAB>  # Shows namespaces
```

**Watch Mode**:
```bash
envtoconf watch --template app.tpl --outfile app.conf
# Re-renders on template file change
# Useful for local development
```

**Better Logging**:
- Structured JSON logs: `--log-format=json`
- Log levels: debug, info, warn, error
- Request ID for tracing: `--request-id=abc123`

**Stdin/Stdout Support**:
```bash
# Read template from stdin
cat app.tpl | envtoconf --template - --outfile app.conf

# Write to stdout
envtoconf --template app.tpl --outfile -

# Pipeline support
cat app.tpl | envtoconf -t - -o - | kubectl apply -f -
```

---

## Technical Architecture

### High-Level Design

```
┌─────────────────────────────────────────────┐
│              CLI (Cobra)                    │
│  - Flag parsing                             │
│  - Config file loading                      │
│  - Command routing                          │
└─────────────────┬───────────────────────────┘
                  │
    ┌─────────────┼─────────────┐
    │             │             │
┌───▼────┐   ┌────▼───┐   ┌────▼────┐
│Template│   │Secret  │   │Validator│
│Engine  │   │Manager │   │         │
│        │   │        │   │         │
│- Parse │◄──┤- Get   │   │- YAML   │
│- Render│   │- Set   │   │- JSON   │
│- Funcs │   │- List  │   │- Schema │
└────────┘   └───┬────┘   └─────────┘
                 │
         ┌───────┴────────┐
         │                │
    ┌────▼─────┐   ┌──────▼──────┐
    │Keychain  │   │Vault/AWS    │
    │Backend   │   │Backend(TBD) │
    │(macOS/   │   │             │
    │ Linux)   │   │             │
    └──────────┘   └─────────────┘
```

### Package Structure

```
envtoconf/
├── cmd/
│   └── envtoconf/
│       └── main.go              # CLI entrypoint
├── pkg/
│   ├── template/
│   │   ├── engine.go            # Template parser/renderer
│   │   ├── functions.go         # Custom template functions
│   │   └── engine_test.go
│   ├── secret/
│   │   ├── manager.go           # Secret CRUD operations
│   │   ├── backend.go           # Backend interface
│   │   ├── keychain_darwin.go   # macOS implementation
│   │   ├── keychain_linux.go    # Linux implementation
│   │   ├── keychain_windows.go  # Windows implementation
│   │   └── vault.go             # Future: Vault backend
│   ├── config/
│   │   ├── config.go            # Config file parsing
│   │   └── schema.go            # Config validation
│   ├── validator/
│   │   ├── yaml.go
│   │   ├── json.go
│   │   └── validator.go
│   └── renderer/
│       ├── renderer.go          # Orchestrates template + secrets
│       └── batch.go             # Multi-file processing
├── internal/
│   ├── cli/
│   │   ├── root.go              # Cobra root command
│   │   ├── render.go            # Render command
│   │   ├── secret.go            # Secret management commands
│   │   ├── analyze.go           # Analysis commands
│   │   └── completion.go        # Shell completion
│   └── logger/
│       └── logger.go            # Structured logging
├── examples/
│   ├── docker/
│   ├── kubernetes/
│   └── simple/
├── docs/
│   ├── migration-guide.md
│   ├── secret-management.md
│   └── api.md
├── .github/
│   └── workflows/
│       ├── test.yml
│       ├── release.yml
│       └── lint.yml
├── go.mod
├── go.sum
├── README.md
├── PRD.md
└── LICENSE
```

### Technology Decisions

| Component | Choice | Rationale |
|-----------|--------|-----------|
| CLI Framework | [Cobra](https://github.com/spf13/cobra) | Industry standard, great docs, completion support |
| Config Format | YAML/JSON | Familiar to ops teams, good validation tools |
| Template Engine | `text/template` + Sprig v3 | Keep existing, just upgrade Sprig |
| Logging | `log/slog` (stdlib) | Native in Go 1.21+, structured, zero deps |
| Keychain (macOS) | [keybase/go-keychain](https://github.com/keybase/go-keychain) | Battle-tested, maintained |
| Keychain (Linux) | [99designs/keyring](https://pkg.go.dev/github.com/99designs/keyring) | Multi-backend support |
| Testing | `testing` + [testify](https://github.com/stretchr/testify) | Assertions library, widely used |
| Validation | [goccy/go-yaml](https://github.com/goccy/go-yaml) | Best YAML validator in Go |

---

## Migration Path

### Breaking Changes from v1.x

1. **CLI Flags Changed**:
   - Old: `--template`, `--outfile`, `-v`
   - New: `-t/--template`, `-o/--output`, `--verbose`
   - Reasoning: Align with common CLI conventions

2. **Behavior Changes**:
   - Undefined variables now error by default (use `--allow-undefined` to revert)
   - Exit codes: 0=success, 1=error, 2=validation failure

3. **Removed**:
   - Vendored dependencies

### Migration Guide

```bash
# Old v1.x command
envtoconf --template file.tpl --outfile out.txt -v

# New v2.0 command (compatibility mode)
envtoconf render --template file.tpl --output out.txt --verbose

# Or use new shorthand
envtoconf render -t file.tpl -o out.txt -v

# Config file approach (recommended)
cat > .envtoconf.yaml <<EOF
templates:
  - input: file.tpl
    output: out.txt
log:
  level: debug
EOF
envtoconf render
```

**Auto-migration tool**:
```bash
# Scan current usage and generate config file
envtoconf migrate --from-flags "envtoconf --template file.tpl --outfile out.txt"
# Outputs: .envtoconf.yaml with equivalent config
```

---

## Success Metrics

### Launch Criteria (v2.0)

- [ ] 90% test coverage (unit + integration)
- [ ] All existing v1.x tests pass (or have migration path)
- [ ] Documentation complete (README, migration guide, examples)
- [ ] GitHub Actions CI passing
- [ ] Binaries built for Linux, macOS, Windows (amd64 + arm64)
- [ ] Docker image published (<5MB Alpine-based)

### Adoption Metrics

**6 months post-launch**:
- 1,000+ downloads via `go install`
- 10+ GitHub stars (from current unknown)
- 5+ community contributions (issues, PRs)
- Usage in 3+ public projects/blog posts

**12 months**:
- 5,000+ downloads
- 50+ GitHub stars
- Package included in 1+ Docker base image (e.g., Alpine, distroless)

### Performance Targets

- Single template render: <10ms
- 100 templates: <1s
- Binary size: <5MB (static binary)
- Memory usage: <50MB for typical workloads
- Cold start (Docker): <100ms

---

## Timeline & Phases

### Phase 1: Foundation (4-6 weeks)
**Goal**: Modernize stack, establish quality baseline

- [ ] Week 1-2: Go modules migration, dependency updates
- [ ] Week 3: Refactor to library + CLI structure
- [ ] Week 4: Cobra integration, new CLI structure
- [ ] Week 5: Testing infrastructure, GitHub Actions
- [ ] Week 6: Documentation, migration guide

**Deliverable**: v2.0-alpha with feature parity to v1.x

### Phase 2: Core Features (6-8 weeks)
**Goal**: Multi-template, config files, better DX

- [ ] Week 1-2: Config file support
- [ ] Week 3-4: Multiple template processing
- [ ] Week 5-6: Error handling improvements
- [ ] Week 7: Stdin/stdout, shell completion
- [ ] Week 8: Testing, docs

**Deliverable**: v2.1-beta

### Phase 3: Secret Management (8-10 weeks)
**Goal**: Keychain integration - the big differentiator

- [ ] Week 1-2: Secret backend interface design
- [ ] Week 3-4: macOS Keychain implementation
- [ ] Week 5-6: Linux Secret Service implementation
- [ ] Week 7: Template functions for secrets
- [ ] Week 8: Interactive secret management CLI
- [ ] Week 9-10: Security audit, testing, docs

**Deliverable**: v2.2-rc1

### Phase 4: Polish (4 weeks)
**Goal**: Validation, analysis, watch mode

- [ ] Week 1-2: Output validation (YAML/JSON)
- [ ] Week 3: Variable analysis, dry-run
- [ ] Week 4: Watch mode, performance tuning

**Deliverable**: v2.3 (stable)

### Total Timeline: ~6 months to v2.3

---

## Open Questions & Decisions Needed

### Strategic
- [ ] **Backwards compatibility**: How strict? Offer v1 compatibility mode?
- [ ] **Naming**: Keep "envtoconf" or rebrand? (e.g., "secconf", "tplsec")
- [ ] **License**: Stay with current license or change?

### Technical
- [ ] **Secret backends priority**: macOS/Linux first, or Windows too?
- [ ] **Vault integration**: MVP for v2.2 or defer to v2.4?
- [ ] **Template caching**: Worth the complexity for performance?
- [ ] **Concurrent rendering**: Parallel template processing for large batches?

### Product
- [ ] **Target audience**: Broaden beyond Docker/containers?
- [ ] **Pricing/Support**: Stay pure OSS or offer enterprise support?
- [ ] **Integrations**: Focus on which platforms? (K8s operator, Helm plugin, etc.)

---

## Risks & Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Breaking changes alienate users | High | Medium | Excellent migration guide, compatibility mode, 6mo deprecation |
| Secret backend bugs expose credentials | Critical | Low | Security audit, extensive testing, gradual rollout |
| Scope creep (trying to compete with gomplate) | Medium | High | Strict adherence to non-goals, focus on simplicity |
| Platform-specific keychain issues | Medium | Medium | Abstraction layer, fallback to env vars, good error messages |
| Adoption slow due to established alternatives | Medium | Medium | Focus on differentiator (secrets), content marketing, examples |

---

## Appendix

### Research References

- [Gomplate Documentation](https://docs.gomplate.ca/)
- [Gomplate GitHub](https://github.com/hairyhenderson/gomplate)
- [Envchain GitHub](https://github.com/sorah/envchain)
- [99designs/keyring](https://pkg.go.dev/github.com/99designs/keyring)
- [Keybase go-keychain](https://github.com/keybase/go-keychain)

### Similar Tools

- **gomplate**: Feature-rich, complex datasources
- **envsubst**: GNU tool, very basic, shell-based
- **confd**: Focused on service discovery (etcd, consul)
- **envchain**: Secret management only, no templating
- **chamber**: AWS Parameter Store wrapper

### Glossary

- **Namespace**: Logical grouping of secrets (e.g., "production-db", "api-keys")
- **Backend**: Secret storage implementation (keychain, vault, etc.)
- **Dry-run**: Preview mode that doesn't write files
- **Strict mode**: Fail on undefined variables instead of rendering empty string
