package secure

import (
	"encoding/base64"
	"strings"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/argon2"
)

// CREDIT: The contents of this file were derived from https://github.com/alexedwards/argon2id.

type Argon2Params struct {
	// The amount of memory used by the algorithm (in kibibytes).
	Memory uint32

	// The number of iterations over the memory.
	Iterations uint32

	// The number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// Length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32

	// Length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

func (p *Argon2Params) toExternalParams() *argon2id.Params {
	return &argon2id.Params{
		Memory:      p.Memory,
		Iterations:  p.Iterations,
		Parallelism: p.Parallelism,
		SaltLength:  p.SaltLength,
		KeyLength:   p.KeyLength,
	}
}

func fromExternalParams(p *argon2id.Params) *Argon2Params {
	return &Argon2Params{
		Memory:      p.Memory,
		Iterations:  p.Iterations,
		Parallelism: p.Parallelism,
		SaltLength:  p.SaltLength,
		KeyLength:   p.KeyLength,
	}
}

func DeriveKey(password, hash string) (key []byte, err error) {
	params, salt, err := decodeHashWithoutKey(hash)
	if err != nil {
		return nil, err
	}

	key = argon2.IDKey([]byte(password), salt,
		params.Iterations, params.Memory, params.Parallelism, params.KeyLength)
	return key, nil
}

func CreateKeyAndHash(password string, params *Argon2Params) (key []byte, hash string, err error) {
	hash, err = argon2id.CreateHash(password, params.toExternalParams())
	if err != nil {
		return nil, "", err
	}

	parts := strings.Split(hash, "$")
	hash = strings.Join(parts[:len(parts)-1], "$")

	key, err = base64.RawStdEncoding.Strict().DecodeString(parts[len(parts)-1])
	if err != nil {
		return nil, "", err
	}

	return key, hash, nil
}

func decodeHashWithoutKey(hash string) (params *Argon2Params, salt []byte, err error) {
	hash += "$YnJvZGll"
	externalParams, salt, _, err := argon2id.DecodeHash(hash)
	if err != nil {
		return nil, nil, err
	}
	return fromExternalParams(externalParams), salt, nil
}
