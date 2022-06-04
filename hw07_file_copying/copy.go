package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Проверим размер входного файла
	stat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	size := stat.Size()
	if size == 0 {
		return ErrUnsupportedFile
	}

	if limit == 0 {
		limit = size - offset
	}

	if offset > 0 && offset > size {
		return ErrOffsetExceedsFileSize
	}

	// Откроем входной файл
	inFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// Откроем выходной файл
	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Если задано смещение, то смести позицию во входном файле
	if offset > 0 {
		inFile.Seek(offset, io.SeekStart)
	}

	// Буфер для чтения
	const bufSize = 1000
	buf := make([]byte, bufSize)

	// Настроим прогресс бар
	bytesToRead := int(limit)
	if offset+limit > size {
		bytesToRead = int(size - offset)
	}
	steps := int(math.Ceil(float64(bytesToRead) / float64(bufSize)))
	bar := pb.StartNew(steps)
	fmt.Println(bytesToRead, bufSize, steps)
	defer bar.Finish()

	// Цикл чтения по кусочкам
	totalRead := 0
	for done := false; !done; {
		readSize, inErr := inFile.Read(buf)

		// Входной файл закончился
		if inErr == io.EOF {
			return nil
		}

		// Какая-то ошибка чтения
		if inErr != nil {
			return inErr
		}

		// Если прочитали меньше чем размер буфер, значит файл закончился,
		// буфер уменьшим
		buf = buf[:readSize]
		totalRead += readSize

		// Если задан лимит и мы его превысили, то уменьшим размер буфера,
		// и укажем что это последняя итерация
		if totalRead > int(limit) {
			tail := bufSize - (totalRead - int(limit))
			buf = buf[:tail]
			done = true
		}

		// Запишем прочитанный кусок
		_, outErr := outFile.Write(buf)
		if outErr != nil {
			return outErr
		}

		bar.Increment()
		// time.Sleep(500 * time.Millisecond) // для проверки прогресс бара
	}

	return nil
}
