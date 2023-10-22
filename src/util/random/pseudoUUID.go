package random

import (
	"crypto/rand"
	"fmt"
)

// Note - NOT RFC4122 compliant
func PseudoUUID() (string, error) {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	// uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	uuid := fmt.Sprintf("%X", b[0:6])
	return uuid, nil
}
