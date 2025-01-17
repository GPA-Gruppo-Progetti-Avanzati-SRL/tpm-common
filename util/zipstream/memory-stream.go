package zipstream

import (
	"archive/zip"
	"bytes"
	"io/fs"
	"os"
)

type memoryZipStreamImpl struct {
	zipStreamImpl
	buffer *bytes.Buffer
}

func (zs *memoryZipStreamImpl) Bytes() []byte {
	return zs.buffer.Bytes()
}

func (zs *memoryZipStreamImpl) Add(fn string, fbody []byte) error {
	if zs.writer == nil {
		zs.writer = zip.NewWriter(zs.buffer)
	}

	return zs.addToStream(fn, fbody)
}

// CloseAndSave
// Convenience for dumping the file somewhere for inspection....
func (zs *memoryZipStreamImpl) CloseAndSave(fn string) error {

	err := zs.Close()
	if err != nil {
		return err
	}

	err = os.WriteFile(fn, zs.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
