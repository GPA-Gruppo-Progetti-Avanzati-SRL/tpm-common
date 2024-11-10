package fileutil

import (
	"errors"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyFolder(dst, src string) (int, error) {
	const semLogContext = "file-util::copy-folder"
	var err error

	src, ok := util.ResolvePath(src)
	if !ok {
		err = errors.New("src folder path not found")
		log.Error().Err(err).Str("file-name", src).Msg(semLogContext)
		return 0, err
	}

	dst, ok = util.ResolvePath(dst)
	if !ok {
		err = errors.New("src folder path not found")
		log.Error().Err(err).Str("file-name", dst).Msg(semLogContext)
		return 0, err
	}

	opts := []FileFindOption{WithFindFileType(FileTypeFile)}
	assetFiles, err := FindFiles(src, opts...)
	if err != nil {
		log.Error().Err(err).Str("file-name", dst).Msg(semLogContext)
		return 0, err
	}

	for _, f := range assetFiles {
		destFile := filepath.Join(dst, filepath.Base(f))
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
