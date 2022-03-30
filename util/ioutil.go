package util

import "io"

// source: https://github.com/Azure/azure-sdk-for-go/blob/fa22cece712452b2ae6ec41751707e1d949fcf3d/sdk/azcore/internal/shared/shared.go#L59

type nopReaderSeekCLoser struct {
	io.ReadSeeker
}

func (n nopReaderSeekCLoser) Close() error {
	return nil
}

func NopReaderSeekCLoser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopReaderSeekCLoser{rs}
}
