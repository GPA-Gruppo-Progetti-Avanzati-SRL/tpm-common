package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/rs/zerolog/log"
)

func BufoReaderReadLineAsString(r *bufio.Reader, lineNo int, maxLength int) (string, error) {
	b, err := BufoReaderReadLine(r, lineNo, maxLength)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func BufoReaderReadLine(r *bufio.Reader, lineNo int, maxLength int) ([]byte, error) {
	const semLogContext = "bufio-util::read-line"

	rawLine, isPrefix, err := r.ReadLine()
	if err != nil {
		return rawLine, err
	}

	if !isPrefix {
		if maxLength > 0 && len(rawLine) > maxLength {
			err = errors.New("max length exceeded")
		}
		return rawLine, err
	}

	var line bytes.Buffer
	var n int
	n, err = line.Write(rawLine)
	if err != nil {
		return nil, err
	}
	longLineLength := n
	for err == nil && isPrefix {
		rawLine, isPrefix, err = r.ReadLine()
		if err == nil {
			n, err = line.Write(rawLine)
			longLineLength += n
		}

		if maxLength > 0 && longLineLength > maxLength {
			err = errors.New("max length exceeded")
			return nil, err
		}
	}

	if err == nil || err == io.EOF {
		log.Warn().Err(err).Int("line-number", lineNo).Int("long-line-length", longLineLength).Msg(semLogContext)
	}

	return line.Bytes(), err
}
