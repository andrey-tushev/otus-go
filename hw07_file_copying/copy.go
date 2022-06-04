package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	stat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if stat.Size() == 0 {
		return ErrUnsupportedFile
	}
	if offset > 0 && offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}

	inFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	const bufSize = 512
	buf := make([]byte, bufSize)

	if offset > 0 {
		inFile.Seek(offset, io.SeekStart)
	}

	totalRead := 0
	for done := false; !done; {
		readSize, inErr := inFile.Read(buf)
		if inErr == io.EOF {
			return nil
		} else if inErr != nil {
			return inErr
		}
		buf = buf[:readSize]
		totalRead += readSize

		if limit > 0 && totalRead > int(limit) {
			tail := bufSize - (totalRead - int(limit))
			buf = buf[:tail]
			done = true
		}

		_, outErr := outFile.Write(buf)
		if outErr != nil {
			return outErr
		}
	}

	return nil
}
