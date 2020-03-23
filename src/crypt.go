package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

func createPBKDF(key string, salt []byte) []byte {
	return pbkdf2.Key([]byte(key), salt, 4096, 32, sha1.New)
}

// RandomSalt can be used to get a randomSalt for use on calls to Encrypt and Decrypt functions
func RandomSalt(size int) (salt []byte, err error) {
	salt = make([]byte, size)
	_, err = rand.Read(salt)
	if err != nil {
		return salt, err
	}
	return salt, nil
}

// Encrypt takes some data and creates cipherText using AES
// A passphrase and a random salt must be provided along the data
// The AES Cipher is created using a key derived from the passphrase and the salt value using the standard PBKDF2 go library
func Encrypt(data []byte, passphrase string, salt []byte) (cipherText []byte, err error) {
	block, _ := aes.NewCipher(createPBKDF(passphrase, salt))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherText = gcm.Seal(nonce, nonce, data, nil)
	return cipherText, nil
}

// Decrypt takes data encrypted with Encrypt function and returns the original decrypted data.
// The original passphrase and salt value used to encrypt the data must be provided.
func Decrypt(data []byte, passphrase string, salt []byte) (plainText []byte, err error) {
	key := []byte(createPBKDF(passphrase, salt))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plainText, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
