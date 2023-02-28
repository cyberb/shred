package shred

import "crypto/rand"

type CryptoRandomGenerator struct{}

func (r *CryptoRandomGenerator) Random() (byte, error) {
	buf := make([]byte, 1)
	_, err := rand.Read(buf)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
