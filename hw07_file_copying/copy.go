package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type ProgressReader struct {
	io.Reader
	bar *pb.ProgressBar
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n > 0 {
		for i := 0; i < n; i++ {
			pr.bar.Increment()
			// time.Sleep(time.Millisecond)
		}
	}
	return n, err
}

func Copy(from, to string, offset, limit int64) error {
	srcFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if offset > srcFileInfo.Size() {
		return fmt.Errorf("offset exceeds file size")
	}

	if limit == 0 || limit > srcFileInfo.Size()-offset {
		limit = srcFileInfo.Size() - offset
	}

	dstFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	var totalCopied int64
	for totalCopied < limit {
		bytesToRead := int64(len(buf))
		if limit-totalCopied < bytesToRead {
			bytesToRead = limit - totalCopied
		}

		n, err := srcFile.Read(buf[:bytesToRead])
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n == 0 {
			break
		}

		_, err = dstFile.Write(buf[:n])
		if err != nil {
			return err
		}

		totalCopied += int64(n)
	}

	return nil
}
