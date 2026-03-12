# Metas

A comprehensive metadata management system that provides application configuration, environment detection, and runtime information. This package serves as the central source of truth for application identity, deployment environment, and configuration settings.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
  - [Creating Metadata](#creating-metadata)
  - [Accessing Metadata](#accessing-metadata)
  - [Environment Variables](#environment-variables)
  - [Configuration Management](#configuration-management)
- [Metadata Components](#metadata-components)
- [Configuration Options](#configuration-options)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Application Identity**: Centralized application name and type management
- **Environment Detection**: Automatic detection of deployment environment (local, test, staging, production)
- **Runtime Information**: Access to node ID, region, zone, and provider information
- **Configuration Management**: Support for YAML and etcd configuration sources
- **Environment Variables**: Automatic loading of configuration from environment variables
- **Singleton Pattern**: Global metadata access throughout the application
- **Validation**: Built-in validation for application names and configuration paths
- **IP Detection**: Automatic detection of local IP address
- **Path Management**: Automatic path resolution for logs and configuration

## Installation

The metas package is part of the core framework and has no external dependencies beyond the standard library.

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func main() {
    // Create metadata
    mt := metas.New(&metas.Options{
        App:  "myapp",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })
    
    // Access metadata
    fmt.Printf("App: %s\n", mt.App())
    fmt.Printf("Kind: %s\n", mt.Kind())
    fmt.Printf("Mode: %s\n", mt.Mode())
    fmt.Printf("Node: %s\n", mt.Node())
    fmt.Printf("IP: %s\n", mt.IP())
}
```

## Architecture Overview

The metas package provides a centralized metadata management system:

```
┌─────────────────┐
│   Metadata      │ - Global metadata instance
└────────┬────────┘
         │
    ┌────┴────┐
    │ Options │ - Configuration options
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Environment    │ - Environment variables
    └─────────────────┘
```

**Key Components:**
- **Metadata**: Main metadata structure with application information
- **Options**: Configuration options for metadata initialization
- **KindType**: Application type (API, Job, Tool, Messaging)
- **ModeType**: Deployment mode (Local, Test, Staging, Prod)

## Usage Guide

### Creating Metadata

Create a new metadata instance:

```go
import (
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func main() {
    // Create metadata with options
    mt := metas.New(&metas.Options{
        App:    "myapp",
        Kind:   metas.KindAPI,
        Mode:   metas.ModeLocal,
        Config: "yaml://@./configs/local",
    })
    
    // Access global metadata instance
    mt := metas.Metadata()
}
```

**Validation Rules:**
- App name must match pattern: `^[a-zA-Z]+[a-zA-Z0-9_-]+[a-zA-Z0-9]+$`
- Kind must be one of: API, Job, Tool, Messaging
- Mode must be one of: Local, Test, Staging, Prod
- Config path must match pattern: `^[a-zA-Z0-9.-_]+/(local|test|staging|prod)[/]?$`

### Accessing Metadata

Access various metadata fields:

```go
mt := metas.Metadata()

// Application information
appName := mt.App()           // Application name
appKind := mt.Kind()          // Application type
appMode := mt.Mode()          // Deployment mode

// Deployment information
node := mt.Node()             // Node ID
region := mt.Region()         // Region
zone := mt.Zone()             // Zone
provider := mt.Provider()     // Service provider
ip := mt.IP()                 // Local IP address

// Configuration information
configType := mt.ConfigType() // Configuration type (yaml/etcd)
configAddr := mt.ConfigAddr() // Configuration address
configPath := mt.ConfigPath() // Configuration path

// Logging information
logDir := mt.LogDir()         // Log directory
```

### Environment Variables

The package automatically reads from environment variables:

```bash
# Set environment variables
export SONG_MODE=prod
export SONG_NODE=node-1
export SONG_REGION=us-east-1
export SONG_ZONE=us-east-1a
export SONG_PROVIDER=aws
export SONG_LOG_DIR=/var/log/myapp
```

These environment variables override the default values:
- `SONG_MODE`: Deployment mode
- `SONG_NODE`: Node identifier
- `SONG_REGION`: Deployment region
- `SONG_ZONE`: Availability zone
- `SONG_PROVIDER`: Cloud provider
- `SONG_LOG_DIR`: Log directory

### Configuration Management

Configure the metadata with different configuration sources:

```go
// YAML configuration
mt := metas.New(&metas.Options{
    App:  "myapp",
    Kind: metas.KindAPI,
    Mode: metas.ModeLocal,
    Config: "yaml://@./configs/local",
})

// Etcd configuration
mt := metas.New(&metas.Options{
    App:  "myapp",
    Kind: metas.KindAPI,
    Mode: metas.ModeProd,
    Config: "etcd://localhost:2379@config/myapp/prod",
})
```

**Configuration DSN Format:**
```
<type>://[address]@<path>
```

Where:
- `type`: Configuration type (yaml or etcd)
- `address`: Configuration server address (optional for yaml)
- `path`: Configuration file or key path

## Metadata Components

### Application Types (KindType)

The package supports different application types:

```go
// API application
metas.KindAPI      // Web API service

// Job application
metas.KindJob      // Background job processor

// Tool application
metas.KindTool     // Command-line tool

// Messaging application
metas.KindMessaging // Message processor
```

### Deployment Modes (ModeType)

The package supports different deployment modes:

```go
// Local development
metas.ModeLocal    // Local development environment

// Testing
metas.ModeTest     // Testing environment

// Staging
metas.ModeStaging  // Staging/pre-production environment

// Production
metas.ModeProd     // Production environment
```

### Mode Validation

Check the current deployment mode:

```go
mt := metas.Metadata()

if mt.Mode().IsModeLocal() {
    // Enable debug features
}

if mt.Mode().IsModeProd() {
    // Use production settings
}

if mt.Mode().IsModeTest() || mt.Mode().IsModeStaging() {
    // Use staging/test settings
}
```

## Configuration Options

### Options Structure

```go
type Options struct {
    App    string   // Application name (required)
    Kind   KindType // Application type (required)
    Mode   ModeType // Deployment mode (optional, default: local)
    Config string   // Configuration address (optional)
}
```

### Configuration Constants

```go
const (
    // Configuration DSN
    ConfigDSNDefault       = "yaml://@./configs/local"
    ConfigDSNRegexpPattern = `^(yaml|etcd)://([a-zA-Z.:0-9]*)@([a-zA-Z0-9/._-]+)$`
    
    // Configuration path
    ConfigPathRegexpPattern = `^[a-zA-Z0-9.-_]+/(local|test|staging|prod)[/]?$`
    ConfigDirDefault        = "./configs/local"
    ConfigTypeYaml          = "yaml"
    ConfigTypeEtcd          = "etcd"
    
    // Log directory
    LogDirDefault = "./logs"
    
    // App name pattern
    FlagAppRegexpPattern = `^[a-zA-Z]+[a-zA-Z0-9_-]+[a-zA-Z0-9]+$`
)
```

### Environment Variable Names

```go
const (
    EnvNameMode     = "SONG_MODE"
    EnvNameNode     = "SONG_NODE"
    EnvNameRegion   = "SONG_REGION"
    EnvNameZone     = "SONG_ZONE"
    EnvNameProvider = "SONG_PROVIDER"
    EnvNameLogDir   = "SONG_LOG_DIR"
)
```

## Examples

### Complete Example: API Application

```go
package main

import (
    "fmt"
    "log"
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize metadata
    mt := metas.New(&metas.Options{
        App:  "user-api",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
        Config: "yaml://@./configs/local",
    })
    
    // Log application information
    log.Printf("Starting %s (%s) in %s mode", 
        mt.App(), mt.Kind(), mt.Mode())
    log.Printf("Node: %s, Region: %s, Zone: %s", 
        mt.Node(), mt.Region(), mt.Zone())
    log.Printf("IP: %s, Log Dir: %s", 
        mt.IP(), mt.LogDir())
    
    // Create HTTP server
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "app":    mt.App(),
                    "kind":   mt.Kind().String(),
                    "mode":   mt.Mode().String(),
                    "node":   mt.Node(),
                    "status": "ok",
                })
            })
            
            eng.GET("/info", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "application": mt.App(),
                    "type":        mt.Kind().String(),
                    "environment": mt.Mode().String(),
                    "node":        mt.Node(),
                    "region":      mt.Region(),
                    "zone":        mt.Zone(),
                    "provider":    mt.Provider(),
                    "ip":          mt.IP(),
                    "config_type": mt.ConfigType(),
                    "config_path": mt.ConfigPath(),
                    "log_dir":     mt.LogDir(),
                })
            })
        }),
    })
    
    // Start server
    server.Serve()
}
```

### Example: Environment-Specific Configuration

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func initializeApplication() {
    mt := metas.Metadata()
    
    // Configure based on mode
    var configPath string
    switch {
    case mt.Mode().IsModeLocal():
        configPath = "yaml://@./configs/local"
        enableDebugFeatures()
        
    case mt.Mode().IsModeTest():
        configPath = "yaml://@./configs/test"
        enableTestFixtures()
        
    case mt.Mode().IsModeStaging():
        configPath = "etcd://staging-etcd:2379@config/app/staging"
        enableStagingFeatures()
        
    case mt.Mode().IsModeProd():
        configPath = "etcd://prod-etcd:2379@config/app/prod"
        enableProductionOptimizations()
    }
    
    // Load configuration
    loadConfig(configPath)
}

func enableDebugFeatures() {
    // Enable verbose logging, debug endpoints, etc.
}

func enableProductionOptimizations() {
    // Enable connection pooling, caching, etc.
}
```

