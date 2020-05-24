package main

import (
	"os"
	"testing"

	utils "github.com/azak-azkaran/putio-go-aria2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ARIA2_ADDRESS            = "http://127.0.0.1:6800/jsonrpc"
	ARIA2_ADDRESS_CONFIGFILE = "http://127.0.0.1:7000/jsonrpc"
	ARIA2_MODE               = "download"
	ARIA2_TOKEN              = "supersecrettoken"
)

func TestConfigFile(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)
	utils.Info.Println("testing: TestConfigFile")

	conf, err := GetArguments("../testdata/test_config.yml")
	require.NoError(t, err)

	assert.NotEqual(t, conf.Url, ARIA2_ADDRESS)
	assert.EqualValues(t, conf.Url, ARIA2_ADDRESS_CONFIGFILE)
	assert.EqualValues(t, conf.Mode, "d")
	assert.EqualValues(t, conf.Oauthtoken, ARIA2_TOKEN)
}
func TestEnvironmentVariables(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)
	utils.Info.Println("testing: TestEnvironmentVariables")

	os.Setenv("ARIA2_ADDRESS", "127.0.0.1:6800")
	os.Setenv("ARIA2_MODE", ARIA2_MODE)
	os.Setenv("ARIA2_OAUTH_TOKEN", ARIA2_TOKEN)

	conf, err := GetArguments("")
	require.NoError(t, err)

	assert.EqualValues(t, conf.Url, ARIA2_ADDRESS)
	assert.NotEqual(t, conf.Url, ARIA2_ADDRESS_CONFIGFILE)
	assert.EqualValues(t, conf.Mode, "d")
	assert.EqualValues(t, conf.Oauthtoken, ARIA2_TOKEN)
}
