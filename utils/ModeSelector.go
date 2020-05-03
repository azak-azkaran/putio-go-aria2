package utils

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Configuration struct {
	Oauthtoken string
	Foldername string
	Mode       string
	Filter     string
	Url        string
}

const (
	ERROR_URL_MISSING        = "Error: Missing url in config file or ARIA2_ADDRESS is not set"
	ERROR_MODE_MISSING       = "Error: Missing mode in config file or ARIA2_MODE is not set"
	ERROR_TOKEN_MISSING      = "Error: Missing oauth token in config file or ARIA2_TOKEN is not set"
	ERROR_FOLDERNAME_MISSING = "Error: Missing foldername in config file or ARIA2_FOLDERNAME is not set"
	ERROR_WRONG_MODE         = "Error: Wrong mode was used, vaild modes: download,organize,d,o"
)

func GetArguments(filename string) (*Configuration, error) {
	viper.SetConfigName(filename)
	viper.SetEnvPrefix("aria2")
	viper.BindEnv("address")
	viper.BindEnv("oauth_token")
	viper.BindEnv("foldername")
	viper.BindEnv("mode")
	viper.BindEnv("filter")

	if _, err := os.Stat("config.yml"); err == nil {
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	var config Configuration

	if viper.InConfig("address") || viper.IsSet("address") {
		config.Url = viper.GetString("address")
	} else {
		Error.Fatalln(ERROR_URL_MISSING)
		return nil, errors.New(ERROR_URL_MISSING)
	}

	if viper.InConfig("mode") || viper.IsSet("mode") {
		modeString := viper.GetString("mode")
		if strings.TrimSpace(modeString) == "download" || strings.TrimSpace(modeString) == "d" {
			config.Mode = "d"
		} else if strings.TrimSpace(modeString) == "organize" || strings.TrimSpace(modeString) == "o" {
			config.Mode = "o"
		} else {
			Error.Fatalln(ERROR_WRONG_MODE)
			return nil, errors.New(ERROR_WRONG_MODE)
		}
	} else {
		Error.Fatalln(ERROR_MODE_MISSING)
		return nil, errors.New(ERROR_MODE_MISSING)
	}

	if viper.InConfig("oauth_token") || viper.IsSet("oauth_token") {
		config.Oauthtoken = viper.GetString("oauth_token")
	} else {
		Error.Fatalln(ERROR_TOKEN_MISSING)
		return nil, errors.New(ERROR_TOKEN_MISSING)
	}

	if config.Mode == "o" {
		if viper.InConfig("foldername") || viper.IsSet("foldername") {
			config.Foldername = viper.GetString("foldername")
		} else {
			Error.Fatalln(ERROR_FOLDERNAME_MISSING)
			return nil, errors.New(ERROR_FOLDERNAME_MISSING)
		}
	}

	if viper.InConfig("filter") || viper.IsSet("filter") {
		config.Foldername = viper.GetString("filter")
	}
	return &config, nil
}
