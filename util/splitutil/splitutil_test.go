package splitutil_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/splitutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

const (
	mb = 1024 * 1024
	gb = 1024 * mb

	//fileToSplit  = "../../local-files/gameofthrones.txt"
	//chunkSize    = 1 * mb
	//chunkEdgSize = 200

	fileToSplit  = "../../local-files/short-file.txt"
	chunkSize    = 5231
	chunkEdgSize = 200
)

func TestSplit(t *testing.T) {
	const semLogContext = "chunk::test-split"

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	chunks, err := splitutil.GetChunksOfFile(fileToSplit, chunkSize, chunkEdgSize)
	require.NoError(t, err)

	file, err := os.Open(fileToSplit)
	require.NoError(t, err)
	defer file.Close()

	totalNumberOfLines := 0
	for chunkNumber, chunk := range chunks {

		log.Info().Int("chunk-number", chunkNumber).Interface("chunk", chunk).Msg(semLogContext + "======================")

		offset, totalSize := chunk.Range()
		actualOffset, err := file.Seek(offset, 0)
		require.NoError(t, err)
		log.Trace().Int64("requested-offset", offset).Int64("actual", actualOffset).Interface("chunk", chunk).Msg(semLogContext)

		b := make([]byte, totalSize)
		actualBytes, err := file.Read(b)
		require.NoError(t, err)
		log.Trace().Int("requested-bytes", len(b)).Int("actual", actualBytes).Interface("chunk", chunk).Msg(semLogContext)

		_, err = chunk.NewReader(b, int64(actualBytes))
		numLines := 0
		firstLine := ""
		lastLine := ""
		for {
			var line string
			line, err = chunk.Read()
			if err == io.EOF {
				log.Info().Int("num-lines", numLines).Int("tot-num-lines", totalNumberOfLines).Msg(semLogContext + " End Of Chunk")
				break
			} else {
				if numLines == 0 {
					firstLine = line
				}
				lastLine = line
			}
			require.NoError(t, err)
			numLines++
			totalNumberOfLines++
			log.Trace().Str("line", line).Msg(semLogContext)
		}

		log.Info().Str("last-line", lastLine).Str("first-line", firstLine).Msg(semLogContext)

	}

	log.Info().Int("tot-num-lines", totalNumberOfLines).Msg(semLogContext + " EOF")
}
