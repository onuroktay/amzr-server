package util

import (
	"fmt"
	"crypto/rand"
	"strings"
)

// GetUUID unique uuid
func GetUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// GetEncryptionKey encrypt key
func GetEncryptionKey() []byte {
	return []byte("mytoken")
}

func CleanQuote(t string) string {
	t = strings.Replace(t, `"`, "", -1)
	t = strings.Replace(t, `'`, "", -1)

	return t
}
