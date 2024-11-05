package secure_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"testing"

	"github.com/Broderick-Westrope/teatime/internal/secure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncryptAESGCM tests the EncryptAESGCM function independently.
func TestEncryptAESGCM(t *testing.T) {
	key := []byte("thisisa16bytekey") // 16-byte key for AES-128
	plaintext := []byte("Test plaintext data for encryption.")

	ciphertext, err := secure.EncryptAESGCM(key, plaintext)
	require.NoError(t, err)

	require.NotEmpty(t, ciphertext)

	block, err := aes.NewCipher(key)
	require.NoError(t, err)
	aesgcm, err := cipher.NewGCM(block)
	require.NoError(t, err)
	nonceSize := aesgcm.NonceSize()

	// Check that the ciphertext length is nonce + plaintext + overhead.
	expectedLength := nonceSize + len(plaintext) + aesgcm.Overhead()
	assert.Len(t, ciphertext, expectedLength)

	assert.NotContains(t, ciphertext, plaintext)
}

// TestDecryptAESGCM tests the DecryptAESGCM function independently.
func TestDecryptAESGCM(t *testing.T) {
	key := []byte("thisisa16bytekey") // 16-byte key for AES-128
	nonce := []byte("unique_nonce")   // 12-byte nonce for AES-GCM
	plaintext := []byte("Test plaintext data for decryption.")

	block, err := aes.NewCipher(key)
	require.NoError(t, err)
	aesgcm, err := cipher.NewGCM(block)
	require.NoError(t, err)
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Prepend the nonce to the ciphertext to match the expected input format.
	nonce = append(nonce, ciphertext...)

	decryptedPlaintext, err := secure.DecryptAESGCM(key, nonce)
	require.NoError(t, err)

	assert.EqualValues(t, decryptedPlaintext, plaintext)
}

// TestEncryptDecryptAESGCM checks if the encryption and decryption process works correctly.
func TestEncryptDecryptAESGCM(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := []byte("This is a test message.")

	ciphertext, err := secure.EncryptAESGCM(key, plaintext)
	require.NoError(t, err)

	// Decrypt the ciphertext.
	decryptedText, err := secure.DecryptAESGCM(key, ciphertext)
	require.NoError(t, err)

	assert.EqualValues(t, plaintext, decryptedText)
}

// TestDecryptAESGCMWithWrongKey ensures decryption fails when using the wrong key.
func TestDecryptAESGCMWithWrongKey(t *testing.T) {
	correctKey := make([]byte, 32)
	_, err := rand.Read(correctKey)
	require.NoError(t, err)

	wrongKey := make([]byte, 32)
	_, err = rand.Read(wrongKey)
	require.NoError(t, err)

	plaintext := []byte("This is a test message.")

	ciphertext, err := secure.EncryptAESGCM(correctKey, plaintext)
	require.NoError(t, err)

	_, err = secure.DecryptAESGCM(wrongKey, ciphertext)
	assert.Error(t, err)
}

// TestDecryptAESGCMWithCorruptedCiphertext checks decryption failure when ciphertext is tampered with.
func TestDecryptAESGCMWithCorruptedCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := []byte("This is a test message.")

	ciphertext, err := secure.EncryptAESGCM(key, plaintext)
	require.NoError(t, err)

	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err = secure.DecryptAESGCM(key, ciphertext)
	assert.Error(t, err)
}

// TestEncryptAESGCMWithInvalidKey checks encryption failure when using an invalid key size.
func TestEncryptAESGCMWithInvalidKey(t *testing.T) {
	key := make([]byte, 20)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := []byte("This is a test message.")

	_, err = secure.EncryptAESGCM(key, plaintext)
	assert.Error(t, err)
}

// TestDecryptAESGCMWithShortCiphertext ensures decryption fails when ciphertext is too short.
func TestDecryptAESGCMWithShortCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	ciphertext := make([]byte, 5)

	_, err = secure.DecryptAESGCM(key, ciphertext)
	assert.Error(t, err)
}
