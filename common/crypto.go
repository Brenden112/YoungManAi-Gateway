package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHMACWithKey(key []byte, data string) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateHMAC(data string) string {
	h := hmac.New(sha256.New, []byte(CryptoSecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Password2Hash(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func ValidatePasswordAndHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// deriveAESKey derives a 32-byte AES-256 key from CryptoSecret via SHA-256.
func deriveAESKey() []byte {
	h := sha256.Sum256([]byte(CryptoSecret))
	return h[:]
}

// EncryptAES encrypts plaintext using AES-256-GCM with a key derived from CryptoSecret.
// Output is base64-encoded (nonce || ciphertext || tag).
func EncryptAES(plaintext string) (string, error) {
	if strings.TrimSpace(CryptoSecret) == "" {
		return "", fmt.Errorf("EncryptAES: crypto secret is not configured")
	}
	block, err := aes.NewCipher(deriveAESKey())
	if err != nil {
		return "", fmt.Errorf("EncryptAES: new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("EncryptAES: new GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("EncryptAES: generate nonce: %w", err)
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// DecryptAES decrypts a base64-encoded AES-256-GCM ciphertext produced by EncryptAES.
func DecryptAES(encoded string) (string, error) {
	if strings.TrimSpace(CryptoSecret) == "" {
		return "", fmt.Errorf("DecryptAES: crypto secret is not configured")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("DecryptAES: base64 decode: %w", err)
	}
	block, err := aes.NewCipher(deriveAESKey())
	if err != nil {
		return "", fmt.Errorf("DecryptAES: new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("DecryptAES: new GCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("DecryptAES: ciphertext too short")
	}
	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("DecryptAES: decrypt: %w", err)
	}
	return string(plaintext), nil
}
