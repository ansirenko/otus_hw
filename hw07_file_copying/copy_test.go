package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {

	var sourceFile = "./testdata/input.txt"

	t.Run("simple copy", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 0, 0)
		require.NoError(t, err, "unexpected copy file error")

		expected, err := ioutil.ReadFile(sourceFile)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("copy with incorrect offset", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		fileStat, err := os.Stat(sourceFile)
		require.NoError(t, err, "can't get testing file stats")

		err = Copy(sourceFile, resultFile.Name(), fileStat.Size()+1, 0)
		require.Error(t, err, "unexpected copy file error")
	})

	t.Run("copy with incorrect limit", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 0, -1)
		require.Error(t, err, "run with incorrect limit")
	})

	t.Run("limit bigger than file size", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		fileStat, err := os.Stat(sourceFile)
		require.NoError(t, err, "can't get testing file stats")

		err = Copy(sourceFile, resultFile.Name(), 0, fileStat.Size()+1)
		require.NoError(t, err, "unexpected copy file error")

		expected, err := ioutil.ReadFile(sourceFile)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("copy with offset", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 100, 0)
		require.NoError(t, err, "unexpected copy file error")

		resultFilePath := "./testdata/out_offset100_limit0.txt"

		expected, err := ioutil.ReadFile(resultFilePath)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("copy with limit 10", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 0, 10)
		require.NoError(t, err, "unexpected copy file error")

		resultFilePath := "./testdata/out_offset0_limit10.txt"

		expected, err := ioutil.ReadFile(resultFilePath)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("copy with limit 1000", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 0, 1000)
		require.NoError(t, err, "unexpected copy file error")

		resultFilePath := "./testdata/out_offset0_limit1000.txt"

		expected, err := ioutil.ReadFile(resultFilePath)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("copy with limit 10000", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 0, 100000)
		require.NoError(t, err, "unexpected copy file error")

		resultFilePath := "./testdata/out_offset0_limit10000.txt"

		expected, err := ioutil.ReadFile(resultFilePath)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("copy with offset 100 and limit 1000", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 100, 1000)
		require.NoError(t, err, "unexpected copy file error")

		resultFilePath := "./testdata/out_offset100_limit1000.txt"

		expected, err := ioutil.ReadFile(resultFilePath)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("copy with offset 6000 and limit 1000", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 6000, 1000)
		require.NoError(t, err, "unexpected copy file error")

		resultFilePath := "./testdata/out_offset6000_limit1000.txt"

		expected, err := ioutil.ReadFile(resultFilePath)
		require.NoError(t, err, "can't read test file")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})

	t.Run("incorrect filename", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy("./testdata/non-existentFile", resultFile.Name(), 0, 0)
		require.Error(t, err, "run with non existed file")
	})

	t.Run("unsupported filename", func(t *testing.T) {
		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy("/dev/urandom", resultFile.Name(), 0, 0)
		require.Error(t, err, "run with non existed file")
	})

	t.Run("incorrect destination", func(t *testing.T) {
		err := Copy(sourceFile, "./", 0, 0)
		require.Error(t, err, "run with incorrect destination")
	})
}
