package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigInvalidFile(t *testing.T) {
	docModel, err := loadSpec("noExist.yml")
	assert.Error(t, err, "Expected Error loading spec")
	assert.Nil(t, docModel, "Document model should be nil")
}

func TestLoadConfigInvalidFormat(t *testing.T) {
	docModel, err := loadSpec("./test/spec/invalid.yml")
	assert.Error(t, err, "Expected Error loading spec")
	assert.Nil(t, docModel, "Document model should be nil")
}

func TestLoadConfigMalformedV2Model(t *testing.T) {
	docModel, err := loadSpec("./test/spec/test.yml")
	assert.Error(t, err, "Expected Error loading spec")
	assert.Nil(t, docModel, "Document model should be nil")
}

func TestLoadConfigInvalidInvalidV3Model(t *testing.T) {
	docModel, err := loadSpec("./test/spec/test2.yml")
	assert.Error(t, err, "Expected Error loading spec")
	assert.Nil(t, docModel, "Document model should be nil")
}

func TestLoadConfigValidFile(t *testing.T) {
	docModel, err := loadSpec("api.yml")
	assert.NoError(t, err, "Error loading config")
	assert.NotNil(t, docModel, "Document model should not be nil")
}

func TestPrintSpecValid(t *testing.T) {
	err := printSpec("api.yml")
	assert.NoError(t, err, "Error loading spec")
}

func TestPrintSpecInvalid(t *testing.T) {
	err := printSpec("NoExist.yml")
	assert.Error(t, err, "Error loading spec")
}
