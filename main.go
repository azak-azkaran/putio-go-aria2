package main

import (
	"os"

	"github.com/azak-azkaran/putio-go-aria2/aria2downloader"
	"github.com/azak-azkaran/putio-go-aria2/organize"
	util "github.com/azak-azkaran/putio-go-aria2/utils"
)

func main() {
	util.Init(os.Stdout, os.Stdout, os.Stderr)

	config, err := GetArguments("./config.yml")
	if err != nil {
		panic(err)
	}

	if config.Mode == "d" {
		aria2downloader.Run(config.Oauthtoken, config.Filter, config.Url)

	} else if config.Mode == "o" {
		organize.Run(config.Oauthtoken, config.Foldername, config.Destination)

	} else {
		panic("mode not detected")
	}
	os.Exit(0)
}
