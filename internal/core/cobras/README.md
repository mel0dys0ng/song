# Cobras

A powerful command-line interface (CLI) framework built on top of `spf13/cobra`, providing simplified command definition, hierarchical command management, and inheritance support. This package makes it easy to build complex CLI applications with clean, maintainable code.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
  - [Creating a Command Runner](#creating-a-command-runner)
  - [Defining Commands](#defining-commands)
  - [Registering Commands](#registering-commands)
  - [Command Hierarchy](#command-hierarchy)
  - [Flags and Arguments](#flags-and-arguments)
  - [Hook Functions](#hook-functions)
- [Command Interface](#command-interface)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Simplified Command Definition**: Cleaner API for defining commands compared to raw Cobra
- **Command Inheritance**: Child commands automatically inherit parent command configurations
- **Hierarchical Structure**: Support for multi-level command trees
- **Hook Functions**: Comprehensive lifecycle hooks (PreRun, PostRun, PersistentPreRun, PersistentPostRun)
- **Flag Binding**: Easy flag definition and binding
- **Error Handling**: Built-in error logging and panic recovery
- **Empty Command Support**: Graceful handling of commands without explicit implementations

## Installation

Ensure you have the required dependencies:

```bash
go get github.com/spf13/cobra
go get github.com/spf13/pflag
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/cobras"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
)

// Define a simple command
type HelloCommand struct{}

func (h *HelloCommand) Short() string {
    return "Say hello"
}

func (h *HelloCommand) Long() string {
    return "A command that says hello to the user"
}

func (h *HelloCommand) BindFlags(set *pflag.FlagSet) {
    set.String("name", "World", "Name to greet")
}

func (h *HelloCommand) Run(cmd *cobra.Command, args []string) {
    name, _ := cmd.Flags().GetString("name")
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    // Create command runner
    runner := cobras.New("myapp")
    
    // Register command
    runner.RegisterCommand("hello", &HelloCommand{})
    
    // Execute
    runner.Execute()
}
```

Run the application:
```bash
./myapp hello --name Alice
# Output: Hello, Alice!
```

## Architecture Overview

The Cobras package provides a layered architecture:

```
┌─────────────────┐
│ CommandRunner   │ - Manages all commands
└────────┬────────┘
         │
    ┌────┴────┐
    │ Command │ - Wraps Cobra commands
    └────┬────┘
         │
    ┌────┴────────────┐
    │ CommandInterface│ - Your command implementation
    └─────────────────┘
```

**Key Components:**
- **CommandRunner**: Central manager for all commands
- **Command**: Wrapper around Cobra commands with metadata
- **CommandInterface**: Interface you implement for each command
- **CommandChild**: Helper for building command hierarchies

## Usage Guide

### Creating a Command Runner

Create a command runner as the entry point for your CLI application:

```go
import "github.com/mel0dys0ng/song/internal/core/cobras"

func main() {
    // Create a new command runner with the application name
    runner := cobras.New("myapp")
    
    // Optionally register a custom root command
    runner.RegisterRoot(func(name string) cobras.CommandInterface {
        return &RootCommand{}
    })
}
```

### Defining Commands

Implement the `CommandInterface` for each command:

```go
type RootCommand struct{}

func (r *RootCommand) Short() string {
    return "Root command of the application"
}

func (r *RootCommand) Long() string {
    return "This is the root command that serves as the entry point for all subcommands"
}

func (r *RootCommand) BindFlags(set *pflag.FlagSet) {
    // Define global flags
    set.BoolP("verbose", "v", false, "Enable verbose output")
    set.String("config", "", "Configuration file path")
}

func (r *RootCommand) PersistentPreRun(cmd *cobra.Command, args []string) {
    // This runs before all subcommands
    verbose, _ := cmd.Flags().GetBool("verbose")
    if verbose {
        fmt.Println("Verbose mode enabled")
    }
}

func (r *RootCommand) Run(cmd *cobra.Command, args []string) {
    // Show help if no subcommand is provided
    cmd.Help()
}
```

### Registering Commands

Register commands with the runner:

```go
// Simple registration
runner.RegisterCommand("hello", &HelloCommand{})

// Register with child commands
runner.RegisterCommand("user").
    RegisterChild("create", &UserCreateCommand{}).
    RegisterChild("delete", &UserDeleteCommand{}).
    RegisterChild("list", &UserListCommand{})

// Multiple levels of hierarchy
runner.RegisterCommand("api").
    RegisterChild("v1").
    RegisterChild("users", &APIV1UsersCommand{})
```

### Command Hierarchy

Build complex command hierarchies:

```go
// Create a command tree
runner := cobras.New("git")

// git add
runner.RegisterCommand("add", &AddCommand{})

// git commit
runner.RegisterCommand("commit", &CommitCommand{})

// git remote add
runner.RegisterCommand("remote").
    RegisterChild("add", &RemoteAddCommand{}).
    RegisterChild("remove", &RemoteRemoveCommand{}).
    RegisterChild("list", &RemoteListCommand{})

// git branch create
runner.RegisterCommand("branch").
    RegisterChild("create", &BranchCreateCommand{}).
    RegisterChild("delete", &BranchDeleteCommand{})
```

### Flags and Arguments

Define and access command-line flags:

```go
type DeployCommand struct{}

func (d *DeployCommand) BindFlags(set *pflag.FlagSet) {
    // String flags
    set.String("env", "development", "Deployment environment")
    set.StringP("region", "r", "us-east-1", "Deployment region")
    
    // Bool flags
    set.BoolP("force", "f", false, "Force deployment")
    set.Bool("dry-run", false, "Simulate deployment")
    
    // Int flags
    set.IntP("replicas", "n", 3, "Number of replicas")
    
    // String slice flags
    set.StringSlice("tags", []string{}, "Deployment tags")
}

func (d *DeployCommand) Run(cmd *cobra.Command, args []string) {
    // Access flag values
    env, _ := cmd.Flags().GetString("env")
    region, _ := cmd.Flags().GetString("region")
    force, _ := cmd.Flags().GetBool("force")
    replicas, _ := cmd.Flags().GetInt("replicas")
    tags, _ := cmd.Flags().GetStringSlice("tags")
    
    fmt.Printf("Deploying to %s (%s) with %d replicas\n", env, region, replicas)
    
    // Access positional arguments
    if len(args) > 0 {
        fmt.Printf("Deploying application: %s\n", args[0])
    }
}
```

### Hook Functions

Use lifecycle hooks for setup and cleanup:

```go
type DatabaseCommand struct{}

func (d *DatabaseCommand) PersistentPreRun(cmd *cobra.Command, args []string) {
    // Runs before this command and all child commands
    fmt.Println("Initializing database connection...")
    // Setup database connection
}

func (d *DatabaseCommand) PreRun(cmd *cobra.Command, args []string) {
    // Runs before this command only (not inherited)
    fmt.Println("Validating inputs...")
}

func (d *DatabaseCommand) Run(cmd *cobra.Command, args []string) {
    // Main command logic
    fmt.Println("Executing database operation...")
}

func (d *DatabaseCommand) PostRun(cmd *cobra.Command, args []string) {
    // Runs after command completion
    fmt.Println("Cleaning up...")
}

func (d *DatabaseCommand) PersistentPostRun(cmd *cobra.Command, args []string) {
    // Runs after this command and all child commands
    fmt.Println("Closing database connection...")
}
```

## Command Interface

The `CommandInterface` defines the contract for all commands:

```go
type CommandInterface interface {
    // Short returns the short description shown in help listings
    Short() string
    
    // Long returns the long description shown in help
    Long() string
    
    // BindFlags binds command-line flags
    BindFlags(set *pflag.FlagSet)
    
    // PersistentPreRun runs before all child commands
    PersistentPreRun(cmd *cobra.Command, args []string)
    
    // PreRun runs before this command (not inherited)
    PreRun(cmd *cobra.Command, args []string)
    
    // Run executes the main command logic
    Run(cmd *cobra.Command, args []string)
    
    // PostRun runs after command completion (not inherited)
    PostRun(cmd *cobra.Command, args []string)
    
    // PersistentPostRun runs after all child commands
    PersistentPostRun(cmd *cobra.Command, args []string)
}
```

### Optional Methods

You can use `EmptyCommand` for commands that only serve as parents:

```go
// For commands that only have children
type ParentCommand struct{}

func (p *ParentCommand) Short() string {
    return "Parent command"
}

func (p *ParentCommand) Long() string {
    return "This command only serves as a parent for subcommands"
}

func (p *ParentCommand) BindFlags(set *pflag.FlagSet) {
    // No flags needed
}

// Other methods can be left empty or use default implementations
```

## Examples

### Complete Example: CLI Application

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/cobras"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
)

// Root command
type RootCommand struct{}

func (r *RootCommand) Short() string {
    return "My CLI Application"
}

func (r *RootCommand) Long() string {
    return "A comprehensive CLI application for managing resources"
}

func (r *RootCommand) BindFlags(set *pflag.FlagSet) {
    set.BoolP("verbose", "v", false, "Enable verbose output")
    set.String("config", "config.yaml", "Configuration file")
}

func (r *RootCommand) PersistentPreRun(cmd *cobra.Command, args []string) {
    verbose, _ := cmd.Flags().GetBool("verbose")
    if verbose {
        fmt.Println("Verbose mode enabled")
    }
}

func (r *RootCommand) Run(cmd *cobra.Command, args []string) {
    cmd.Help()
}

// User create command
type UserCreateCommand struct{}

func (u *UserCreateCommand) Short() string {
    return "Create a new user"
}

func (u *UserCreateCommand) Long() string {
    return "Create a new user with the specified name and email"
}

func (u *UserCreateCommand) BindFlags(set *pflag.FlagSet) {
    set.String("name", "", "User name (required)")
    set.String("email", "", "User email (required)")
    set.String("role", "user", "User role")
}

func (u *UserCreateCommand) Run(cmd *cobra.Command, args []string) {
    name, _ := cmd.Flags().GetString("name")
    email, _ := cmd.Flags().GetString("email")
    role, _ := cmd.Flags().GetString("role")
    
    if name == "" || email == "" {
        fmt.Println("Error: name and email are required")
        return
    }
    
    fmt.Printf("Creating user: %s (%s) with role: %s\n", name, email, role)
}

// User list command
type UserListCommand struct{}

func (u *UserListCommand) Short() string {
    return "List all users"
}

func (u *UserListCommand) Long() string {
    return "Display a list of all users in the system"
}

func (u *UserListCommand) BindFlags(set *pflag.FlagSet) {
    set.Int("limit", 10, "Maximum number of users to display")
    set.Int("offset", 0, "Number of users to skip")
}

func (u *UserListCommand) Run(cmd *cobra.Command, args []string) {
    limit, _ := cmd.Flags().GetInt("limit")
    offset, _ := cmd.Flags().GetInt("offset")
    
    fmt.Printf("Listing users (limit: %d, offset: %d)\n", limit, offset)
}

func main() {
    // Create command runner
    runner := cobras.New("myapp")
    
    // Register root command
    runner.RegisterRoot(func(name string) cobras.CommandInterface {
        return &RootCommand{}
    })
    
    // Register user commands
    runner.RegisterCommand("user").
        RegisterChild("create", &UserCreateCommand{}).
        RegisterChild("list", &UserListCommand{}).
        RegisterChild("delete", &cobras.EmptyCommand{}).
        RegisterChild("update", &cobras.EmptyCommand{})
    
    // Register other commands
    runner.RegisterCommand("version", &cobras.EmptyCommand{})
    runner.RegisterCommand("config", &cobras.EmptyCommand{})
    
    // Execute
    runner.Execute()
}
```

### Example: Database Migration Tool

```go
// Migration root command
type MigrationCommand struct{}

func (m *MigrationCommand) Short() string {
    return "Database migration tools"
}

func (m *MigrationCommand) Long() string {
    return "Commands for managing database migrations"
}

func (m *MigrationCommand) BindFlags(set *pflag.FlagSet) {
    set.String("database", "default", "Database name")
    set.String("env", "development", "Environment")
}

func (m *MigrationCommand) PersistentPreRun(cmd *cobra.Command, args []string) {
    // Initialize database connection
    env, _ := cmd.Flags().GetString("env")
    fmt.Printf("Connecting to database (env: %s)...\n", env)
}

// Migrate up command
type MigrateUpCommand struct{}

func (u *MigrateUpCommand) Short() string {
    return "Run pending migrations"
}

func (u *MigrateUpCommand) Long() string {
    return "Execute all pending database migrations"
}

func (u *MigrateUpCommand) BindFlags(set *pflag.FlagSet) {
    set.Int("steps", 0, "Number of migrations to run (0 = all)")
}

func (u *MigrateUpCommand) Run(cmd *cobra.Command, args []string) {
    steps, _ := cmd.Flags().GetInt("steps")
    fmt.Printf("Running migrations (steps: %d)...\n", steps)
}

// Migrate down command
type MigrateDownCommand struct{}

func (d *MigrateDownCommand) Short() string {
    return "Rollback migrations"
}

func (d *MigrateDownCommand) Long() string {
    return "Rollback the last database migration"
}

func (d *MigrateDownCommand) BindFlags(set *pflag.FlagSet) {
    set.Int("steps", 1, "Number of migrations to rollback")
}

func (d *MigrateDownCommand) Run(cmd *cobra.Command, args []string) {
    steps, _ := cmd.Flags().GetInt("steps")
    fmt.Printf("Rolling back %d migrations...\n", steps)
}

// Usage
func main() {
    runner := cobras.New("migrate")
    
    runner.RegisterCommand("migrate").
        RegisterChild("up", &MigrateUpCommand{}).
        RegisterChild("down", &MigrateDownCommand{}).
        RegisterChild("status", &cobras.EmptyCommand{})
    
    runner.Execute()
}
```

### Example: Build System CLI

```go
type BuildCommand struct{}

func (b *BuildCommand) Short() string {
    return "Build the project"
}

func (b *BuildCommand) Long() string {
    return "Build the project with specified options"
}

func (b *BuildCommand) BindFlags(set *pflag.FlagSet) {
    set.StringP("output", "o", "dist", "Output directory")
    set.BoolP("minify", "m", false, "Minify output")
    set.Bool("sourcemap", false, "Generate sourcemaps")
    set.StringP("target", "t", "es2020", "Compilation target")
}

func (b *BuildCommand) Run(cmd *cobra.Command, args []string) {
    output, _ := cmd.Flags().GetString("output")
    minify, _ := cmd.Flags().GetBool("minify")
    sourcemap, _ := cmd.Flags().GetBool("sourcemap")
    
    fmt.Printf("Building project to %s (minify: %v, sourcemap: %v)\n", 
        output, minify, sourcemap)
}

func main() {
    runner := cobras.New("build")
    
    runner.RegisterCommand("build", &BuildCommand{}).
        RegisterChild("clean", &cobras.EmptyCommand{}).
        RegisterChild("watch", &cobras.EmptyCommand{})
    
    runner.RegisterCommand("test").
        RegisterChild("unit", &cobras.EmptyCommand{}).
        RegisterChild("integration", &cobras.EmptyCommand{})
    
    runner.Execute()
}
```

## Best Practices

1. **Use Descriptive Names**: Choose clear, concise command names that describe their purpose.

2. **Provide Good Help Text**: Write helpful short and long descriptions for all commands.

3. **Use Flags Appropriately**: Use flags for optional parameters and arguments for required values.

4. **Leverage Inheritance**: Use `PersistentPreRun` and `PersistentPostRun` for shared setup/cleanup logic.

5. **Validate Early**: Perform input validation in `PreRun` before executing main logic.

6. **Keep Commands Focused**: Each command should do one thing well. Split complex commands into subcommands.

7. **Use Empty Commands for Grouping**: Use `EmptyCommand` for commands that only serve as parents.

8. **Handle Errors Gracefully**: Always check for errors and provide helpful error messages.

9. **Follow CLI Conventions**: Follow established CLI conventions (e.g., `-h` for help, `-v` for verbose).

10. **Test Commands**: Write tests for your commands to ensure they work as expected.

11. **Document Examples**: Include usage examples in the long description for complex commands.

12. **Use Consistent Flag Names**: Maintain consistency in flag naming across commands.

## Additional Resources

- [Cobra Documentation](https://github.com/spf13/cobra)
- [pflag Documentation](https://github.com/spf13/pflag)
- [Song Framework Documentation](../../README.md)
