package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
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

type Response struct {
	Id      string
	Jsonrpc string
	Result  string
}

type Answer struct {
	id      string  `json:"id"`
	link    string  `json:"link"`
	name    string  `json:"name"`
	ariaId  string  `json:"ariaId"`
	request Request `json:"request"`
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

func CreateLink(conf Configuration, value putio.File, links []string, answers []Answer) ([]string, []Answer) {
	if value.IsDir() {
		fmt.Println(value.Name, "is Folder adding to check for contents")
		links, answers = AddLink(conf, value.ID, links, answers)
	} else {
		var currentAnswer Answer
		var builder strings.Builder
		builder.WriteString("https://api.put.io/v2/files/")
		builder.WriteString(strconv.FormatInt(value.ID, 10))
		builder.WriteString("/download?oauth_token=")
		builder.WriteString(conf.oauthToken)
		currentlink := builder.String()
		fmt.Println(value.ID, ": ", value.Name, "\nlink: ", currentlink)
		currentAnswer.id = strconv.FormatInt(value.ID, 10)
		currentAnswer.link = currentlink
		links = append(links, currentlink)
		answers = append(answers, currentAnswer)

	}
	return links, answers
}

func AddLink(conf Configuration, dir int64, links []string, answers []Answer) ([]string, []Answer) {
	fmt.Println("Checking folder: ", strconv.FormatInt(dir, 10))
	list, _, err := conf.client.Files.List(context.Background(), dir)
	if err != nil {
		log.Fatal("error:", err)
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

func Read(filename string, filter string) Configuration {
	file, err := os.Open(filename)
	var conf Configuration
	if err != nil {
		panic(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			conf.oauthToken = text
		}
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.oauthToken})
		conf.oauthClient = oauth2.NewClient(oauth2.NoContext, tokenSource)
		conf.client = putio.NewClient(conf.oauthClient)
		fmt.Println("Using oauth Token: ", conf.oauthToken)
	}
	if len(filter) > 0 {
		conf.filter = filter
	}
	return conf
}

func Send(answer Answer) bool {
	//url := "http://172.17.0.2:6800/jsonrpc"
	url := "http://localhost:6800/jsonrpc"
	b, err := json.Marshal(answer.request)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	aria := http.Client{}
	resp, err := aria.Do(req)
	fmt.Println("Status: ", resp.Status)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var m Response
		err = decoder.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("result: ", m.Result)
		answer.ariaId = m.Result
	}

	return true
}

func main() {
	var links []string
	var answers []Answer
	var conf Configuration
	argsWithProg := os.Args

	i := len(argsWithProg)
	switch i {
	case 2:
		fmt.Println("reading config file: ", argsWithProg[1])
		conf = Read(argsWithProg[1], "")
		fmt.Println("Running without filter")
	case 3:
		fmt.Println("reading config file: ", argsWithProg[1])
		conf = Read(argsWithProg[1], argsWithProg[2])
		fmt.Println("Running with filter", conf.filter)
	default:
		panic("script was used wrong:\n putio-go-aria2 oauth.secret\nputio-go-aria2 oauth.secret filter")
	}

	links, answers = AddLink(conf, 0, links, answers)
	for _, v := range answers {
		v.request = NewRequest(v.link)
		Send(v)
	}
}
