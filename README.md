# SQLSet

[![Go Reference](https://pkg.go.dev/badge/github.com/stoi/sqlset.svg)](https://pkg.go.dev/github.com/stoi/sqlset)

SQLSet is a simple Go library that provides a convenient way to manage and access SQL queries stored in `.sql` files. It allows you to separate your SQL code from your Go code, making it cleaner and more maintainable.

## Features

- **Decouple SQL from Go code**: Keep your SQL queries in separate `.sql` files.
- **Easy to use**: A simple API to get your queries.
- **Flexible**: Works with any `fs.FS`, including `embed.FS` for bundling queries with your application.
- **Query Metadata**: Associate names and descriptions with your query sets.
- **Organized**: Structure your queries into logical sets.

## Installation

```bash
go get github.com/stoi/sqlset
```

## Usage

1.  **Create your SQL files**.

Create a directory (e.g., `queries`) and add your `.sql` files. Each file represents a "query set". The name of the file (without the `.sql` extension) becomes the query set ID.

Inside each file, define your queries using a special `--- meta` comment for metadata and `--- query` comments to mark the beginning of each query.

**`queries/users.sql`**
```sql
--- meta
{
    "name": "User Queries",
    "description": "A set of queries for user management."
}
---

--- query: GetUserByID
SELECT id, name, email FROM users WHERE id = ?;

--- query: CreateUser
INSERT INTO users (name, email) VALUES (?, ?);
```

2.  **Embed and load the queries in your Go application**.

Use Go's `embed` package to bundle the SQL files directly into your application binary.

```go
package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/stoi/sqlset"
)

//go:embed queries
var queriesFS embed.FS

func main() {
	// Create a new SQLSet from the embedded filesystem.
	// We pass "queries" as the subdirectory to look into.
	sqlSet, err := sqlset.New(queriesFS)
	if err != nil {
		log.Fatalf("Failed to create SQL set: %v", err)
	}

	// Get a specific query
	query, err := sqlSet.GetQuery("users", "GetUserByID")
	if err != nil {
		log.Fatalf("Failed to get query: %v", err)
	}
	fmt.Println("GetUserByID query:", query)

	// Or, panic if the query is not found
	query = sqlSet.MustGetQuery("users", "CreateUser")
	fmt.Println("CreateUser query:", query)

    // You can also retrieve metadata for all query sets
    metas := sqlSet.GetAllMetas()
    for _, meta := range metas {
        fmt.Printf("Set ID: %s, Name: %s, Description: %s\n", meta.ID, meta.Name, meta.Description)
    }
}
```

### File Format Specification

-   **Metadata Block (Optional)**:
    -   Starts with `--- meta`.
    -   Followed by a JSON object containing `name` (string) and `description` (string, optional).
    -   There can be only one metadata block per file.

-   **Query Block (Required)**:
    -   Starts with `--- query: <query_id>`, where `<query_id>` is the unique identifier for the query within the file.
    -   The SQL statement follows on the next lines.
    -   All text until the next `--- query:` block or the end of the file is considered part of the query.

## Contributing

Contributions are welcome! If you find a bug or have a feature request, please open an issue. If you want to contribute code, please open a pull request.

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature`).
3.  Make your changes.
4.  Commit your changes (`git commit -am 'Add some feature'`).
5.  Push to the branch (`git push origin feature/your-feature`).
6.  Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
