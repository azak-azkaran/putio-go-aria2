package organize

import (
	"context"
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"github.com/orcaman/concurrent-map"
	"github.com/putdotio/go-putio/putio"
	"golang.org/x/oauth2"
	"net/http"
)

type PutIoFiles struct {
	PutIoID int64
	Folder  string
	Name    string
	CRC32   string
	Size    int64
}
type Configuration struct {
	oauthToken  string
	oauthClient *http.Client
	client      *putio.Client
	filter      string
}

func CreateConfiguration(oauthToken string, filter string) Configuration {
	var conf Configuration
	if len(oauthToken) > 0 {
		conf.oauthToken = oauthToken
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.oauthToken})
		conf.oauthClient = oauth2.NewClient(context.TODO(), tokenSource)
		conf.client = putio.NewClient(conf.oauthClient)
		utils.Info.Println("Using oauth Token: ", conf.oauthToken)
	} else {
		utils.Error.Fatalln("No Token found")
		panic("No Token found")
	}

	if len(filter) > 0 {
		conf.filter = filter
	}
	return conf
}

func GetFolderInformation(conf Configuration, folderName string, dir int64, folders cmap.ConcurrentMap) {
	utils.Info.Println("Checking folder: ", folderName)
	list, _, err := conf.client.Files.List(context.Background(), dir)
	if err != nil {
		utils.Error.Fatalln("error:", err)
	}

	for _, value := range list {
		if value.IsDir() {
			GetFolderInformation(conf, folderName+"/"+value.Name, value.ID, folders)
		} else {
			var currentAnswer PutIoFiles
			currentAnswer.PutIoID = value.ID
			currentAnswer.Name = value.Name
			currentAnswer.Folder = folderName
			currentAnswer.CRC32 = value.CRC32
			currentAnswer.Size = value.Size
			folders.Set(value.Name, currentAnswer)
		}
	}
}

func RemoveOnlineFile(conf Configuration, file PutIoFiles) {
	utils.Info.Println("Deleting File: ", file.Name)
	err := conf.client.Files.Delete(context.Background(), file.PutIoID)
	if err != nil {
		utils.Error.Fatalln("error:", err)
	}
}
