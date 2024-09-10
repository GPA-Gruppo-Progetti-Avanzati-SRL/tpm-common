package splitutil

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
)

type Chunk struct {
	Offset        int64 `json:"offset,omitempty" yaml:"offset,omitempty" mapstructure:"offset,omitempty"`
	Size          int64 `json:"size,omitempty" yaml:"size,omitempty" mapstructure:"size,omitempty"`
	IsLast        bool  `json:"is-last,omitempty" yaml:"is-last,omitempty" mapstructure:"is-last,omitempty"`
	ChunkEdgeSize int64 `json:"edge-size,omitempty" yaml:"edge-size,omitempty" mapstructure:"edge-size,omitempty"`
	scanner       *bufio.Scanner
	consumedBytes int64
}

func (chunk *Chunk) Range() (int64, int64) {
	return chunk.Offset, chunk.Size + chunk.ChunkEdgeSize
}

func (chunk *Chunk) NewReader(b []byte, size int64) (*bufio.Scanner, error) {
	const semLogContext = "chunk::new-reader"
	/*
		if size < (chunk.Size + chunk.ChunkEdgeSize) {
			err := fmt.Errorf("invalid chunk: %v. size is %d instead of %d", chunk, size, chunk.Size+chunk.ChunkEdgeSize)
			return nil, err
		}
	*/

	if size != int64(len(b)) {
		b = b[:size]
	}

	chunk.scanner = bufio.NewScanner(bytes.NewReader(b))
	chunk.scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		const semLogContext = "chunk::split-func"
		advance, token, err = bufio.ScanLines(data, atEOF)
		chunk.consumedBytes += int64(advance)
		var evt *zerolog.Event
		if err != nil {
			evt = log.Info().Err(err)
		} else {
			evt = log.Trace()
		}
		evt.Int("advance", advance).Str("token", string(token)).Msg(semLogContext)
		return
	})

	if chunk.Offset != 0 {
		ok := chunk.scanner.Scan()
		log.Trace().Bool("unused-first-scan", ok).Msg(semLogContext)
	}

	return chunk.scanner, nil
}

func (chunk *Chunk) Read() (string, error) {
	const semLogContext = "chunk::read"
	var err error
	if chunk.scanner == nil {
		err = errors.New("chunk scanner is nil")
		return "", err
	}

	if chunk.consumedBytes <= chunk.Size {
		ok := chunk.scanner.Scan()
		if ok {
			text := chunk.scanner.Text()
			log.Trace().Int("scanned-bytes", len(chunk.scanner.Bytes())).Str("scanned-text", text).Msg(semLogContext)
			return text, nil
		} else {
			if err := chunk.scanner.Err(); err != nil {
				return "", err
			}
		}
	}
	return "", io.EOF

	/*
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // Println will add back the final '\n'
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	*/
}

func (chunk *Chunk) Close() error {
	if chunk.scanner != nil {
		chunk.scanner = nil
	}

	return nil
}

func GetChunksFromSize(sz int64, chunkSize int64, chunkEdgeSize int64) []Chunk {
	numOfChunks := int((sz-1)/chunkSize + 1)
	if numOfChunks == 1 {
		chunkSize = sz
	}

	if chunkSize+chunkEdgeSize > sz {
		chunkEdgeSize = sz - chunkSize
	}

	var chunks []Chunk
	var offset int64
	for i := 0; i < numOfChunks; i++ {

		isLast := false
		edgeSize := chunkEdgeSize
		if i == (numOfChunks - 1) {
			edgeSize = 0
			isLast = true
			chunkSize = sz - offset
		}

		chunks = append(chunks, Chunk{
			Offset:        offset,
			Size:          chunkSize,
			ChunkEdgeSize: edgeSize,
			IsLast:        isLast,
		})

		offset += chunkSize
	}

	return chunks
}
