package sqlset

import (
	"fmt"
	"io/fs"
	"strings"
)

// New creates a new SQLSet by walking the directory tree of the provided fsys.
// It parses all .sql files it finds and adds them to the SQLSet.
// The walk starts from the root of the fsys. If you are using embed.FS
// and your queries are in a subdirectory, you should create a sub-filesystem
// using fs.Sub.
//
// Example with embed.FS:
//
//	//go:embed queries
//	var queriesFS embed.FS
//
//	sqlSet, err := sqlset.New(queriesFS)
func New(fsys fs.FS) (*SQLSet, error) {
	sqlSet := &SQLSet{}

	if err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		return handleDirEntry(fsys, sqlSet, path, entry)
	}); err != nil {
		return nil, fmt.Errorf("failed build SQL set: %w", err)
	}

	return sqlSet, nil
}

func handleDirEntry(fsys fs.FS, set *SQLSet, path string, entry fs.DirEntry) error {
	if entry.IsDir() {
		return nil
	}

	setID, ok := strings.CutSuffix(strings.ToLower(entry.Name()), filesExt)
	if !ok {
		return nil
	}

	f, err := fsys.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}

	defer func() {
		_ = f.Close()
	}()

	qs, err := parse(setID, f)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	set.registerQuerySet(qs.GetMeta().ID, qs)

	return nil
}
