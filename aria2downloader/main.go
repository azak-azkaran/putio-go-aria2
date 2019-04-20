package aria2downloader

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	cmap "github.com/orcaman/concurrent-map"
	"os"
)

func Run(configfile string, filter string) {
	url := "http://localhost:6800/jsonrpc"
	answers := cmap.New()
	var conf Configuration

	conf = CreateConfiguration(Read(configfile), filter)

	AddLink(conf, 0, answers)
	respond := make(chan Answer)
	for item := range answers.IterBuffered() {
		v := item.Val.(Answer)
		go AddUriToAria(respond, v.Request, v, url)
	}

	for item := range answers.IterBuffered() {
		v := item.Val.(Answer)
		utils.Info.Println("File: ", v.Name)
		utils.Info.Println("Respond: ", v.AriaID)
	}
}

func _() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	configfile, _, filter, err := utils.GetArguments()
	if err != nil {
		panic(err)
	}
	Run(configfile, filter)
}
