package signup

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptStringValue(t *testing.T) {
	originalString := "thisOriginalString"
	encryptKey, _ := GenerateEncryptKey()
	encryptedString, _ := EncryptStringValue(originalString, encryptKey)

	assert.NotEqual(t, encryptKey, "", "encryptKey should not be an empty string")
	assert.NotEqual(t, encryptedString, originalString, "encrypted string should not equal original string")
	assert.Equal(t, reflect.TypeOf(encryptedString).Kind(), reflect.String, "encryptedString should be a string")
	assert.NotEqual(t, len(encryptedString), len(originalString), "encrypted string and original string should be different lengths")
}

func assertFailedEncrypt(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("EncryptStringValueWithEmptyKey failed: No panic thrown")
		}
	}()
	f()
}

func TestEncryptStringValueWithEmptyKey(t *testing.T) {
	originalString := "thisOriginalString"
	encryptKey := ""

	assertFailedEncrypt(t, func() { EncryptStringValue(originalString, encryptKey) })
}

func TestGenerateRandomId(t *testing.T) {
	resultID := GenerateRandomId()

	assert.NotEqual(t, resultID, "", "resultID should be a string")
}

func TestDecryptStringValue(t *testing.T) {
	originalString := "thisOriginalString"
	encryptKey, _ := GenerateEncryptKey()
	encryptedString, _ := EncryptStringValue(originalString, encryptKey)
	decryptedString, _ := DecryptStringValue(encryptedString, encryptKey)

	assert.NotEqual(t, encryptKey, "", "encryptKey should not be an empty string")
	assert.NotEqual(t, encryptedString, decryptedString, "encrypted string should not equal decrypted string")
	assert.Equal(t, originalString, decryptedString, "decrypted string should be the original string")
	assert.Equal(t, reflect.TypeOf(decryptedString).Kind(), reflect.String, "decrypted string should be a string")
}
