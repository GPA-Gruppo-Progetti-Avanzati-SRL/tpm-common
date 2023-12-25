package util

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func FindEmbeddedFiles(embeddedFS embed.FS, folderPath string, opts ...FileFindOption) ([]string, error) {

	cfg := fileFindConfig{fileType: FileTypeAll}
	for _, o := range opts {
		o(&cfg)
	}

	var files []string
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
						}
					}
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
		p := fi.Name()
		if !cfg.excludeRootFolderInNames {
			p = filepath.Join(folderPath, fi.Name())
		}

		if cfg.acceptFileName(fi.Name(), fi.IsDir()) {
			files = append(files, p)
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
		fn := filepath.Join(path, e.Name())
		err = visitor(fn, info, err)

		if e.IsDir() {
			err = walkEmbeddedFS(embedFS, fn, visitor)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
