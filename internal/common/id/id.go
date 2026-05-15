package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func New() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err == nil {
		return hex.EncodeToString(bytes)
	}
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
}
