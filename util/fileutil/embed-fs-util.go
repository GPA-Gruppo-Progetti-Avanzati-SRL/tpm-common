package fileutil

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FoundFile struct {
	Path    string
	Info    fs.FileInfo
	Content []byte
}

func FindEmbeddedFiles(embeddedFS embed.FS, folderPath string, opts ...FileFindOption) ([]FoundFile, error) {

	cfg := FileFindConfig{fileType: FileTypeAll}
	for _, o := range opts {
		o(&cfg)
	}

	var blockedPath string
	var files []FoundFile
	if cfg.recurse {
		err := walkEmbeddedFS(embeddedFS, folderPath,
			func(fPath string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if blockedPath != "" && strings.HasPrefix(fPath, blockedPath) {
					return nil
				} else {
					blockedPath = ""
				}

				includeList := cfg.filesIncludeList
				ignoreList := cfg.filesIgnoreList
				if info.IsDir() {
					includeList = cfg.foldersIncludeList
					ignoreList = cfg.foldersIgnoreList
				}

				if cfg.acceptFileName(info.Name(), info.IsDir(), includeList, ignoreList) {

					var data []byte
					if cfg.preloadContent && !info.IsDir() {
						data, err = embeddedFS.ReadFile( /* file */ path.Join(fPath, info.Name()))
					}

					// this option in principle should replace the excludeRootFolderInNames.
					// the other is still there for backward compatibility
					if cfg.trimRootFolderFromNames {
						fPath = strings.TrimPrefix(fPath, folderPath)
					}

					if cfg.excludeRootFolderInNames {
						ndx := strings.Index(fPath, "/")
						if ndx >= 0 {
							fPath = fPath[ndx+1:]
						} else {
							fPath = ""
						}
					}

					files = append(files, FoundFile{Path: fPath, Info: info, Content: data})
				} else {
					if info.IsDir() {
						blockedPath = filepath.Join(fPath, info.Name())
					}
				}

				return nil
			})

		return files, err
	}

	fis, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	p := ""
	if !cfg.excludeRootFolderInNames {
		p = folderPath
	}

	for _, fi := range fis {
		info, err := fi.Info()
		if err != nil {
			return files, err
		}

		includeList := cfg.filesIncludeList
		ignoreList := cfg.filesIgnoreList
		if info.IsDir() {
			includeList = cfg.foldersIncludeList
			ignoreList = cfg.foldersIgnoreList
		}

		if cfg.acceptFileName(fi.Name(), fi.IsDir(), includeList, ignoreList) {
			var data []byte
			if cfg.preloadContent && !info.IsDir() {
				data, err = embeddedFS.ReadFile( /* file */ path.Join(folderPath, info.Name()))
			}
			files = append(files, FoundFile{Path: p, Info: info, Content: data})
		}
	}

	return files, nil
}

type EmbeddedFSVisitor func(fPath string, info fs.FileInfo, err error) error

func walkEmbeddedFS(embedFS embed.FS, fPath string, visitor EmbeddedFSVisitor) error {
	entries, err := embedFS.ReadDir(fPath)
	if err != nil {
		return err
	}

	var info fs.FileInfo
	for _, e := range entries {

		info, err = e.Info()
		err = visitor(fPath, info, err)

		if e.IsDir() {
			fn := /* file */ path.Join(fPath, e.Name())
			err = walkEmbeddedFS(embedFS, fn, visitor)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
