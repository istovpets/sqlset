package sqlset

import (
	"errors"
	"fmt"
)

var (
	// ErrEmpty is the base error for when an object is empty.
	ErrEmpty = errors.New("empty")
	// ErrQuerySetsEmpty indicates that a query sets is empty.
	ErrQuerySetsEmpty = fmt.Errorf("query sets %w", ErrEmpty)
	// ErrQuerySetEmpty indicates that a query sets is empty.
	ErrQuerySetEmpty = fmt.Errorf("query set %w", ErrEmpty)
	// ErrNotFound is the base error for when an item is not found.
	ErrNotFound = errors.New("not found")
	// ErrQuerySetNotFound indicates that a specific query set was not found.
	ErrQuerySetNotFound = fmt.Errorf("query set %w", ErrNotFound)
	// ErrQueryNotFound indicates that a specific query was not found within a set.
	ErrQueryNotFound = fmt.Errorf("query %w", ErrNotFound)
	// ErrInvalidSyntax is returned when the parser encounters a syntax error in a .sql file.
	ErrInvalidSyntax = errors.New("invalid SQLSetList syntax")
	// ErrMaxLineLenExceeded is returned when a line in a .sql file is too long,
	// which may indicate a corrupted file.
	ErrMaxLineLenExceeded = errors.New("line too long, possible line corruption")
)
