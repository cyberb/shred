package shred

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

type GeneratorStub struct {
	values  []byte
	current int
}

func (g *GeneratorStub) Random() (byte, error) {
	if g.current > len(g.values)-1 {
		return 0, fmt.Errorf("unexpected call to Random")
	}
	value := g.values[g.current]
	g.current++
	return value, nil
}

func TestWriteRandom(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := path.Join(tempDir, "file")
	file, err := os.Create(tempFile)
	assert.NoError(t, err)
	in := []byte{1, 2, 3}
	_, err = file.Write(in)
	assert.NoError(t, err)

	_ = file.Close()

	random := &GeneratorStub{values: []byte{2, 3, 4}}
	shredder := New(random)
	err = shredder.WriteRandom(tempFile)
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

func TestWriteRandom_SameByte_TryNextRandom(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := path.Join(tempDir, "file")
	file, err := os.Create(tempFile)
	assert.NoError(t, err)
	in := []byte{1}
	_, err = file.Write(in)
	assert.NoError(t, err)

	_ = file.Close()

	random := &GeneratorStub{values: []byte{1, 2}}
	shredder := New(random)
	err = shredder.WriteRandom(tempFile)
	assert.NoError(t, err)
	assert.Equal(t, 2, random.current)

}

func TestShred_ThreeTimes(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := path.Join(tempDir, "file")
	file, err := os.Create(tempFile)
	assert.NoError(t, err)
	in := []byte{1}
	_, err = file.Write(in)
	assert.NoError(t, err)

	_ = file.Close()

	random := &GeneratorStub{values: []byte{2, 3, 4}}
	shredder := New(random)
	err = shredder.Shred(tempFile)
	assert.NoError(t, err)
	assert.Equal(t, 3, random.current)

	_, err = os.Stat(tempFile)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestShred_NotFile(t *testing.T) {
	tempDir := t.TempDir()
	random := &GeneratorStub{values: []byte{}}
	shredder := New(random)
	err := shredder.Shred(tempDir)
	assert.Error(t, err)

}
