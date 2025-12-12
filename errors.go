package sqlset

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrQuerySetNotFound   = fmt.Errorf("query set %w", ErrNotFound)
	ErrQueryNotFound      = fmt.Errorf("query %w", ErrNotFound)
	ErrInvalidSyntax      = errors.New("invalid SQLSetList syntax")
	ErrMaxLineLenExceeded = errors.New("line too long, possible line corruption")
)
