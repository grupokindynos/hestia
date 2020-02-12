package bitcou

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	_ = godotenv.Load()
}

func TestProviders(t *testing.T) {
	bitcouService := InitService()
	prodProv, err := bitcouService.GetProviders(false)
	assert.Nil(t, err)
	assert.NotNil(t, prodProv)

	devProv, err := bitcouService.GetProviders(true)
	assert.Nil(t, err)
	assert.NotNil(t, devProv)

	assert.NotEqual(t, devProv, prodProv)
}