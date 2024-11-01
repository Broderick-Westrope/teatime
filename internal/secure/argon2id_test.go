package secure

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/argon2"
)

func TestDeriveKey(t *testing.T) {
	tt := map[string]struct {
		password      string
		encodedParams string
		keyLenght     uint32
		wantKey       []byte
	}{
		"key length 16": {
			password:      "pa$$word",
			encodedParams: "$argon2id$v=19$m=65536,t=1,p=10$c29tZV90ZXN0X3NhbHQ",
			keyLenght:     16,
			wantKey:       []byte{197, 165, 8, 93, 210, 140, 41, 207, 48, 145, 123, 18, 169, 209, 155, 127},
		},
		"key length 24": {
			password:      "pa$$word",
			encodedParams: "$argon2id$v=19$m=65536,t=1,p=10$c29tZV90ZXN0X3NhbHQ",
			keyLenght:     24,
			wantKey:       []byte{250, 127, 55, 191, 157, 157, 158, 167, 112, 237, 192, 78, 113, 49, 170, 62, 182, 54, 206, 187, 169, 103, 207, 232},
		},
		"key length 32": {
			password:      "pa$$word",
			encodedParams: "$argon2id$v=19$m=65536,t=1,p=10$c29tZV90ZXN0X3NhbHQ",
			keyLenght:     32,
			wantKey:       []byte{224, 155, 147, 58, 185, 138, 18, 132, 178, 70, 157, 131, 177, 128, 182, 123, 43, 195, 131, 39, 29, 214, 67, 196, 29, 234, 255, 166, 211, 9, 95, 191},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			key, err := DeriveKey(tc.password, tc.encodedParams, tc.keyLenght)

			require.NoError(t, err)
			assert.Equal(t, tc.wantKey, key)
		})
	}
}

func TestCreateKey(t *testing.T) {
	inParams := &ArgonParams{
		Memory:      65536,
		Iterations:  1,
		Parallelism: 10,
		SaltLength:  14,
	}

	key, hash, err := CreateKey("pa$$word", inParams, 16)
	require.NoError(t, err)

	params, salt, err := decodeArgonParams(hash)
	require.NoError(t, err)

	assert.Len(t, key, 16)
	assert.Len(t, salt, int(inParams.SaltLength))
	assert.EqualValues(t, inParams, params)
}

func TestDecodeHashWithoutKey(t *testing.T) {
	wantSalt := []byte("some_test_salt")
	inSalt := base64.RawStdEncoding.EncodeToString(wantSalt)

	inParams := &ArgonParams{
		Memory:      65536,
		Iterations:  1,
		Parallelism: 10,
		SaltLength:  uint32(len(wantSalt)),
	}

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		argon2.Version, inParams.Memory, inParams.Iterations, inParams.Parallelism, inSalt)

	params, salt, err := decodeArgonParams(hash)

	require.NoError(t, err)
	assert.EqualValues(t, inParams, params)
	assert.Equal(t, wantSalt, salt)
}

func TestDecryptionFlow(t *testing.T) {
	password := "pa$$word"
	encodedParams := "$argon2id$v=19$m=65536,t=1,p=10$RHAoBAEOFt+P/4bxPIY9RA"
	ciphertext := "n7H4HQ0vJjp68nyMI8/kdXzIoICD+ebtTGh5L+t9"

	key, err := DeriveKey(password, encodedParams, 24)
	require.NoError(t, err)

	decodedCiphertext, err := base64.RawStdEncoding.Strict().DecodeString(ciphertext)
	require.NoError(t, err)

	plaintextBytes, err := DecryptAESGCM(key, decodedCiphertext)
	require.NoError(t, err)

	assert.Equal(t, string(plaintextBytes), "[]")
}
