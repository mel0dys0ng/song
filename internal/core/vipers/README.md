# Vipers

A powerful configuration management system built on top of Viper, providing support for multiple configuration sources (YAML, JSON, TOML, etcd, Consul), real-time configuration updates, and type-safe configuration access with default values.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
  - [Creating Configuration](#creating-configuration)
  - [Reading Configuration Values](#reading-configuration-values)
  - [Configuration Providers](#configuration-providers)
  - [Configuration Updates](#configuration-updates)
- [Configuration Options](#configuration-options)
- [Type-Safe Access Methods](#type-safe-access-methods)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Multiple Configuration Sources**: Support for YAML, JSON, TOML, etcd, and Consul
- **Real-Time Updates**: Automatic configuration reload on changes
- **Type-Safe Access**: Strongly-typed getter methods with default values
- **Configuration Validation**: Built-in validation and error handling
- **Hierarchical Configuration**: Support for nested configuration structures
- **Environment Variable Integration**: Automatic environment variable binding
- **Remote Configuration**: Support for distributed configuration management
- **Configuration Change Hooks**: Callback functions for configuration changes
- **Default Values**: Graceful fallback to default values
- **Unified Interface**: Consistent API across all configuration sources

## Installation

The vipers package requires the following dependencies:

```bash
go get github.com/spf13/viper
go get github.com/fsnotify/fsnotify
go get go.etcd.io/etcd/client/v3
go get github.com/hashicorp/consul/api
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/vipers"
)

func main() {
    // Create configuration from YAML file
    config, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/app.yaml"),
    )
    if err != nil {
        panic(err)
    }
    
    // Read configuration values
    port := config.GetInt("server.port", 8080)
    host := config.GetString("server.host", "localhost")
    debug := config.GetBool("server.debug", false)
    
    fmt.Printf("Server: %s:%d (debug: %v)\n", host, port, debug)
}
```

## Architecture Overview

The vipers package provides a unified configuration management system:

```
┌─────────────────┐
│    Config       │ - Main configuration structure
└────────┬────────┘
         │
    ┌────┴────┐
    │ Viper   │ - Underlying Viper instance
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Provider       │ - Configuration source provider
    └─────────────────┘
         │
    ┌────┴────┬────────┬────────┬──────────┐
    │  YAML   │  JSON  │  TOML  │  Etcd    │ Consul │
    └─────────┴────────┴────────┴──────────┴────────┘
```

**Key Components:**
- **Config**: Main configuration structure with Viper integration
- **Provider**: Abstract interface for configuration sources
- **Options**: Configuration options and settings
- **ProviderInterface**: Interface for different configuration providers

## Usage Guide

### Creating Configuration

Create a new configuration instance:

```go
import (
    "github.com/mel0dys0ng/song/internal/core/vipers"
)

// YAML configuration
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderYaml),
    vipers.OnPath("./configs/app.yaml"),
)

// JSON configuration
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderJson),
    vipers.OnPath("./configs/app.json"),
)

// TOML configuration
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderToml),
    vipers.OnPath("./configs/app.toml"),
)

// Etcd configuration
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderEtcd),
    vipers.OnEndpoints("localhost:2379"),
    vipers.OnPath("config/myapp/prod"),
)

// Consul configuration
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderConsul),
    vipers.OnEndpoints("localhost:8500"),
    vipers.OnPath("config/myapp/prod"),
)
```

### Reading Configuration Values

Access configuration values with type-safe methods:

```go
// String values
host := config.GetString("server.host", "localhost")
name := config.GetString("app.name", "MyApp")

// Integer values
port := config.GetInt("server.port", 8080)
maxConn := config.GetInt("database.max_connections", 100)

// Boolean values
debug := config.GetBool("server.debug", false)
enabled := config.GetBool("features.new_ui", true)

// Float values
timeout := config.GetFloat64("server.timeout", 30.0)
ratio := config.GetFloat64("sampling.rate", 0.1)

// Duration values
readTimeout := config.GetDuration("server.read_timeout", 30*time.Second)
retryDelay := config.GetDuration("retry.delay", 5*time.Second)

// Slice values
hosts := config.GetStringSlice("server.hosts", []string{"localhost"})
ports := config.GetIntSlice("server.ports", []int{8080, 8081})

// Map values
labels := config.GetStringMapString("labels", map[string]string{})
metadata := config.GetStringMap("metadata", map[string]any{})
```

**All Getter Methods:**
- `Get(key, defaultValue)` - Generic getter
- `GetString(key, defaultValue)` - String value
- `GetInt(key, defaultValue)` - Integer value
- `GetInt32(key, defaultValue)` - 32-bit integer
- `GetInt64(key, defaultValue)` - 64-bit integer
- `GetUint(key, defaultValue)` - Unsigned integer
- `GetUint32(key, defaultValue)` - 32-bit unsigned integer
- `GetUint64(key, defaultValue)` - 64-bit unsigned integer
- `GetBool(key, defaultValue)` - Boolean value
- `GetFloat64(key, defaultValue)` - Float64 value
- `GetDuration(key, defaultValue)` - Duration value
- `GetTime(key, defaultValue)` - Time value
- `GetStringSlice(key, defaultValue)` - String slice
- `GetIntSlice(key, defaultValue)` - Integer slice
- `GetStringMap(key, defaultValue)` - Generic map
- `GetStringMapString(key, defaultValue)` - String map
- `GetStringMapStringSlice(key, defaultValue)` - Map of string slices
- `GetSizeInBytes(key, defaultValue)` - Size in bytes

### Configuration Providers

#### YAML Provider

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderYaml),
    vipers.OnPath("./configs/app.yaml"),
)
```

**Example YAML file:**
```yaml
server:
  host: localhost
  port: 8080
  debug: true
  timeout: 30s

database:
  host: localhost
  port: 5432
  name: myapp
  max_connections: 100
  
features:
  new_ui: true
  beta_features:
    - feature1
    - feature2
```

#### JSON Provider

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderJson),
    vipers.OnPath("./configs/app.json"),
)
```

**Example JSON file:**
```json
{
  "server": {
    "host": "localhost",
    "port": 8080,
    "debug": true,
    "timeout": "30s"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "myapp",
    "max_connections": 100
  }
}
```

#### TOML Provider

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderToml),
    vipers.OnPath("./configs/app.toml"),
)
```

**Example TOML file:**
```toml
[server]
host = "localhost"
port = 8080
debug = true
timeout = "30s"

[database]
host = "localhost"
port = 5432
name = "myapp"
max_connections = 100
```

#### Etcd Provider

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderEtcd),
    vipers.OnEndpoints("localhost:2379"),
    vipers.OnPath("config/myapp/prod"),
    vipers.OnUsername("etcd_user"),
    vipers.OnPassword("etcd_password"),
)
```

#### Consul Provider

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderConsul),
    vipers.OnEndpoints("localhost:8500"),
    vipers.OnPath("config/myapp/prod"),
    vipers.OnToken("consul_token"),
)
```

### Configuration Updates

Watch for configuration changes:

```go
// Set up configuration change handler
config.OnConfigChange(func(event fsnotify.Event, options *vipers.Options) {
    fmt.Println("Configuration changed!")
    fmt.Printf("Event: %v\n", event)
    
    // Reload configuration or take action
    newPort := config.GetInt("server.port", 8080)
    log.Printf("New port: %d", newPort)
})
```

## Configuration Options

### Option Functions

```go
// Provider configuration
vipers.OnProvider(vipers.ConfigProviderYaml)  // Set configuration provider
vipers.OnPath("./configs/app.yaml")           // Set configuration path
vipers.OnEndpoints("localhost:2379")          // Set remote endpoints