### Example: Multi-Region Deployment

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func main() {
    // Initialize metadata
    mt := metas.New(&metas.Options{
        App:  "global-service",
        Kind: metas.KindAPI,
        Mode: metas.ModeProd,
        Config: "etcd://etcd.global:2379@config/global/prod",
    })
    
    // Check deployment location
    region := mt.Region()
    zone := mt.Zone()
    
    fmt.Printf("Deployed in region: %s, zone: %s\n", region, zone)
    
    // Configure region-specific settings
    switch region {
    case "us-east-1":
        configureUSEast()
    case "eu-west-1":
        configureEUWest()
    case "ap-northeast-1":
        configureAPNortheast()
    }
    
    // Configure zone-specific settings
    configureForZone(zone)
}

func configureUSEast() {
    // US East configuration
}

func configureEUWest() {
    // EU West configuration
}

func configureAPNortheast() {
    // Asia Pacific configuration
}

func configureForZone(zone string) {
    // Zone-specific configuration
}
```

### Example: Application Type Detection

```go
func initializeBasedOnKind() {
    mt := metas.Metadata()
    
    switch mt.Kind() {
    case metas.KindAPI:
        initializeAPIServer()
    case metas.KindJob:
        initializeJobProcessor()
    case metas.KindTool:
        initializeCLI()
    case metas.KindMessaging:
        initializeMessageHandler()
    }
}

