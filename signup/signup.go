package signup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	uuid "github.com/satori/go.uuid"
)

type Signup struct {
	ActionID string   `json:"action_id"`
	Signals  []string `json:"signals"`
	Name     string   `json:"name"`
}

// Encrypt the signup values to send to the virtualorb backend
func EncryptStringValue(value string, encryptKey string) (string, error) {
	myKey, _ := hex.DecodeString(encryptKey)
	plaintext := []byte(value)

	block, err := aes.NewCipher(myKey)
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", cipherText), nil
}

func DecryptStringValue(value string, encryptKey string) (string, error) {
	myKey, _ := hex.DecodeString(encryptKey)
	enc, _ := hex.DecodeString(value)

	block, err := aes.NewCipher(myKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, cipherText := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}

func GenerateRandomId() string {
	return uuid.NewV4().String()
}

func GenerateEncryptKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	encryptKey := hex.EncodeToString(bytes)
	return encryptKey, nil
}
