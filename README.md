# envtoconf

[![Go Report Card](https://goreportcard.com/badge/github.com/benoahriz/envtoconf)](https://goreportcard.com/report/github.com/benoahriz/envtoconf)

A simple, lightweight CLI tool for processing Go templates with environment variables. Designed primarily for Docker containers and configuration management, `envtoconf` makes it easy to generate configuration files at runtime based on environment variables.

## Features

- 🚀 Simple and fast - minimal dependencies, quick execution
- 🐳 Perfect for Docker containers - generate configs at container startup
- 🔧 Powerful templating - uses Go's `text/template` with [Sprig](https://github.com/Masterminds/sprig) functions
- 📝 Environment variable expansion - `env` and `expandenv` functions built-in
- 🎯 Production-ready - well-tested and battle-hardened

## Installation

### Using Go Install

```bash
go install github.com/benoahriz/envtoconf@latest
```

### Download Binary

Download pre-built binaries from the [releases page](https://github.com/benoahriz/envtoconf/releases).

### Building from Source

```bash
git clone https://github.com/benoahriz/envtoconf.git
cd envtoconf
go build -o envtoconf
```

## Usage

### Basic Usage

```bash
envtoconf --template myconfig.tpl --outfile myconfig.conf
```

### Command-Line Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--template` | | `file.tpl` | Path to the source template file |
| `--outfile` | | `outfile.txt` | Path to the output file |
| `--verbose` | `-v` | `false` | Enable verbose/debug output |
| `--version` | | | Show version information |
| `--help` | | | Show help message |

### Examples

#### Example 1: Basic Environment Variable Substitution

**Template file (`app.conf.tpl`):**
```text
database_host={{ env "DB_HOST" }}
database_port={{ env "DB_PORT" }}
app_name={{ env "APP_NAME" }}
```

**Command:**
```bash
export DB_HOST="localhost"
export DB_PORT="5432"
export APP_NAME="MyApp"
envtoconf --template app.conf.tpl --outfile app.conf
```

**Output (`app.conf`):**
```text
database_host=localhost
database_port=5432
app_name=MyApp
```

#### Example 2: Using expandenv for Variable Expansion

**Template file (`script.sh.tpl`):**
```bash
#!/bin/bash
{{ expandenv "Your PATH is: $PATH" }}
{{ expandenv "Home directory: $HOME" }}
```

**Command:**
```bash
envtoconf --template script.sh.tpl --outfile script.sh
chmod +x script.sh
```

#### Example 3: Using Sprig Functions

The [Sprig library](http://masterminds.github.io/sprig/) provides 100+ template functions for strings, dates, math, and more:

**Template file (`config.yaml.tpl`):**
```yaml
app:
  name: {{ env "APP_NAME" | lower }}
  version: {{ env "APP_VERSION" | default "1.0.0" }}
  created: {{ now | date "2006-01-02" }}

database:
  host: {{ env "DB_HOST" | default "localhost" }}
  port: {{ env "DB_PORT" | default "5432" | int }}
  name: {{ env "DB_NAME" | required }}

features:
  {{- range list "logging" "metrics" "tracing" }}
  - {{ . }}
  {{- end }}
```

**Command:**
```bash
export APP_NAME="MyService"
export DB_NAME="production_db"
envtoconf --template config.yaml.tpl --outfile config.yaml --verbose
```

#### Example 4: Docker Container Usage

**Dockerfile:**
```dockerfile
FROM alpine:latest

COPY envtoconf /usr/local/bin/
COPY nginx.conf.tpl /etc/nginx/

CMD envtoconf --template /etc/nginx/nginx.conf.tpl --outfile /etc/nginx/nginx.conf && \
    nginx -g "daemon off;"
```

**Template (`nginx.conf.tpl`):**
```nginx
server {
    listen {{ env "NGINX_PORT" | default "80" }};
    server_name {{ env "SERVER_NAME" | default "localhost" }};

    location / {
        proxy_pass {{ env "BACKEND_URL" }};
    }
}
```

## Template Functions

### Built-in Environment Functions

- `env "VAR_NAME"` - Get the value of an environment variable
- `expandenv "text with $VAR"` - Expand environment variables in a string

### Sprig Functions

All [Sprig functions](http://masterminds.github.io/sprig/) are available, including:

- **String Functions**: `trim`, `upper`, `lower`, `substr`, `replace`, `split`, etc.
- **Default Values**: `default`, `empty`, `coalesce`, `required`
- **Date Functions**: `now`, `date`, `dateModify`, `dateInZone`
- **Math Functions**: `add`, `sub`, `mul`, `div`, `mod`, `max`, `min`
- **Type Conversions**: `int`, `float64`, `toString`, `atoi`
- **Lists**: `list`, `first`, `rest`, `append`, `prepend`
- **And many more!**

See the [Sprig documentation](http://masterminds.github.io/sprig/) for the complete list.

## Use Cases

- **Docker Container Configuration**: Generate configuration files when containers start
- **Kubernetes ConfigMaps**: Process templates with pod-specific environment variables
- **CI/CD Pipelines**: Create environment-specific configurations during deployment
- **Service Configuration**: Manage application configs across different environments
- **Secret Injection**: Combine with secret management tools to inject sensitive values

## Development

### Running Tests

```bash
go test -v ./...
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See [LICENSE](LICENSE) file for details.

## Author

Benjamin Rizkowsky

## TODO

- [ ] Create tests for malformed template handling
- [ ] Create option for required vars strict mode
- [ ] Add support for multiple template files in one run
- [ ] Add JSON/YAML validation for output files
