package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
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

type FindFileType string

const (
	FileTypeFile FindFileType = "file"
	FileTypeDir               = "dir"
	FileTypeAll               = "file&dir"
)

type FileFindConfig struct {
	filesIncludeList         []*regexp.Regexp
	filesIgnoreList          []*regexp.Regexp
	foldersIncludeList       []*regexp.Regexp
	foldersIgnoreList        []*regexp.Regexp
	recurse                  bool
	fileType                 FindFileType
	excludeRootFolderInNames bool // Used currently in the find from embed.FS objects
	preloadContent           bool // Used currently in the find from embed.FS objects
}

type FileFindOption func(cfg *FileFindConfig)

func WithFindFileType(ft FindFileType) FileFindOption {
	return func(cfg *FileFindConfig) {
		cfg.fileType = ft
	}
}

func WithFindOptionNavigateSubDirs() FileFindOption {
	return func(cfg *FileFindConfig) {
		cfg.recurse = true
	}
}

// WithExcludeRootFolderInNames used only in the processing of embed.FS structures
func WithFindOptionExcludeRootFolderInNames() FileFindOption {
	return func(cfg *FileFindConfig) {
		cfg.excludeRootFolderInNames = true
	}
}

// WithPreloadContent used only in the processing of embed.FS structures
func WithFindOptionPreloadContent() FileFindOption {
	return func(cfg *FileFindConfig) {
		cfg.preloadContent = true
	}
}

func WithFindOptionFoldersIncludeList(p []string) FileFindOption {
	return func(cfg *FileFindConfig) {
		if len(p) == 0 {
			cfg.foldersIncludeList = nil
		} else {
			for _, s := range p {
				cfg.foldersIncludeList = append(cfg.foldersIncludeList, regexp.MustCompile(s))
			}
		}
	}
}

func WithFindOptionFoldersIgnoreList(p []string) FileFindOption {
	return func(cfg *FileFindConfig) {
		if len(p) == 0 {
			cfg.foldersIgnoreList = nil
		} else {
			for _, s := range p {
				cfg.foldersIgnoreList = append(cfg.foldersIgnoreList, regexp.MustCompile(s))
			}
		}
	}
}

func WithFindOptionFilesIncludeList(p []string) FileFindOption {
	return func(cfg *FileFindConfig) {
		if len(p) == 0 {
			cfg.filesIncludeList = nil
		} else {
			for _, s := range p {
				cfg.filesIncludeList = append(cfg.filesIncludeList, regexp.MustCompile(s))
			}
		}
	}
}

func WithFindOptionFilesIgnoreList(p []string) FileFindOption {
	return func(cfg *FileFindConfig) {
		if len(p) == 0 {
			cfg.filesIgnoreList = nil
		} else {
			for _, s := range p {
				cfg.filesIgnoreList = append(cfg.filesIgnoreList, regexp.MustCompile(s))
			}
		}
	}
}

// WithFindOptionIncludeList For backward compatibility it assigns same stuff to files and folders...
func WithFindOptionIncludeList(p []string) FileFindOption {
	return func(cfg *FileFindConfig) {
		if len(p) == 0 {
			cfg.filesIncludeList = nil
		} else {
			for _, s := range p {
				cfg.filesIncludeList = append(cfg.filesIncludeList, regexp.MustCompile(s))
			}
		}

		cfg.foldersIncludeList = cfg.filesIncludeList
	}
}

// WithFindOptionIgnoreList For backward compatibility it assigns same stuff to files and folders...
func WithFindOptionIgnoreList(p []string) FileFindOption {
	return func(cfg *FileFindConfig) {
		if len(p) == 0 {
			cfg.filesIgnoreList = nil
		} else {
			for _, s := range p {
				cfg.filesIgnoreList = append(cfg.filesIgnoreList, regexp.MustCompile(s))
			}
		}

		cfg.foldersIgnoreList = cfg.filesIgnoreList
	}
}

func (cfg *FileFindConfig) isIncluded(n string, includeList []*regexp.Regexp) bool {

	if len(includeList) == 0 {
		return true
	}

	for _, r := range includeList {
		if r.Match([]byte(n)) {
			return true
		}
	}

	return false
}

func (cfg *FileFindConfig) isExcluded(n string, ignoreList []*regexp.Regexp) bool {

	if len(ignoreList) == 0 {
		return false
	}

	for _, r := range ignoreList {
		if r.Match([]byte(n)) {
			return true
		}
	}

	return false
}

func (cfg *FileFindConfig) acceptFileName(n string, isDir bool, includeList, ignoreList []*regexp.Regexp) bool {
	if !cfg.isExcluded(n, ignoreList) {
		if cfg.isIncluded(n, includeList) {
			if (isDir && cfg.fileType != FileTypeFile) || (!isDir && cfg.fileType != FileTypeDir) {
				return true
			}
		}
	}

	return false
}

func FindFiles(folderPath string, opts ...FileFindOption) ([]string, error) {

	cfg := FileFindConfig{fileType: FileTypeAll}
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

				includeList := cfg.filesIncludeList
				ignoreList := cfg.filesIgnoreList
				if info.IsDir() {
					includeList = cfg.foldersIncludeList
					ignoreList = cfg.foldersIgnoreList
				}

				if cfg.acceptFileName(info.Name(), info.IsDir(), includeList, ignoreList) {
					files = append(files, path)
				}

				return nil
			})

		return files, err
	}

	fis, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		includeList := cfg.filesIncludeList
		ignoreList := cfg.filesIgnoreList
		if fi.IsDir() {
			includeList = cfg.foldersIncludeList
			ignoreList = cfg.foldersIgnoreList
		}

		p := filepath.Join(folderPath, fi.Name())
		if cfg.acceptFileName(fi.Name(), fi.IsDir(), includeList, ignoreList) {
			files = append(files, p)
		}
	}

	return files, nil
}
