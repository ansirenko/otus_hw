package main

import (
	"errors"
	"io"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrUnsupportedFile
	}

	fileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !fileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = source.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	readySize := (fileStat.Size() - offset)
	if readySize > limit && limit > 0{
		readySize = limit
	}

	buffSize := readySize / 100
	if buffSize == 0 {
		buffSize = 1
	}

	bar, err := NewProgressBar(int(readySize))
	if err != nil {
		return err
	}

	buf := make([]byte, buffSize)

	summ := 0

	bar.Play()
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		//limit -= int64(n)
		//if limit <= 0 {
		//	n -= int(limit)
		//}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}

		summ += n
		time.Sleep(100*time.Millisecond)

		bar.currentValue <- summ
		if int64(summ) == readySize {
			break
		}
	}

	return nil
}

