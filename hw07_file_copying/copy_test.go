package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {

	var sourceFile = "./testdata/input.txt"

	t.Run("simple copy", func(t *testing.T) {
		stdin, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		os.Stdin,_ = os.Open(stdin.Name())
		expected, err := ioutil.ReadFile(sourceFile)
		require.NoError(t, err, "can't read test file")

		resultFile, err := ioutil.TempFile("", "")
		require.NoError(t, err, "can't create temporary file")

		err = Copy(sourceFile, resultFile.Name(), 0, 0)
		require.NoError(t, err, "unexpected copy file error")

		actual, err := ioutil.ReadFile(resultFile.Name())
		require.NoError(t, err, "can't read result file")
		require.Equal(t, expected, actual, "copied file is incorrect")
	})
}
