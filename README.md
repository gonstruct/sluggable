# Sluggable

[![Go Version](https://img.shields.io/github/go-mod/go-version/gonstruct/sluggable)](https://golang.org/doc/devel/release.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/gonstruct/sluggable)](https://goreportcard.com/report/github.com/gonstruct/sluggable)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful and flexible Go library for generating unique database-safe slugs with collision detection and resolution.

## Features

- **Unique Slug Generation**: Automatically generates unique slugs by checking existing database entries
- **Collision Resolution**: Intelligently handles duplicate slugs by appending numeric suffixes
- **Database Agnostic**: Works with any SQL database through Go's `database/sql` interface
- **Flexible Configuration**: Customizable separators, column names, and slug generation methods
- **Soft Delete Support**: Built-in support for soft delete patterns
- **Custom Queries**: Add custom WHERE clauses for advanced filtering
- **Global & Instance Configuration**: Use global defaults or create custom instances

## Installation

```bash
go get github.com/gonstruct/sluggable
```

## Quick Start

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    "github.com/gonstruct/sluggable"
    _ "github.com/lib/pq" // or your preferred database driver
)

func main() {
    db, err := sql.Open("postgres", "your-connection-string")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Generate a unique slug
    slug, err := sluggable.Generate(db, "Hello World", 
        sluggable.WithTableName("articles"),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(slug) // Output: "hello-world" or "hello-world-2" if duplicate exists
}
```

## Usage

### Basic Usage

The simplest way to use sluggable is with the global `Generate` function:

```go
slug, err := sluggable.Generate(db, "My Article Title", 
    sluggable.WithTableName("articles"),
)
```

### Global Configuration

Configure global defaults that apply to all slug generation:

```go
func init() {
    sluggable.Configure(
        sluggable.WithColumnName("slug"),        // Default: "slug"
        sluggable.WithSeperator("-"),            // Default: "-"
        sluggable.WithFirstUniqueSuffix(2),      // Default: 2
        sluggable.WithDeleted("deleted_at"),     // For soft deletes
    )
}
```

### Custom Instance

Create a custom sluggable instance with specific configuration:

```go
mySlugger := sluggable.New(
    sluggable.WithSeperator("_"),
    sluggable.WithFirstUniqueSuffix(1),
)

slug, err := mySlugger.Generate(db, "Article Title",
    sluggable.WithTableName("articles"),
)
```

### Advanced Options

#### Custom Slug Generation Method

```go
import "strings"

slug, err := sluggable.Generate(db, "Article Title",
    sluggable.WithTableName("articles"),
    sluggable.WithMethod(func(value, separator string) string {
        return strings.ToLower(strings.ReplaceAll(value, " ", separator))
    }),
)
```

#### Updating Existing Records

When updating an existing record, provide the identifier to avoid unnecessary suffix increments:

```go
slug, err := sluggable.Generate(db, "Updated Article Title",
    sluggable.WithTableName("articles"),
    sluggable.WithIdentifier("123"), // ID of the record being updated
)
```

#### Custom WHERE Clauses

Add additional filtering conditions:

```go
slug, err := sluggable.Generate(db, "Article Title",
    sluggable.WithTableName("articles"),
    sluggable.WithQuery(`"user_id" = $1`, userID),
    sluggable.WithDeleted("deleted_at"), // Exclude soft-deleted records
)
```

#### Soft Delete Support

Automatically exclude soft-deleted records:

```go
slug, err := sluggable.Generate(db, "Article Title",
    sluggable.WithTableName("articles"),
    sluggable.WithDeleted("deleted_at"), // NULL values will be included
)
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithTableName(string)` | Database table name (required) | `""` |
| `WithColumnName(string)` | Column name for slugs | `"slug"` |
| `WithSeperator(string)` | Separator for words and suffixes | `"-"` |
| `WithMethod(func)` | Custom slug generation function | Uses `github.com/gosimple/slug` |
| `WithFirstUniqueSuffix(int)` | Starting number for duplicate resolution | `2` |
| `WithIdentifier(string)` | ID of record being updated | `""` |
| `WithDeleted(string)` | Soft delete column name | N/A |
| `WithQuery(string, ...interface{})` | Custom WHERE clause | N/A |

## How It Works

1. **Generate Base Slug**: Converts input text to a URL-safe slug
2. **Check Database**: Queries for existing slugs matching the pattern
3. **Resolve Conflicts**: If duplicates exist, appends numeric suffix
4. **Return Unique Slug**: Guarantees uniqueness within the specified constraints

### Collision Resolution Example

Given the input "Hello World" and existing slugs:
- `hello-world` (exists)
- `hello-world-2` (exists)
- `hello-world-3` (exists)

The function will return `hello-world-4`.

## Database Requirements

Your table must have:
- An `id` column (any type that can be scanned into a string)
- A slug column (default name: `slug`)

Example table structure:

```sql
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Error Handling

The library returns descriptive errors for common issues:

```go
slug, err := sluggable.Generate(db, "Title")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "table name cannot be empty"):
        // Handle missing table name
    case strings.Contains(err.Error(), "failed to query"):
        // Handle database errors
    default:
        // Handle other errors
    }
}
```

## Performance Considerations

- Uses prepared statements internally for better performance
- Single database query per slug generation
- Efficient numeric suffix detection using string operations
- Minimal memory allocation for large result sets

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Development

```bash
# Clone the repository
git clone https://github.com/gonstruct/sluggable.git
cd sluggable

# Run tests
go test -v ./...

# Run linter
golangci-lint run

# Run linter with auto-fix
golangci-lint run --fix
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [gosimple/slug](https://github.com/gosimple/slug) for the default slug generation algorithm
- Inspired by Laravel's Eloquent slug generation patterns