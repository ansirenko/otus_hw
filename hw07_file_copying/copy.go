package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIncorrectInputData    = errors.New("offset or limit is incorrect")
)

func prepareSources(from, to string, offset int64) (*os.File, *os.File, error) {
	source, err := os.Open(from)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open file: %w", err)
	}

	destination, err := os.Create(to)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create file: %w, ", err)
	}

	_, err = source.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, nil, fmt.Errorf("can't use offset: %w", err)
	}

	return source, destination, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileSize, err := checkConditions(fromPath, offset, limit)
	if err != nil {
		return err
	}

	source, destination, err := prepareSources(fromPath, toPath, offset)
	if err != nil {
		return err
	}
	defer source.Close()
	defer destination.Close()

	readySize := fileSize - offset
	if readySize > limit && limit > 0 {
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
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("can't read file: %w", err)
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return fmt.Errorf("can't write to file: %w", err)
		}

		summ += n

		bar.currentValue <- summ
		if int64(summ) == readySize {
			break
		}
	}

	return nil
}

func checkConditions(from string, offset, limit int64) (int64, error) {
	if from == "" {
		return 0, ErrUnsupportedFile
	}

	if offset < 0 || limit < 0 {
		return 0, ErrIncorrectInputData
	}

	fileStat, err := os.Stat(from)
	if err != nil {
		return 0, fmt.Errorf("can't get file stats: %w", err)
	}

	if !fileStat.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}

	if offset > fileStat.Size() {
		return 0, ErrOffsetExceedsFileSize
	}

	return fileStat.Size(), nil
}
