package utils

import (
	"crypto/rand"
	"math/big"
)

const projectCodeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const projectCodeLength = 5

// GenerateProjectCode produces a random 5 character code for projects.
func GenerateProjectCode() (string, error) {
	result := make([]byte, projectCodeLength)
	alphabetLength := big.NewInt(int64(len(projectCodeAlphabet)))

	for i := 0; i < projectCodeLength; i++ {
		n, err := rand.Int(rand.Reader, alphabetLength)
		if err != nil {
			return "", err
		}
		result[i] = projectCodeAlphabet[n.Int64()]
	}

	return string(result), nil
}
