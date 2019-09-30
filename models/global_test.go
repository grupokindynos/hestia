package models

import (
	"github.com/grupokindynos/common/hestia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigModel_UpdateConfigData(t *testing.T) {
	err := configModel.UpdateConfigData(TestConfigData)
	assert.Nil(t, err)
}

func TestConfigModel_GetConfigData(t *testing.T) {
	configData, err := configModel.GetConfigData()
	assert.Nil(t, err)
	assert.IsType(t, hestia.Config{}, configData)
	assert.Equal(t, TestConfigData, configData)
}
