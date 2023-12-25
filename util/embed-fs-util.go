package util

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FoundFile struct {
	Path string
	Info fs.FileInfo
}

func FindEmbeddedFiles(embeddedFS embed.FS, folderPath string, opts ...FileFindOption) ([]FoundFile, error) {

	cfg := fileFindConfig{fileType: FileTypeAll}
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

				if cfg.acceptFileName(info.Name(), info.IsDir()) {
					if cfg.excludeRootFolderInNames {
						ndx := strings.Index(path, "/")
						if ndx >= 0 {
							path = path[ndx+1:]
						} else {
							path = ""
						}
					}
					files = append(files, FoundFile{Path: path, Info: info})
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

		if cfg.acceptFileName(fi.Name(), fi.IsDir()) {
			files = append(files, FoundFile{Path: p, Info: info})
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
