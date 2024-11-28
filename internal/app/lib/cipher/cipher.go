package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

type Cipher struct {
	aesgcm cipher.AEAD
	key    []byte
}

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func New() (*Cipher, error) {
	key, err := GenerateRandom(2 * aes.BlockSize)
	if err != nil {
		return nil, err
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	return &Cipher{
		aesgcm: aesgcm,
		key:    key,
	}, nil
}

func (c *Cipher) Sile(src string) string {
	nonce := c.key[len(c.key)-c.aesgcm.NonceSize():]
	dst := c.aesgcm.Seal(nil, nonce, []byte(src), nil)
	return hex.EncodeToString(dst)
}

func (c *Cipher) Open(data string) (string, error) {
	nonce := c.key[len(c.key)-c.aesgcm.NonceSize():]
	dst, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	src, err := c.aesgcm.Open(nil, nonce, dst, nil)
	if err != nil {
		return "", err
	}
	return string(src), nil
}
