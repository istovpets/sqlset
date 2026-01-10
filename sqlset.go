// Package sqlset is a way to store SQL queries separated from the go code.
// Query sets are stored in the .sql files, every filename without extension is an SQL set ID.
// Every file contains queries, marked with query IDs using special syntax,
// see `testdata/valid/*.sql` files for examples.
// Also file may contain JSON-encoded query set metadata with name and description.
package sqlset

import (
	"fmt"
	"sort"
	"strings"
)

// SQLQueriesProvider is the interface for getting SQL queries.
type SQLQueriesProvider interface {
	// Get returns a query by set ID and query ID.
	// If the set or query is not found, it returns an error.
	Get(ids ...string) (string, error)
	// MustGet returns a query by set ID and query ID.
	// It panics if the set or query is not found.
	MustGet(ids ...string) string
}

// SQLSetsProvider is the interface for getting information about query sets.
type SQLSetsProvider interface {
	// GetSetsMetas returns metadata for all registered query sets.
	GetSetsMetas() []QuerySetMeta
	// GetQueryIDs returns a slice of all query IDs.
	GetQueryIDs(setID string) ([]string, error)
}

// SQLSet is a container for multiple query sets, organized by set ID.
// It provides methods to access SQL queries and metadata.
// Use New to create a new instance.
type SQLSet struct {
	sets map[string]QuerySet
}

// Get returns an SQL query by its identifiers.
//
// Supported forms:
//
//   - Get(setID, queryID)
//     Returns the query identified by queryID from the query set setID.
//
//   - Get("setID.queryID")
//     Equivalent to Get(setID, queryID).
//
//   - Get(queryID)
//     Returns the query identified by queryID from the only available query set.
//     If there is more than one query set, an error is returned.
//
// An error is returned if:
//   - the number of arguments is invalid,
//   - any identifier is empty,
//   - the query set or query cannot be found.
func (s *SQLSet) Get(ids ...string) (string, error) {
	for i, id := range ids {
		if id == "" {
			return "", fmt.Errorf("empty argument: %d", i)
		}
	}

	l := len(ids)
	if l == 0 || l > 2 {
		return "", fmt.Errorf("invalid number of arguments: %d", l)
	}

	if l == 1 {
		left, right, ok := strings.Cut(ids[0], ".")
		if ok {
			ids = []string{left, right}
		}
	}

	return s.findQuery(ids...)
}

// MustGet is like Get but panics if the query set or query is not found.
// This is useful for cases where the query is expected to exist and its absence is a critical error.
func (s *SQLSet) MustGet(ids ...string) string {
	q, err := s.Get(ids...)
	if err != nil {
		panic(err)
	}

	return q
}

// GetSetsMetas returns a slice of metadata for all the query sets loaded.
// The order of the returned slice is not guaranteed.
func (s *SQLSet) GetSetsMetas() []QuerySetMeta {
	metas := make([]QuerySetMeta, 0, len(s.sets))

	for _, qs := range s.sets {
		metas = append(metas, qs.GetMeta())
	}

	return metas
}

// GetQueryIDs returns a sorted slice of all query IDs within a specific query set.
func (s *SQLSet) GetQueryIDs(setID string) ([]string, error) {
	if s.sets == nil {
		return nil, fmt.Errorf("%s: %w", setID, ErrQuerySetNotFound)
	}

	qs, ok := s.sets[setID]
	if !ok {
		return nil, fmt.Errorf("%s: %w", setID, ErrQuerySetNotFound)
	}

	if qs.queries == nil {
		return []string{}, nil
	}

	ids := make([]string, 0, len(qs.queries))
	for id := range qs.queries {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	return ids, nil
}

func (s *SQLSet) findQuery(ids ...string) (string, error) {
	if s.sets == nil {
		return "", ErrQuerySetsEmpty
	}

	var (
		qs      QuerySet
		queryID string
		ok      bool
	)

	if len(ids) == 1 {
		if len(s.sets) > 1 {
			return "", fmt.Errorf("query set not specified")
		}

		queryID = ids[0]

		for _, v := range s.sets {
			qs = v
			break
		}
	} else if len(ids) == 2 {
		queryID = ids[1]

		qs, ok = s.sets[ids[0]]
		if !ok {
			return "", fmt.Errorf("%s: %w", ids[0], ErrQuerySetNotFound)
		}
	} else {
		return "", fmt.Errorf("invalid number of arguments: %d", len(ids))
	}

	q, err := qs.findQuery(queryID)
	if err != nil {
		return "", err
	}

	return q, nil
}

func (s *SQLSet) registerQuerySet(setID string, qs QuerySet) {
	if s.sets == nil {
		s.sets = make(map[string]QuerySet)
	}

	s.sets[setID] = qs
}

// QuerySet represents a single set of queries, usually from a single .sql file.
type QuerySet struct {
	meta    QuerySetMeta
	queries map[string]string
}

// GetMeta returns the metadata associated with the query set.
func (qs *QuerySet) GetMeta() QuerySetMeta {
	return qs.meta
}

func (qs *QuerySet) registerQuery(id string, query string) {
	if qs.queries == nil {
		qs.queries = make(map[string]string)
	}

	qs.queries[id] = query
}

func (qs *QuerySet) findQuery(id string) (string, error) {
	if qs.queries == nil {
		return "", fmt.Errorf("%s: %w", qs.meta.ID, ErrQuerySetEmpty)
	}

	q, ok := qs.queries[id]
	if !ok {
		return "", fmt.Errorf("%s: %w", id, ErrQueryNotFound)
	}

	return q, nil
}

// QuerySetMeta holds the metadata for a query set.
type QuerySetMeta struct {
	// ID is the unique identifier for the set, derived from the filename.
	ID string `json:"id"`
	// Name is a human-readable name for the query set, from the metadata block.
	Name string `json:"name"`
	// Description provides more details about the query set, from the metadata block.
	Description string `json:"description,omitempty"`
}
