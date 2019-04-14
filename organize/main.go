package organize

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	cmap "github.com/orcaman/concurrent-map"
	"os"
)

func Run(configfile string, foldername string) {
	var conf Configuration
	folders := cmap.New()

	conf = CreateConfiguration(Read(configfile), "")
	GetFolderInformation(conf, "", 0, folders)
	if foldername == "" {
		foldername = "~/Downloads"
	}
	OrganizeFolder(foldername, folders, conf)
}

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	configfile, _, foldername, err := utils.GetArguments()
	if err != nil {
		panic(err)
	}
	Run(configfile, foldername)
}
