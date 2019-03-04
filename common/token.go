package common

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ofonimefrancis/safeboda/common/must"
)

func GenerateRandomToken() string {
	token := make([]byte, 6)
	must.DoF(func() error {
		_, err := rand.Read(token)
		return err
	})
	return hex.EncodeToString(token)
}
