package t

import (
	"os"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	code := m.Run()
	s.Shutdown(context.TODO())
	os.Exit(code)
}

func TestGetAllPromo(t *testing.T) {
	ds := getDataStoreSession()
	_, err := ds.GetAllPromos("2")
	assert.NoError(t, err)
}

func TestGetAllActivePromos(t *testing.T) {
	ds := getDataStoreSession()
	_, err := ds.GetAllActivePromos("5")
	assert.NoError(t, err)
}
