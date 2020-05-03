package organize

import (
	cmap "github.com/orcaman/concurrent-map"
)

func Run(configfile string, foldername string) {
	var conf Configuration
	folders := cmap.New()

	conf = CreateConfiguration(Read(configfile), "")
	GetFolderInformation(conf, "", 0, folders)
	if foldername == "" {
		foldername = "~/Downloads"
	}
	GoOrganizeFolder(foldername, folders, conf)
}
