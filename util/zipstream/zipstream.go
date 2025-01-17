package zipstream

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/fs"
	"time"
)

type ZipStream interface {
	Close() error
	Dispose() error
	Bytes() []byte
	Add(fn string, fbody []byte) error
	CloseAndSave(fn string) error
	IsEmpty() bool
	Size() int
}

type zipStreamImpl struct {
	numberOfItems int
	siz           int
	writer        *zip.Writer
}

func (zs *zipStreamImpl) IsEmpty() bool {
	return zs.numberOfItems == 0

}

func (zs *zipStreamImpl) Size() int {
	return zs.siz
}

func (zs *zipStreamImpl) Dispose() error {
	err := zs.Close()
	return err
}

func (zs *zipStreamImpl) Close() error {

	var err error
	if zs.writer != nil {
		log.Trace().Msg("closing zip stream writer")
		err = zs.writer.Close()
		zs.writer = nil
	}

	return err
}

func (zs *zipStreamImpl) addToStream(fn string, fbody []byte) error {

	size := len(fbody)
	fh := &zip.FileHeader{
		Name:               fn,
		UncompressedSize64: uint64(size),
	}
	fh.Modified = time.Now()
	fh.SetMode(fs.ModePerm)
	fh.Method = zip.Deflate

	f, err := zs.writer.CreateHeader(fh)
	if err != nil {
		return err
	}

	n, err := f.Write(fbody)
	if err != nil {
		return err
	}

	zs.siz += n
	zs.numberOfItems++
	return nil
}

const (
	MemZipStream  = "mem"
	DiskZipStream = "disk"
)

func NewZipStream(streamType string, diskFolderName string) (ZipStream, error) {

	log.Trace().Str("stream-type", streamType).Str("folder", diskFolderName).Msg("new zip stream instance")

	var zs ZipStream
	switch streamType {
	case MemZipStream:
		zs = &memoryZipStreamImpl{buffer: new(bytes.Buffer)}
	case DiskZipStream:
		zs = &fileZipStreamImpl{rootFolder: diskFolderName}
	default:
		return nil, fmt.Errorf("NewZipStream: unrecognized streamType param %s", streamType)
	}

	return zs, nil
}