func initializeAPIServer() {
    fmt.Println("Starting API server...")
    // Initialize HTTP server, routers, handlers
}

func initializeJobProcessor() {
    fmt.Println("Starting job processor...")
    // Initialize job queue, workers
}

func initializeCLI() {
    fmt.Println("Starting CLI tool...")
    // Initialize command-line interface
}

func initializeMessageHandler() {
    fmt.Println("Starting message handler...")
    // Initialize message subscribers, handlers
}
```

### Example: Health Check with Metadata

```go
func HealthCheckHandler(c *gin.Context) {
    mt := metas.Metadata()
    
    response := gin.H{
        "status":     "healthy",
        "timestamp":  time.Now().UTC(),
        "app":        mt.App(),
        "version":    getVersion(),
        "environment": mt.Mode().String(),
        "node":       mt.Node(),
        "region":     mt.Region(),
        "uptime":     getUptime(),
    }
    
    // Add provider-specific information
    if mt.Provider() != "" {
        response["provider"] = mt.Provider()
    }
    
    c.JSON(200, response)
}

func getVersion() string {
    // Return application version from build info
    return "1.0.0"
}

func getUptime() string {
    // Calculate and return uptime
    return time.Since(startTime).String()
}
```

## Best Practices

1. **Initialize Early**: Initialize metadata at the very beginning of your application startup.

2. **Use Singleton Pattern**: Use `metas.Metadata()` to access the global metadata instance instead of creating multiple instances.

3. **Validate Configuration**: Always validate configuration paths and application names.

4. **Use Environment Variables**: Leverage environment variables for deployment-specific settings.

5. **Separate Environments**: Use different modes (local, test, staging, prod) to separate environments.

6. **Include Metadata in Logs**: Include application metadata in log entries for better observability.

7. **Use in Health Checks**: Expose metadata in health check endpoints for monitoring.

8. **Respect Mode Settings**: Adjust behavior based on the deployment mode (e.g., enable debug in local, optimize for prod).

9. **Document Configuration**: Document all configuration options and environment variables.

10. **Use Consistent Naming**: Follow the naming pattern for applications across your organization.

11. **Secure Sensitive Information**: Don't expose sensitive metadata (like credentials) in logs or responses.

12. **Monitor Configuration Changes**: Track configuration changes in production environments.

## Additional Resources

- [Song Framework Documentation](../../README.md)
- [Viper Configuration Module](../vipers/README.md)
