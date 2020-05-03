package aria2downloader

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	cmap "github.com/orcaman/concurrent-map"
)

func Run(oauthToken string, filter string, url string) {
	answers := cmap.New()
	var conf Configuration

	conf = CreateConfiguration(oauthToken, filter)

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
