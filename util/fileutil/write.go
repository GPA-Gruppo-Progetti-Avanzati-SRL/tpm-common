package fileutil

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type WriteOptions struct {
	createFolderIfMissing bool
}

type WriteOption func(*WriteOptions)

func WithWriteOptionCreateFolderIfMissing() WriteOption {
	return func(o *WriteOptions) {
		o.createFolderIfMissing = true
	}
}

func WriteFile(fn string, b []byte, fileMode os.FileMode, writeOpts ...WriteOption) error {
	const semLogContext = "file-util::write-2-file"
	var err error

	options := WriteOptions{}
	for _, funOpt := range writeOpts {
		funOpt(&options)
	}

	outFn, _ := ResolvePath(fn)
	outFolder := filepath.Dir(outFn)
	if !FileExists(outFolder) {
		if options.createFolderIfMissing {
			err = os.MkdirAll(outFolder, fileMode)
			if err != nil {
				log.Error().Err(err).Str("folder-name", outFolder).Msg(semLogContext)
				return err
			}
		} else {
			err = errors.New("out folder path not found")
			log.Error().Err(err).Str("folder-name", outFolder).Msg(semLogContext)
			return err
		}
	}

	err = os.WriteFile(outFn, b, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Str("file-name", fn).Msg(semLogContext)
	}

	return err
}
