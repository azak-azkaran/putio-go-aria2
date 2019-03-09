package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"encoding/json"
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

type Request struct {
	Jsonrpc string     `json:"jsonrpc"`
	Id      string     `json:"id"`
	Method  string     `json:"method"`
	Params  [][]string `json:"params"`
}

func NewRequest(link string) Request {
	request := Request{}
	request.Jsonrpc = "2.0"
	request.Id = "qwer"
	request.Method = "aria2.addUri"
	var nested []string
	nested = append(nested, link)
	request.Params = append(request.Params, nested)
	return request
}

func CreateRequests(links []string) []Request {
	var requests []Request
	for _, v := range links {
		requests = append(requests, NewRequest(v))
	}
	return requests
}

func AddLink(conf Configuration, dir int64, links []string) []string {
	fmt.Println("Checking folder: ", strconv.FormatInt(dir, 10))
	link := "https://api.put.io/v2/files/"
	var builder strings.Builder
	list, _, err := conf.client.Files.List(context.Background(), dir)
	if err != nil {
		log.Fatal("error:", err)
	}

	for _, value := range list {
		if len(conf.filter) == 0 && strings.Contains(value.Name, conf.filter) {
			if value.IsDir() {
				fmt.Println(value.Name, "is Folder adding to check for contents")
				links = AddLink(conf, value.ID, links)
			} else {
				builder.WriteString(link)
				builder.WriteString(strconv.FormatInt(value.ID, 10))
				//builder.WriteString("609933704")
				builder.WriteString("/download?oauth_token=")
				builder.WriteString(conf.oauthToken)
				currentlink := builder.String()
				builder.Reset()
				fmt.Println(value.ID, ": ", value.Name, "\nlink: ", currentlink)
				links = append(links, currentlink)

			}
		}
	}
	return links
}

func Read(filename string, filter string) Configuration {
	file, err := os.Open(filename)
	var conf Configuration
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Reading configuration file")
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			conf.oauthToken = text
		}
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.oauthToken})
		conf.oauthClient = oauth2.NewClient(oauth2.NoContext, tokenSource)
		conf.client = putio.NewClient(conf.oauthClient)
	}
	if len(filter) > 0 {
		conf.filter = filter
	}
	return conf
}

func Send(jsonRequest Request) bool {
	//url := "http://172.17.0.2:6800/jsonrpc"
	url := "http://localhost:6800/jsonrpc"
	b, err := json.Marshal(jsonRequest)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	aria := http.Client{}
	resp, err := aria.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	test, err := io.Copy(os.Stdout, resp.Body)
	fmt.Println(test)
	return true
}

func main() {
	var links []string
	conf := Read("secret.conf", "")
	links = AddLink(conf, 0, links)
	request := CreateRequests(links)
	//fmt.Println("preparing to send:\n", request)
	for _, v := range request {
		Send(v)
	}
}
