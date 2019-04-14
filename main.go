package main

import (
	"github.com/azak-azkaran/putio-go-aria2/aria2downloader"
	"github.com/azak-azkaran/putio-go-aria2/organize"
	util "github.com/azak-azkaran/putio-go-aria2/utils"
	"os"
	"strings"
)

func main() {
	util.Init(os.Stdout, os.Stdout, os.Stderr)

	configfile, mode, add, err := util.GetArguments()
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(mode) == "d" || strings.TrimSpace(mode) == "downloader" {
		aria2downloader.Run(configfile, add)

	} else if strings.TrimSpace(mode) == "o" || strings.TrimSpace(mode) == "organize" {
		organize.Run(configfile, add)

	} else {
		panic("mode not detected")
	}
}
