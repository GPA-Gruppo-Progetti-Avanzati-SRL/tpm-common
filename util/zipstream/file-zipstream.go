package zipstream

import (
	"archive/zip"
	"github.com/rs/zerolog/log"
	"io/fs"
	"io/ioutil"
	"os"
)

type fileZipStreamImpl struct {
	zipStreamImpl
	rootFolder string
	filename   string
	file       *os.File
}

//func (zs *fileZipStreamImpl) Close() error {
//	if zs.writer != nil {
//		if err := zs.writer.Close(); err != nil {
//			return err
//		}
//	}
//
//	zs.writer = nil
//	return nil
//}

func (zs *fileZipStreamImpl) Bytes() []byte {

	if zs.filename != "" {
		b, err := ioutil.ReadFile(zs.filename)
		if err != nil {
			log.Error().Err(err).Str("fileName", zs.filename).Msg("fileZipStreamImpl error reading from file")
		}
		return b
	}
	return nil
}

func (zs *fileZipStreamImpl) Add(fn string, fbody []byte) error {

	if zs.writer == nil {
		f, err := ioutil.TempFile(zs.rootFolder, "archive-*.zip")
		if err != nil {
			return err
		}

		zs.file = f
		zs.filename = f.Name()
		zs.writer = zip.NewWriter(f)

		log.Trace().Str("file-name", zs.filename).Msg("zip stream file created")
	}

	return zs.addToStream(fn, fbody)
}

func (zs *fileZipStreamImpl) Close() error {

	var errZip, errFile error
	errZip = zs.zipStreamImpl.Close()
	if zs.file != nil {
		log.Trace().Str("file-name", zs.filename).Msg("closing zip stream file")
		errFile = zs.file.Close()
		zs.file = nil
	}

	if errZip != nil || errFile != nil {
		zs.Dispose()
	}

	if errZip != nil {
		return errZip
	}

	return errFile
}

func (zs *fileZipStreamImpl) Dispose() error {

	closeErr := zs.Close()
	if closeErr != nil {
		log.Error().Err(closeErr).Msg("fileZipStreamImpl release resources")
	}

	if zs.filename != "" {
		log.Trace().Str("file-name", zs.filename).Msg("disposing zip stream file")

		err := os.Remove(zs.filename)
		if err != nil {
			log.Error().Err(err).Msg("fileZipStreamImpl release resources")
		}

		zs.filename = ""

		if closeErr != nil {
			return closeErr
		}

		return err
	}

	return closeErr
}

func (zs *fileZipStreamImpl) CloseAndSave(fn string) error {

	err := zs.Close()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, zs.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
