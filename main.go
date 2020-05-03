package main

import (
	"github.com/azak-azkaran/putio-go-aria2/aria2downloader"
	"github.com/azak-azkaran/putio-go-aria2/organize"
	util "github.com/azak-azkaran/putio-go-aria2/utils"
	"os"
)

func main() {
	util.Init(os.Stdout, os.Stdout, os.Stderr)

	config, err := util.GetArguments("config")
	if err != nil {
		panic(err)
	}

	if config.Mode == "d" {
		aria2downloader.Run(config.Oauthtoken, config.Filter, config.Url)

	} else if config.Mode == "o" {
		organize.Run(config.Oauthtoken, config.Foldername)

	} else {
		panic("mode not detected")
	}
	os.Exit(0)
}
