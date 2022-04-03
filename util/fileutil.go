package util

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

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

type FindFileType string

const (
	FileTypeFile FindFileType = "file"
	FileTypeDir               = "dir"
	FileTypeAll               = "file&dir"
)

type fileFindConfig struct {
	includeList []*regexp.Regexp
	ignoreList  []*regexp.Regexp
	recurse     bool
	fileType    FindFileType
}

type FileFindOption func(cfg *fileFindConfig)

func WithFindFileType(ft FindFileType) FileFindOption {
	return func(cfg *fileFindConfig) {
		cfg.fileType = ft
	}
}

func WithFindOptionNavigateSubDirs() FileFindOption {
	return func(cfg *fileFindConfig) {
		cfg.recurse = true
	}
}

func WithFindOptionIncludeList(p []string) FileFindOption {
	return func(cfg *fileFindConfig) {
		if len(p) == 0 {
			cfg.includeList = nil
		} else {
			for _, s := range p {
				cfg.includeList = append(cfg.includeList, regexp.MustCompile(s))
			}
		}
	}
}

func WithFindOptionIgnoreList(p []string) FileFindOption {
	return func(cfg *fileFindConfig) {
		if len(p) == 0 {
			cfg.ignoreList = nil
		} else {
			for _, s := range p {
				cfg.ignoreList = append(cfg.ignoreList, regexp.MustCompile(s))
			}
		}
	}
}

func (cfg *fileFindConfig) isIncluded(n string) bool {

	if len(cfg.includeList) == 0 {
		return true
	}

	for _, r := range cfg.includeList {
		if r.Match([]byte(n)) {
			return true
		}
	}

	return false
}

func (cfg *fileFindConfig) isExcluded(n string) bool {

	if len(cfg.ignoreList) == 0 {
		return false
	}

	for _, r := range cfg.ignoreList {
		if r.Match([]byte(n)) {
			return true
		}
	}

	return false
}

func (cfg *fileFindConfig) acceptFileName(n string, isDir bool) bool {
	if !cfg.isExcluded(n) {
		if cfg.isIncluded(n) {
			if (isDir && cfg.fileType != FileTypeFile) || (!isDir && cfg.fileType != FileTypeDir) {
				return true
			}
		}
	}

	return false
}

func FindFiles(folderPath string, opts ...FileFindOption) ([]string, error) {

	cfg := fileFindConfig{fileType: FileTypeAll}
	for _, o := range opts {
		o(&cfg)
	}

	var files []string
	if cfg.recurse {
		err := filepath.Walk(folderPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if cfg.acceptFileName(info.Name(), info.IsDir()) {
					files = append(files, path)
				}

				return nil
			})

		return files, err
	}

	fis, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		p := filepath.Join(folderPath, fi.Name())
		if cfg.acceptFileName(fi.Name(), fi.IsDir()) {
			files = append(files, p)
		}
	}

	return files, nil
}
