package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSetting
func TestSetting(t *testing.T) {
	err := SetSetting("", "content")
	assert.NotNil(t, err, "Fail to set wrong setting")

	err = SetSetting("a", "content")
	assert.Nil(t, err, "Fail to set correct setting")

	err = SetSetting("a", "new content")
	assert.Nil(t, err, "Fail to set correct setting")

	_, err = GetSetting("")
	assert.NotNil(t, err, "Fail to get empty key")

	_, err = GetSetting("b")
	assert.NotNil(t, err, "Fail to get non exist key")

	v, err := GetSetting("a")
	assert.Nil(t, err, "Fail to get exist key")
	assert.Equal(t, v, "new content", "Fail to get correct value")
}
