package fileutil

import (
	"errors"
	"os"
)

func FileInfo(fn string) (os.FileInfo, bool) {
	if fi, err := os.Stat(fn); err == nil {
		return fi, true

	} else if errors.Is(err, os.ErrNotExist) {
		return nil, false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return nil, false
	}
}

func FileExists(fn string) bool {
	if _, err := os.Stat(fn); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

func FileSize(fn string) int64 {
	if fi, err := os.Stat(fn); err == nil {
		return fi.Size()

	} else if errors.Is(err, os.ErrNotExist) {
		return -1

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return -1
	}
}

// emptyDirAllowedEntries are the entries a pre-existing project folder may hold and still count as
// "effectively empty" for the new-project guard: a README.md placeholder and/or a .git folder (both
// optional). A folder with only these — or truly empty — is scaffolded into without --force.
var emptyDirAllowedEntries = map[string]bool{"README.md": true, ".git": true}

// RepoIsEmpty reports whether dir contains nothing but the allowed placeholder entries
// (README.md, .git). A truly empty directory is effectively empty too.
func RepoIsEmpty(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if !emptyDirAllowedEntries[e.Name()] {
			return false, nil
		}
	}
	return true, nil
}