// Authentication
vipers.OnUsername("user")                     // Set username
vipers.OnPassword("password")                 // Set password
vipers.OnToken("token")                       // Set authentication token

// Change monitoring
vipers.OnChangeConfig(func(event, options) {  // Set change callback
    // Handle configuration change
})
```

### Provider Constants

```go
const (
    ConfigProviderJson   = "json"   // JSON provider
    ConfigProviderYaml   = "yaml"   // YAML provider
    ConfigProviderToml   = "toml"   // TOML provider
    ConfigProviderEtcd   = "etcd"   // Etcd provider
    ConfigProviderConsul = "consul" // Consul provider
)
```

### Default Options

```go
opts := vipers.DefaultOptions()
// Returns:
// - Provider: ConfigProviderYaml
// - Path: "./configs/local"
// - Endpoints: ""
// - Username: ""
// - Password: ""
// - Token: ""
// - OnChangeConfig: nil
```

## Type-Safe Access Methods

All getter methods follow the same pattern:

```go
func (c *Config) Get<Type>(key string, defaultValue <Type>) <Type>
```

**Behavior:**
- If the key exists and can be converted to the target type, return the value
- If the key doesn't exist or conversion fails, return the default value
- Never panics - always returns a valid value

**Example:**
```go
// Will return 8080 if key doesn't exist
port := config.GetInt("server.port", 8080)

