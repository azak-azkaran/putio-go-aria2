package aria2downloader

import (
	"encoding/json"

	"github.com/azak-azkaran/putio-go-aria2/utils"
	cmap "github.com/orcaman/concurrent-map"
)

func Run(oauthToken string, filter string, url string) {
	answers := cmap.New()
	conf := CreateConfiguration(oauthToken, filter)

	AddLink(conf, 0, answers)
	for item := range answers.IterBuffered() {
		v := item.Val.(Answer)
		AddUriToAria(v.Request, v, url)
	}

	//utils.Info.Println("Sending PurgeDownloadResult Request")

	//b, err := json.Marshal(PurgeDownload())
	//if err != nil {
	//	utils.Error.Println("Error while Marshaling the PurgeRequest: ", err)
	//	return
	//}
	//result, err := Send(b, url)
	//if err != nil {
	//	utils.Error.Println(result, err)
	//}
}
