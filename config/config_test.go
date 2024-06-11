package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Test case 1: Valid config file
	config, err := LoadConfig("../stubs/etc/serupmon/serupmon.hcl")
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Test case 2: Non-existent config file
	_, err = LoadConfig("../stubs/etc/serupmon/serupmon.hclx")
	assert.Error(t, err)
	assert.EqualError(t, err, "config file not found: ../stubs/etc/serupmon/serupmon.hclx")

	// Test case 3: Config file is a directory
	_, err = LoadConfig("../stubs/etc/serupmon")
	assert.Error(t, err)
	assert.EqualError(t, err, "config file is a directory: ../stubs/etc/serupmon")

	// Add more test cases as needed
}
