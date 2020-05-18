package organize

import (
	cmap "github.com/orcaman/concurrent-map"
)

func Run(oauthToken string, foldername string, moveFileToFolder string) {
	var conf Configuration
	folders := cmap.New()

	conf = CreateConfiguration(oauthToken, "")
	GetFolderInformation(conf, "", 0, folders)
	if foldername == "" {
		foldername = "~/Downloads"
	}
	if moveFileToFolder == "" {
		moveFileToFolder = foldername
	}
	GoOrganizeFolder(foldername, folders, conf, moveFileToFolder)
}