// Will return "localhost" if key doesn't exist
host := config.GetString("server.host", "localhost")

// Will return false if key doesn't exist
debug := config.GetBool("server.debug", false)
```

## Examples

### Complete Example: Application Configuration

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/mel0dys0ng/song/internal/core/vipers"
)

type ServerConfig struct {
    Host         string        `mapstructure:"host"`
    Port         int           `mapstructure:"port"`
    Debug        bool          `mapstructure:"debug"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Host           string `mapstructure:"host"`
    Port           int    `mapstructure:"port"`
    Name           string `mapstructure:"name"`
    User           string `mapstructure:"user"`
    Password       string `mapstructure:"password"`
    MaxConnections int    `mapstructure:"max_connections"`
}

type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}

func main() {
    // Load configuration
    config, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/app.yaml"),
    )
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Read individual values
    host := config.GetString("server.host", "localhost")
    port := config.GetInt("server.port", 8080)
    debug := config.GetBool("server.debug", false)
    
    log.Printf("Server: %s:%d (debug: %v)", host, port, debug)
    
    // Unmarshal into struct
    var appConfig AppConfig
    if err := config.Unmarshal(&appConfig); err != nil {
        log.Fatalf("Failed to unmarshal configuration: %v", err)
    }
    
    log.Printf("Server Config: %+v", appConfig.Server)
    log.Printf("Database Config: %+v", appConfig.Database)
    
    // Set up configuration change monitoring
    config.OnConfigChange(func(event fsnotify.Event, options *vipers.Options) {
        log.Println("Configuration file changed, reloading...")
        
        // Re-read configuration
        newPort := config.GetInt("server.port", 8080)
        log.Printf("New port: %d", newPort)
    })
    
    // Keep application running to receive configuration changes
    select {}
}
```

### Example: Environment-Specific Configuration

```go
func loadEnvironmentConfig(env string) (*vipers.Config, error) {
    var path string
    var provider string
    
    switch env {
    case "local":
        path = "./configs/local.yaml"
        provider = vipers.ConfigProviderYaml
        
    case "test":
        path = "./configs/test.yaml"
        provider = vipers.ConfigProviderYaml
        
    case "staging":
        path = "config/app/staging"
        provider = vipers.ConfigProviderEtcd
        
    case "prod":
        path = "config/app/prod"
        provider = vipers.ConfigProviderEtcd
        
    default:
        return nil, fmt.Errorf("unknown environment: %s", env)
    }
    
    // Build configuration options
    opts := []vipers.Option{
        vipers.OnProvider(provider),
        vipers.OnPath(path),
    }
    
    // Add etcd endpoints for remote environments
    if provider == vipers.ConfigProviderEtcd {
        opts = append(opts, 
            vipers.OnEndpoints("etcd.cluster:2379"),
            vipers.OnUsername("etcd_user"),
            vipers.OnPassword("etcd_password"),
        )
    }
    
    return vipers.New(opts...)
}
```

### Example: Feature Flags

```go
type FeatureFlags struct {
    NewUI         bool     `mapstructure:"new_ui"`
    BetaFeatures  []string `mapstructure:"beta_features"`
    EnabledRegions []string `mapstructure:"enabled_regions"`
    RolloutRate   float64  `mapstructure:"rollout_rate"`
}

func loadFeatureFlags(config *vipers.Config) (*FeatureFlags, error) {
    var flags FeatureFlags
    
    // Use type-safe getters with defaults
    flags.NewUI = config.GetBool("features.new_ui", false)
    flags.BetaFeatures = config.GetStringSlice("features.beta_features", []string{})
    flags.EnabledRegions = config.GetStringSlice("features.enabled_regions", []string{})
    flags.RolloutRate = config.GetFloat64("features.rollout_rate", 0.0)
    
    // Or unmarshal entire section
    if err := config.UnmarshalKey("features", &flags); err != nil {
        return nil, err
    }
    
    return &flags, nil
}

func isEnabled(feature string, flags *FeatureFlags) bool {
    switch feature {
    case "new_ui":
        return flags.NewUI
    case "beta_feature1", "beta_feature2":
        return contains(flags.BetaFeatures, feature)
    default:
        return false
    }
}

