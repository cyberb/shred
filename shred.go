package shred

import (
	"os"
)

type Shredder struct {
	random RandomGenerator
}

type RandomGenerator interface {
	Random() (byte, error)
}

func New(random RandomGenerator) *Shredder {
	return &Shredder{
		random: random,
	}
}

func (s *Shredder) Shred(path string) error {
	i := 0
	for i < 3 {
		err := s.WriteRandom(path)
		if err != nil {
			return err
		}
		i++
	}
	return os.Remove(path)
}

func (s *Shredder) WriteRandom(path string) error {
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
		err := s.writeRandomBytes(file, current)
		if err != nil {
			return err
		}
		current++
	}
	return nil
}

func (s *Shredder) writeRandomBytes(file *os.File, offset int64) error {
	actual := make([]byte, 1)
	_, err := file.ReadAt(actual, offset)
	if err != nil {
		return err
	}

	rnd, err := s.random.Random()
	if err != nil {
		return err
	}
	for actual[0] == rnd {
		rnd, err = s.random.Random()
		if err != nil {
			return err
		}
	}

	_, err = file.WriteAt([]byte{rnd}, offset)
	if err != nil {
		return err
	}
	return nil
}
