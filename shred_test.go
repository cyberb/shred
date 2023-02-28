package shred

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestWriteRandom(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := path.Join(tempDir, "file")
	file, err := os.Create(tempFile)
	assert.NoError(t, err)
	in := []byte{1, 2, 3}
	_, err = file.Write(in)
	assert.NoError(t, err)

	_ = file.Close()

	err = WriteRandom(tempFile)
	assert.NoError(t, err)

	file, err = os.OpenFile(tempFile, os.O_RDONLY, 0666)
	assert.NoError(t, err)
	offset := 0
	for offset < 3 {
		out := make([]byte, 1)
		_, err = file.ReadAt(out, int64(offset))
		assert.NoError(t, err)
		assert.NotEqual(t, in[offset], out[0])
		offset++
	}
	_ = file.Close()

}

func TestShred(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := path.Join(tempDir, "file")
	file, err := os.Create(tempFile)
	assert.NoError(t, err)
	in := []byte{1, 2, 3}
	_, err = file.Write(in)
	assert.NoError(t, err)

	_ = file.Close()

	err = Shred(tempFile)
	assert.NoError(t, err)

	_, err = os.Stat(tempFile)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestShred_NotFile(t *testing.T) {
	tempDir := t.TempDir()
	err := Shred(tempDir)
	assert.Error(t, err)

}
