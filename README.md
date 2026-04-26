# Go Project Scaffolding CLI

A command-line tool to quickly scaffold new Go projects from templates.

## Features

- 🚀 **Fast Project Setup**: Create new Go projects from pre-built templates in seconds
- 🔄 **Smart Replacements**: Automatically replaces template placeholders with your project details
- 📂 **Directory Renaming**: Renames template directories to match your project name
- 🔧 **Go Module Handling**: Updates go.mod files and import paths correctly
- 🐳 **Docker Support**: Updates Dockerfiles and docker-compose files
- 📝 **Comprehensive Coverage**: Handles .go files, documentation, configuration, and more
- ⚡ **Built-in Validation**: Validates project names and module paths
- 🎯 **Force Override**: Option to overwrite existing directories

## Installation

### Global Installation (recommended)

Build and install the binary as `nugz` so it's available anywhere on your system:

```bash
# Build and install to your Go bin directory (must be in PATH)
go build -o nugz . && mv nugz $(go env GOPATH)/bin/nugz
```

If `$(go env GOPATH)/bin` is not in your PATH, add it to your shell profile:

```bash
# For zsh (~/.zshrc) or bash (~/.bash_profile)
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then reload your shell (`source ~/.zshrc`) and use the tool from any directory:

```bash
nugz create -name myproject -module github.com/user/myproject
nugz module -name product
```

### Local Build

```bash
# Build from source
go build -o goScaffold

# Or run directly
go run main.go [options]
```

## Usage

### Basic Usage

```bash
# Create a new project with default settings
./goScaffold create -name myproject -module github.com/user/myproject

# Create project in specific directory
./goScaffold create -name myapi -output /path/to/projects -module github.com/company/myapi

# Use custom template directory
./goScaffold create -name mycli -template templates/cliTemplate -module github.com/user/mycli

# Force overwrite existing directory
./goScaffold create -name myproject -module github.com/user/myproject -force

# Verbose output
./goScaffold create -name myproject -module github.com/user/myproject -v
```

### Command Line Options

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `-name` | Project name | - | ✅ |
| `-module` | Go module path | `github.com/user/<name>` | - |
| `-output` | Output directory | `.` (current dir) | - |
| `-template` | Template directory | `templates/goTemplate` | - |
| `-force` | Force overwrite existing directory | `false` | - |
| `-v` | Verbose output | `false` | - |
| `-h` | Show help | - | - |

## Project Name Requirements

- Must start with a letter
- Can contain letters, numbers, underscores, and hyphens
- Cannot be a reserved name (`main`, `test`, `vendor`, `internal`)
- Must be a valid Go package name

## Module Path Requirements

- Must be in format `domain.com/user/project`
- Should follow Go module naming conventions
- Can contain letters, numbers, dots, slashes, underscores, and hyphens

## Template Structure

The tool expects a template directory structure like:

```
templates/goTemplate/
├── go.mod                 # Go module file
├── cmd/
│   └── goTemplate/        # Main application (renamed to project name)
├── internal/              # Internal packages
├── README.md             # Project documentation
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Docker Compose configuration
└── ...                   # Other project files
```

### What Gets Replaced

The tool performs intelligent replacements throughout your template:

**In File Contents:**
- `goTemplate` → `{projectName}`
- `github.com/dmawardi/goTemplate` → `{modulePath}`
- `GoTemplate` → `{ProjectName}` (title case)
- `GO_TEMPLATE` → `{PROJECTNAME}` (upper case)
- `{{.ProjectName}}` → `{projectName}` (template variables)
- `{{.ModulePath}}` → `{modulePath}` (template variables)

**In Directory/File Names:**
- `goTemplate` → `{projectName}`

**File Type Handling:**
- **Text files**: Content replacement (`.go`, `.mod`, `.md`, `.yml`, etc.)
- **Binary files**: Direct copy without modification
- **Special handling**: `go.mod`, `Dockerfile`, Go source files

## Examples

### Create a REST API Project

```bash
./goScaffold \
  -name userapi \
  -module github.com/company/userapi \
  -output ~/projects
```

This creates:
- Directory: `~/projects/userapi/`
- Module: `github.com/company/userapi`
- Main binary: `cmd/userapi/main.go`

### Create a CLI Tool

```bash
./goScaffold \
  -name mytool \
  -module github.com/username/mytool \
  -template templates/cliTemplate
```

### Development/Testing

```bash
# Create a test project locally
./goScaffold -name testproject -module github.com/test/project

# Run the generated project
cd testproject
go run cmd/testproject/main.go
```

## After Project Creation

Once your project is scaffolded:

1. **Navigate to project directory:**
   ```bash
   cd {projectName}
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   ```bash
   go run cmd/{projectName}/main.go
   ```

4. **Build for production:**
   ```bash
   go build -o bin/{projectName} cmd/{projectName}/main.go
   ```

## Troubleshooting

### Common Issues

**Error: "Directory already exists"**
- Use `-force` flag to overwrite existing directories
- Or choose a different output location with `-output`

**Error: "Template directory not found"**
- Verify the template path with `-template`
- Check that the template directory exists and is accessible

**Error: "Invalid project name"**
- Ensure project name follows Go package naming rules
- Avoid reserved words and start with a letter

**Import path issues after generation**
- Verify the module path is correct
- Run `go mod tidy` in the generated project
- Check that all import statements use the new module path

### Verbose Mode

Use `-v` flag to see detailed processing information:

```bash
./goScaffold -name myproject -module github.com/user/myproject -v
```

This shows:
- Files being processed
- Replacements being made
- Any skipped files or directories

## Contributing

To add new templates or improve the tool:

1. Templates go in `templates/{templateName}/`
2. Update replacement logic in `internal/config/config.go` if needed
3. Test with various project names and module paths
4. Ensure Docker and documentation files are handled correctly

## License

[Add your license information here]
