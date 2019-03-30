package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/putdotio/go-putio/putio"
	"golang.org/x/oauth2"
	"net/http"
)

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
		Info.Println("Using oauth Token: ", conf.oauthToken)
	} else {
		Error.Fatalln("No Token found")
		panic("No Token found")
	}

	if len(filter) > 0 {
		conf.filter = filter
	}
	return conf
}

func CreateLink(conf Configuration, value putio.File, links []string, answers []Answer) ([]string, []Answer) {
	if value.IsDir() {
		Info.Println(value.Name, "is Folder adding to check for contents")
		links, answers = AddLink(conf, value.ID, links, answers)
	} else {
		var currentAnswer Answer
		var builder strings.Builder
		builder.WriteString("https://api.put.io/v2/files/")
		builder.WriteString(strconv.FormatInt(value.ID, 10))
		builder.WriteString("/download?oauth_token=")
		builder.WriteString(conf.oauthToken)
		currentlink := builder.String()
		currentAnswer.ID = strconv.FormatInt(value.ID, 10)
		currentAnswer.Link = currentlink
		links = append(links, currentlink)
		answers = append(answers, currentAnswer)

		Info.Println(value.ID, ": ", value.Name, "\nlink: ", currentlink)

	}
	return links, answers
}

func AddLink(conf Configuration, dir int64, links []string, answers []Answer) ([]string, []Answer) {
	Info.Println("Checking folder: ", strconv.FormatInt(dir, 10))
	list, _, err := conf.client.Files.List(context.Background(), dir)
	if err != nil {
		Error.Fatalln("error:", err)
	}

	for _, value := range list {
		if len(conf.filter) == 0 {
			links, answers = CreateLink(conf, value, links, answers)
		} else if strings.Contains(value.Name, conf.filter) {
			links, answers = CreateLink(conf, value, links, answers)
		}
	}
	return links, answers
}
