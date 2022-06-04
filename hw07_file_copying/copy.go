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
	// Проверим размер входного файла
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
	const bufSize = 512
	buf := make([]byte, bufSize)

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
		if limit > 0 && totalRead > int(limit) {
			tail := bufSize - (totalRead - int(limit))
			buf = buf[:tail]
			done = true
		}

		// Запишем прочитанный кусок
		_, outErr := outFile.Write(buf)
		if outErr != nil {
			return outErr
		}
	}

	return nil
}
