# H3S Utility Packages

This directory contains utility packages that provide common functionality used throughout the H3S codebase.

## Packages

### `cli`

Provides utilities for creating and managing CLI commands with a standardized approach.

- `NewCommand` - Creates a new cobra command with standardized configuration
- `CommandConfig` - Configuration for creating a command
- `Flag` - Represents a command flag

### `cloud`

Provides utilities for interacting with cloud providers.

- `ResourceClient` - Generic interface for cloud resource operations
- `NewClient` - Creates a new cloud resource client
- `Get`, `Create`, `Delete`, `GetOrCreate` - Standard resource operations

### `naming`

Provides utilities for generating consistent resource names.

- `ResourceName` - Generates a consistent name for a resource based on the cluster context
- `FormatName` - Joins name components with a consistent separator

### `resource`

Provides utilities for managing resources with consistent patterns.

- `ResourceManager` - Generic interface for resource operations
- `NewManager` - Creates a new resource manager
- `Get`, `Create`, `Delete`, `GetOrCreate` - Standard resource operations

### `validation`

Provides utilities for validating input data.

- Common validation functions: `StringNotEmpty`, `StringLength`, `StringMatches`
- Domain-specific validation: `Name`, `Email`, `Domain`, `IP`, `Number`, etc.
- Common validation errors and regular expressions

## Usage

### CLI Commands

```go
// Create a command with the standardized approach
cmd := cli.NewCommand(cli.CommandConfig{
    Use:   "example",
    Short: "Example command",
    Long:  "Example command long description",
    RunE:  runExample,
    Args:  cobra.NoArgs,
    Flags: []cli.Flag{
        {
            Name:         "flag",
            Shorthand:    "f",
            Value:        "",
            DefaultValue: "default",
            Usage:        "Example flag",
            Required:     false,
        },
    },
})
```

### Resource Management

```go
// Create a resource manager
manager := resource.NewManager[*MyResource](ctx, "resource-type", "resource-name")

// Get a resource
resource, err := manager.Get(func() (*MyResource, error) {
    // Implementation to get the resource
    return getResource()
})

// Create a resource
resource, err := manager.Create(func() (*MyResource, error) {
    // Implementation to create the resource
    return createResource()
})

// Delete a resource
err := manager.Delete(func() error {
    // Implementation to delete the resource
    return deleteResource()
})

// Get or create a resource
resource, err := manager.GetOrCreate(
    func() (*MyResource, error) {
        // Implementation to get the resource
        return getResource()
    },
    func() (*MyResource, error) {
        // Implementation to create the resource
        return createResource()
    },
)
```

### Validation

```go
// Validate a name
if err := validation.Name(name); err != nil {
    return err
}

// Validate an email
if err := validation.Email(email); err != nil {
    return err
}

// Validate a number in range
if err := validation.NumberInRange(value, 1, 100); err != nil {
    return err
}
```

### Naming

```go
// Generate a resource name
name := naming.ResourceName(ctx, "component1", "component2")

// Format a name
name := naming.FormatName("prefix", "component1", "component2")
```
