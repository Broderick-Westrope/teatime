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
	password := "pa$$word"
	hash := "$argon2id$v=19$m=65536,t=1,p=10$c29tZV90ZXN0X3NhbHQ"
	wantKey := []byte{239, 8, 160, 14, 108, 197}

	key, err := DeriveKey(password, hash)

	require.NoError(t, err)
	assert.Equal(t, wantKey, key)
}

func TestCreateKeyAndHash(t *testing.T) {
	inParams := &Argon2Params{
		Memory:      65536,
		Iterations:  1,
		Parallelism: 10,
		SaltLength:  14,
		KeyLength:   8,
	}

	key, hash, err := CreateKeyAndHash("pa$$word", inParams)
	require.NoError(t, err)

	params, salt, err := decodeHashWithoutKey(hash)
	require.NoError(t, err)

	assert.Len(t, key, int(inParams.KeyLength))
	assert.Len(t, salt, int(inParams.SaltLength))
	assert.Equal(t, inParams.Memory, params.Memory)
	assert.Equal(t, inParams.Iterations, params.Iterations)
	assert.Equal(t, inParams.Parallelism, params.Parallelism)
	assert.Equal(t, inParams.SaltLength, params.SaltLength)
	// NOTE: we don't verify key length.
	// 	this field is only accurate when creating a key since the stored hash contains no key to infer the length of.
}

func TestDecodeHashWithoutKey(t *testing.T) {
	wantSalt := []byte("some_test_salt")
	inSalt := base64.RawStdEncoding.EncodeToString(wantSalt)

	inParams := &Argon2Params{
		Memory:      65536,
		Iterations:  1,
		Parallelism: 10,
		SaltLength:  uint32(len(wantSalt)),
		KeyLength:   6,
	}

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		argon2.Version, inParams.Memory, inParams.Iterations, inParams.Parallelism, inSalt)

	params, salt, err := decodeHashWithoutKey(hash)

	require.NoError(t, err)
	assert.EqualValues(t, inParams, params)
	assert.Equal(t, wantSalt, salt)
}
