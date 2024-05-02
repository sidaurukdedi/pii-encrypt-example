package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type Crypto interface {
	Encrypt(plainText string) (cipherText []byte, err error)
	Decrypt(encryptedTextBytes []byte) (plainText []byte, err error)
	Hash(s string) []byte
}

// AESGCM concrete struct of Crypto PII Use Case.
type AESGCM struct {
	secret string
	pepper string
}

// NewAESGCM is constructor.
func NewAESGCM(secret string, pepper string) Crypto {
	return &AESGCM{
		secret: secret,
		pepper: pepper,
	}
}

// Encrypt returns encrypted string.
func (a AESGCM) Encrypt(plainText string) (cipherText []byte, err error) {
	// The key argument should be the AES key, either 16 or 32 bytes to select AES-128 or AES-256.
	// https://pkg.go.dev/crypto/aes#pkg-overview
	key := []byte(a.secret)
	bPlainText := []byte(plainText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText = aesgcm.Seal(nil, nonce, bPlainText, nil)
	cipherText = append(nonce[:], cipherText[:]...)

	return
}

// Decrypt returns plaintext of encrypted string.
func (a AESGCM) Decrypt(encryptedTextBytes []byte) (plainText []byte, err error) {
	// The key argument should be the AES key, either 16 or 32 bytes to select AES-128 or AES-256.
	key := []byte(a.secret)
	nonce := encryptedTextBytes[:12]
	encryptedTextBytes = encryptedTextBytes[12:]

	if len(encryptedTextBytes) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return plainText, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return plainText, err
	}

	plainText, err = aesGcm.Open(nil, nonce, encryptedTextBytes, nil)
	if err != nil {
		return plainText, err
	}

	return plainText, nil
}

func (a AESGCM) Hash(s string) []byte {
	return a.CreateHashByte([]byte(s), []byte(a.pepper))
}

// CreateHashByte from string, using sha256.
func (a AESGCM) CreateHashByte(key ...[]byte) []byte {
	var data []byte
	for _, s := range key {
		data = append(data[:], s[:]...)
	}
	hash := sha256.New()

	hash.Write(data)
	return hash.Sum(nil)
}
