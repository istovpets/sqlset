package sqlset

import (
	"fmt"
	"io/fs"
	"strings"
)

// New builds new SQLSet instance from the directory inside the given fs.FS.
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
