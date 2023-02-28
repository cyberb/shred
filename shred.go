package shred

import (
	"crypto/rand"
	"os"
)

func Shred(path string) error {
	i := 0
	for i < 3 {
		err := WriteRandom(path)
		if err != nil {
			return err
		}
		i++
	}
	return os.Remove(path)
}

func WriteRandom(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	current := int64(0)
	for current < size {
		err := writeRandomBytes(file, current)
		if err != nil {
			return err
		}
		current++
	}
	return nil
}

func writeRandomBytes(file *os.File, offset int64) error {
	actual := make([]byte, 1)
	_, err := file.ReadAt(actual, offset)
	if err != nil {
		return err
	}

	buf, err := random(1)
	if err != nil {
		return err
	}
	for actual[0] == buf[0] {
		buf, err = random(1)
		if err != nil {
			return err
		}
	}

	_, err = file.WriteAt(buf, offset)
	if err != nil {
		return err
	}
	return nil
}

func random(len int64) ([]byte, error) {
	buf := make([]byte, len)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
