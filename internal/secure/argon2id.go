package secure

import (
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// CREDIT: The contents of this file were derived from https://github.com/alexedwards/argon2id.

var (
	// ErrInvalidHash is returned if the provided hash isn't in the expected format.
	ErrInvalidHash = errors.New("argon2id: hash is not in the correct format")

	// ErrIncompatibleVariant is returned if the provided hash was created using a unsupported variant of Argon2.
	// Currently only "argon2id" is supported by this package.
	ErrIncompatibleVariant = errors.New("argon2id: incompatible variant of argon2")

	// ErrIncompatibleVersion is returned if the provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")
)

type ArgonParams struct {
	// The amount of memory used by the algorithm (in kibibytes).
	Memory uint32

	// The number of iterations over the memory.
	Iterations uint32

	// The number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// Length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32
}

func DeriveKey(password, encodedParams string, keyLength uint32) ([]byte, error) {
	params, salt, err := DecodeArgonParams(encodedParams)
	if err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(password), salt,
		params.Iterations, params.Memory, params.Parallelism, keyLength)
	return key, nil
}

func CreateKey(password string, params *ArgonParams, keyLength uint32) ([]byte, string, error) {
	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return nil, "", err
	}
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	key := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, keyLength)

	encodedParams := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt)

	return key, encodedParams, nil
}

func DecodeArgonParams(hash string) (*ArgonParams, []byte, error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 5 {
		return nil, nil, ErrInvalidHash
	}

	if vals[1] != "argon2id" {
		return nil, nil, ErrIncompatibleVariant
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	params := &ArgonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	return params, salt, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
