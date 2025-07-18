package zipper

import (
	"compress/gzip"
	"io"
)

func Zip(file io.Reader) (io.Reader, func() error, chan error) {
	pipeR, pipeW := io.Pipe()
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		compressed := gzip.NewWriter(pipeW)
		if _, err := io.Copy(compressed, file); err != nil {
			defer pipeR.Close()
			errCh <- err
			return
		}
		if err := compressed.Close(); err != nil {
			defer pipeR.Close()
			errCh <- err
			return
		}
		if err := pipeW.Close(); err != nil {
			defer pipeR.Close()
			errCh <- err
			return
		}
	}()
	return pipeR, pipeR.Close, errCh
}
