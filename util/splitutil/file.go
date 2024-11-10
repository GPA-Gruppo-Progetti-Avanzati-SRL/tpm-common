package splitutil

import (
	"errors"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
)

func GetChunksOfFile(fileName string, chunkSize int64, chunkEdgeSize int64) ([]Chunk, error) {
	fi, ok := fileutil.FileInfo(fileName)
	if !ok {
		err := errors.New("file not found")
		return nil, err
	}

	if fi.IsDir() {
		err := errors.New("file is not regular file")
		return nil, err
	}

	sz := fi.Size()
	chunks := GetChunksFromSize(sz, chunkSize, chunkEdgeSize)

	return chunks, nil
}
