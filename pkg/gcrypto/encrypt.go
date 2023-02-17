package gcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
)

func PasswordToSha256(password string) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hash.Sum(nil)
}

func Encryption(key []byte, privateKey []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	return encrypt(block, privateKey)
}

func Decryption(key []byte, cipherText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	return decrypt(block, cipherText)
}

func encrypt(b cipher.Block, plaintext []byte) []byte {
	gcm, err := cipher.NewGCM(b)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func decrypt(b cipher.Block, ciphertext []byte) []byte {
	gcm, err := cipher.NewGCM(b)
	if err != nil {
		log.Fatal(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Fatal("len(ciphertext) < nonceSize")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	return plaintext
}
