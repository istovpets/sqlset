package sqlset

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	maxCapacity = 1024

	tokenPrefix  = "--"
	tokenKeySep  = ":"
	tokenComment = tokenPrefix
	tokenSQL     = "SQL"
	tokenMeta    = "META"
	tokenEnd     = "end"

	filesExt   = ".sql"
	lineEnding = "\r\n"
)

type parserToken struct {
	Type    string
	Key     string
	Content strings.Builder
}

//nolint:funlen
func parse(setID string, inp io.Reader) (QuerySet, error) {
	scanner := bufio.NewScanner(inp)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var (
		openedToken *parserToken
		lineN       int
		metaBuf     []byte
	)

	qs := QuerySet{}

	for scanner.Scan() {
		lineN++

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		token, key, err := detectToken(line)
		if err != nil {
			return QuerySet{}, fmt.Errorf("line %d: %w", lineN, err)
		}

		if openedToken != nil && (token == tokenSQL || token == tokenMeta) {
			return QuerySet{}, fmt.Errorf(
				"line %d: %w: unexpected %s inside %s",
				lineN, ErrInvalidSyntax, token, openedToken.Type,
			)
		}

		switch token {
		case tokenComment:
			continue
		case tokenSQL:
			openedToken = &parserToken{
				Type: tokenSQL,
				Key:  key,
			}

			continue
		case tokenMeta:
			if metaBuf != nil {
				return QuerySet{}, fmt.Errorf("line %d: %w: unexpected multiple metadata", lineN, ErrInvalidSyntax)
			}
			openedToken = &parserToken{Type: tokenMeta}

			continue
		}

		if token == tokenEnd {
			if openedToken == nil {
				return QuerySet{}, fmt.Errorf(
					"line %d: %w: unexpected '%s' token",
					lineN, ErrInvalidSyntax, tokenEnd,
				)
			}

			switch {
			case openedToken.Type == tokenSQL:
				qs.registerQuery(
					openedToken.Key,
					strings.TrimSuffix(openedToken.Content.String(), lineEnding),
				)
			case openedToken.Type == tokenMeta:
				metaBuf = []byte(openedToken.Content.String())
			}

			openedToken.Content.Reset()
			openedToken = nil

			continue
		}

		if openedToken == nil {
			continue
		}

		openedToken.Content.WriteString(line + lineEnding)
	}

	if err := scanner.Err(); err != nil {
		if errors.Is(err, bufio.ErrTooLong) {
			return QuerySet{}, fmt.Errorf("line %d: %w", lineN+1, ErrMaxLineLenExceeded)
		}

		return QuerySet{}, fmt.Errorf("scanning error: %w", err)
	}

	if openedToken != nil {
		return QuerySet{}, fmt.Errorf(
			"%w: no closing tag found for '%s:%s'",
			ErrInvalidSyntax, openedToken.Type, openedToken.Key,
		)
	}

	meta, err := parseMeta(setID, metaBuf)
	if err != nil {
		return qs, fmt.Errorf("parse meta: %w", err)
	}

	qs.meta = meta

	return qs, nil
}

func detectToken(line string) (token string, key string, err error) {
	var ok bool

	line, ok = strings.CutPrefix(line, tokenPrefix)
	if !ok {
		// Not a token nor comment, skipping.
		return "", "", nil
	}

	// SQL:key
	key, ok = strings.CutPrefix(line, tokenSQL+tokenKeySep)
	if ok {
		key = strings.TrimSpace(key)
		if key == "" {
			return "", "", fmt.Errorf("%w: no SQL set query key given", ErrInvalidSyntax)
		}

		return tokenSQL, key, nil
	}

	// META
	if strings.HasPrefix(line, tokenMeta) {
		return tokenMeta, "", nil
	}

	// --end
	if strings.HasPrefix(line, tokenEnd) {
		return tokenEnd, "", nil
	}

	// Just a comment
	return tokenComment, "", nil
}

func parseMeta(setID string, jsonData []byte) (QuerySetMeta, error) {
	meta := QuerySetMeta{
		ID:   setID,
		Name: setID,
	}

	if jsonData == nil {
		return meta, nil
	}

	var parsed QuerySetMeta

	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		return QuerySetMeta{}, fmt.Errorf("%w: %s", ErrInvalidSyntax, err.Error())
	}

	if parsed.ID != "" {
		meta.ID = parsed.ID
	}

	if parsed.Name != "" {
		meta.Name = parsed.Name
	}

	meta.Description = parsed.Description

	return meta, nil
}
