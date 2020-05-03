package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	ARIA2_ADDRESS = "127.0.0.1:6800"
	ARIA2_MODE    = "download"
	ARIA2_TOKEN   = "supersecrettoken"
)

func TestEnvironmentVariables(t *testing.T) {
	Init(os.Stdout, os.Stdout, os.Stdout)

	os.Setenv("ARIA2_ADDRESS", ARIA2_ADDRESS)
	os.Setenv("ARIA2_MODE", ARIA2_MODE)
	os.Setenv("ARIA2_OAUTH_TOKEN", ARIA2_TOKEN)

	conf, err := GetArguments("./wrong.conf")
	require.NoError(t, err)

	assert.EqualValues(t, conf.Url, ARIA2_ADDRESS)
	assert.EqualValues(t, conf.Mode, "d")
	assert.EqualValues(t, conf.Oauthtoken, ARIA2_TOKEN)
}