func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}
```

### Example: Database Configuration

```go
func loadDatabaseConfig(config *vipers.Config) (*DatabaseConfig, error) {
    var dbConfig DatabaseConfig
    
    // Read configuration with defaults
    dbConfig.Host = config.GetString("database.host", "localhost")
    dbConfig.Port = config.GetInt("database.port", 5432)
    dbConfig.Name = config.GetString("database.name", "myapp")
    dbConfig.User = config.GetString("database.user", "postgres")
    dbConfig.Password = config.GetString("database.password", "")
    dbConfig.MaxConnections = config.GetInt("database.max_connections", 100)
    
    // Validate required fields
    if dbConfig.Host == "" {
        return nil, fmt.Errorf("database host is required")
    }
    if dbConfig.Name == "" {
        return nil, fmt.Errorf("database name is required")
    }
    
    return &dbConfig, nil
}

func createDatabaseConnection(dbConfig *DatabaseConfig) (*sql.DB, error) {
    // Build connection string
    connStr := fmt.Sprintf(
        "host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
        dbConfig.Host,
        dbConfig.Port,
        dbConfig.Name,
        dbConfig.User,
        dbConfig.Password,
    )
    
    // Create connection pool
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(dbConfig.MaxConnections)
    db.SetMaxIdleConns(dbConfig.MaxConnections / 2)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return db, nil
}
```

### Example: Multi-Source Configuration

```go
func loadMultiSourceConfig() (*vipers.Config, error) {
    // Load base configuration from file
    baseConfig, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/base.yaml"),
    )
    if err != nil {
        return nil, err
    }
    
    // Override with environment-specific configuration
    env := os.Getenv("APP_ENV")
    if env != "" {
        envConfig, err := vipers.New(
            vipers.OnProvider(vipers.ConfigProviderYaml),
            vipers.OnPath(fmt.Sprintf("./configs/%s.yaml", env)),
        )
        if err != nil {
            return nil, err
        }
        
        // Merge configurations (envConfig takes precedence)
        mergeConfigs(baseConfig, envConfig)
    }
    
    // Override with environment variables
    bindEnvironmentVariables(baseConfig)
    
    return baseConfig, nil
}

func mergeConfigs(base, override *vipers.Config) {
    // Get all keys from override config
    for _, key := range override.AllKeys() {
        value := override.Get(key, nil)
        if value != nil {
            // Set in base config (override takes precedence)
            base.Viper.Set(key, value)
        }
    }
}

func bindEnvironmentVariables(config *vipers.Config) {
    // Automatically bind environment variables
    config.Viper.AutomaticEnv()
    
    // Example: SONG_SERVER_PORT overrides server.port
    config.Viper.SetEnvPrefix("SONG")
    config.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
```

### Example: Configuration Validation

```go
type ValidatableConfig interface {
    Validate() error
}

func loadAndValidateConfig(config *vipers.Config) (*AppConfig, error) {
    var appConfig AppConfig
    
    // Unmarshal configuration
    if err := config.Unmarshal(&appConfig); err != nil {
        return nil, err
    }
    
    // Validate configuration
    if err := validateAppConfig(&appConfig); err != nil {
        return nil, err
    }
    
    return &appConfig, nil
}

func validateAppConfig(config *AppConfig) error {
    // Validate server configuration
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    
    // Validate database configuration
    if config.Database.MaxConnections < 1 {
        return fmt.Errorf("database max connections must be positive")
    }
    
    if config.Database.MaxConnections > 1000 {
        return fmt.Errorf("database max connections cannot exceed 1000")
    }
    
    return nil
}
```

## Best Practices

1. **Use Type-Safe Getters**: Always use type-safe getter methods with default values to avoid panics.

2. **Validate Configuration**: Validate critical configuration values before using them.

3. **Use Struct Tags**: Use `mapstructure` tags when unmarshaling into structs for better type safety.

4. **Separate Environments**: Use different configuration files/sources for different environments.

5. **Monitor Changes**: Set up configuration change handlers for dynamic configuration updates.

6. **Secure Sensitive Data**: Never store sensitive data (passwords, tokens) in plain text configuration files.

7. **Use Environment Variables**: Use environment variables for environment-specific overrides.

8. **Document Configuration**: Document all configuration options, their types, and default values.

9. **Provide Sensible Defaults**: Always provide sensible default values for optional configuration.

10. **Test Configuration**: Test configuration loading and validation in your test suite.

11. **Fail Fast**: Fail early if required configuration is missing or invalid.

12. **Use Hierarchical Configuration**: Organize configuration in a logical hierarchy (e.g., `server.port`, `database.host`).

## Additional Resources

- [Song Framework Documentation](../../README.md)
- [Metas Metadata Module](../metas/README.md)
- [Viper Documentation](https://github.com/spf13/viper)
