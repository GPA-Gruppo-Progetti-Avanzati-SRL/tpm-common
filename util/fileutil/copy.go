package fileutil

import (
	"errors"

	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
	"path/filepath"
)

type CopyOptions struct {
	createIfMissing   bool
	includeSubFolders bool
}

type CopyOption func(*CopyOptions)

func WithCopyOptionCreateIfMissing() CopyOption {
	return func(o *CopyOptions) {
		o.createIfMissing = true
	}
}

func WithCopyOptionIncludeSubFolder() CopyOption {
	return func(o *CopyOptions) {
		o.includeSubFolders = true
	}
}

func CopyFolder(dst, src string, copyOpts ...CopyOption) (int, error) {
	const semLogContext = "file-util::copy-folder"
	var err error

	options := CopyOptions{}
	for _, funOpt := range copyOpts {
		funOpt(&options)
	}

	src, ok := ResolvePath(src)
	if !ok {
		err = errors.New("src folder path not found")
		log.Error().Err(err).Str("file-name", src).Msg(semLogContext)
		return 0, err
	}

	dst, ok = ResolvePath(dst)
	if !ok {
		if options.createIfMissing {
			err = os.MkdirAll(dst, os.ModePerm)
			if err != nil {
				log.Error().Err(err).Str("file-name", dst).Msg(semLogContext)
				return 0, err
			}
		} else {
			err = errors.New("src folder path not found")
			log.Error().Err(err).Str("file-name", dst).Msg(semLogContext)
			return 0, err
		}
	}

	opts := []FileFindOption{WithFindFileType(FileTypeFile)}
	if options.includeSubFolders {
		opts = append(opts, WithFindOptionNavigateSubDirs())
	}

	assetFiles, err := FindFiles(src, opts...)
	if err != nil {
		log.Error().Err(err).Str("file-name", dst).Msg(semLogContext)
		return 0, err
	}

	for _, f := range assetFiles {
		fnDir, fn, err := relPath(src, f)
		if err != nil {
			log.Error().Err(err).Str("file-name", f).Msg(semLogContext)
			return 0, err
		}

		destFolder := dst
		if fnDir != "" {
			destFolder = filepath.Join(dst, fnDir)
			if !FileExists(destFolder) {
				if options.createIfMissing {
					err = os.MkdirAll(destFolder, os.ModePerm)
					if err != nil {
						log.Error().Err(err).Str("folder-name", destFolder).Msg(semLogContext)
						return 0, err
					}
				}
			}
		}
		destFile := filepath.Join(destFolder, fn)
		log.Info().Str("src-file", f).Str("dst-file", destFile).Msg(semLogContext)
		b, err := os.ReadFile(f)
		if err != nil {
			log.Error().Err(err).Str("file-name", dst).Msg(semLogContext)
			return 0, err
		}
		err = os.WriteFile(destFile, b, fs.ModePerm)
		if err != nil {
			log.Error().Err(err).Str("file-name", destFile).Msg(semLogContext)
			return 0, err
		}
	}

	return 0, nil
}

func relPath(src, f string) (string, string, error) {
	const semLogContext = "file-util::copy-folder-rel-path"
	p, err := filepath.Rel(src, f)
	if err != nil {
		log.Error().Err(err).Str("src", src).Str("f", f).Msg(semLogContext)
		return "", "", err
	}

	d := filepath.Dir(p)
	if d == "." {
		d = ""
	}

	return d, filepath.Base(p), nil
}
