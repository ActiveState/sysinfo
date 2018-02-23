package sysinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOS(t *testing.T) {
	assert.NotEqual(t, UnknownOs, OS())
}

func TestOSVersion(t *testing.T) {
	version, err := OSVersion()
	assert.Nil(t, err, "Determined OS version")
	assert.NotEmpty(t, version.Version)
	assert.NotEqual(t, 0, version.Major)
	assert.NotEqual(t, 0, version.Minor)
	assert.NotEmpty(t, version.Name)
}

func TestArchitecture(t *testing.T) {
	assert.NotEqual(t, UnknownArch, Architecture())
}

func TestLibc(t *testing.T) {
	libc, err := Libc()
	assert.Nil(t, err, "Determined Libc version")
	assert.NotEqual(t, UnknownLibc, libc)
}

func TestCompiler(t *testing.T) {
	compilers, err := Compilers()
	assert.Nil(t, err, "Determined system compilers")
	assert.NotEqual(t, 0, len(compilers), "More than one compiler was found")
}
