// Package sqlset is a way to store SQL queries separated from the go code.
// Query sets are stored in the .sql files, every filename without extension is an SQL set ID.
// Every file contains queries, marked with query IDs using special syntax,
// see `testdata/valid/*.sql` files for examples.
// Also file may contain JSON-encoded query set metadata with name and description.
package sqlset

import "fmt"

type SQLQueriesProvider interface {
	GetQuery(setID string, queryID string) (string, error)
	MustGetQuery(setID string, queryID string) string
}

type SQLSetsProvider interface {
	GetAllMetas() []QuerySetMeta
}

type SQLSet struct {
	sets map[string]QuerySet
}

func (s *SQLSet) GetQuery(setID string, queryID string) (string, error) {
	return s.findQuery(setID, queryID)
}

func (s *SQLSet) MustGetQuery(setID string, queryID string) string {
	q, err := s.findQuery(setID, queryID)
	if err != nil {
		panic(err)
	}

	return q
}

func (s *SQLSet) GetAllMetas() []QuerySetMeta {
	metas := make([]QuerySetMeta, 0, len(s.sets))

	for _, qs := range s.sets {
		metas = append(metas, qs.GetMeta())
	}

	return metas
}

func (s *SQLSet) findQuery(setID string, queryID string) (string, error) {
	if s.sets == nil {
		return "", fmt.Errorf("%s: %w", setID, ErrQuerySetNotFound)
	}

	qs, ok := s.sets[setID]
	if !ok {
		return "", fmt.Errorf("%s: %w", setID, ErrQuerySetNotFound)
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

type QuerySet struct {
	meta    QuerySetMeta
	queries map[string]string
}

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
		return "", fmt.Errorf("%s: %w", id, ErrQueryNotFound)
	}

	q, ok := qs.queries[id]
	if !ok {
		return "", fmt.Errorf("%s: %w", id, ErrQueryNotFound)
	}

	return q, nil
}

type QuerySetMeta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
