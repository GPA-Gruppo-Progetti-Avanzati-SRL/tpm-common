package fileutil

import (
	"embed"
	"io/fs"
	"os"
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

	var files []FoundFile
	if cfg.recurse {
		err := walkEmbeddedFS(embeddedFS, folderPath,
			func(path string, info fs.FileInfo, err error) error {
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

					var data []byte
					if cfg.preloadContent && !info.IsDir() {
						data, err = embeddedFS.ReadFile(filepath.Join(path, info.Name()))
					}

					if cfg.excludeRootFolderInNames {
						ndx := strings.Index(path, "/")
						if ndx >= 0 {
							path = path[ndx+1:]
						} else {
							path = ""
						}
					}

					files = append(files, FoundFile{Path: path, Info: info, Content: data})
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
				data, err = embeddedFS.ReadFile(filepath.Join(folderPath, info.Name()))
			}
			files = append(files, FoundFile{Path: p, Info: info, Content: data})
		}
	}

	return files, nil
}

type EmbeddedFSVisitor func(path string, info fs.FileInfo, err error) error

func walkEmbeddedFS(embedFS embed.FS, path string, visitor EmbeddedFSVisitor) error {
	entries, err := embedFS.ReadDir(path)
	if err != nil {
		return err
	}

	var info fs.FileInfo
	for _, e := range entries {

		info, err = e.Info()
		err = visitor(path, info, err)

		if e.IsDir() {
			fn := filepath.Join(path, e.Name())
			err = walkEmbeddedFS(embedFS, fn, visitor)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
