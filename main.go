package main

import (
	"bufio"
	"encoding/json"
	"github.com/orcaman/concurrent-map"
	"io/ioutil"
	"os"
	"strings"
)

func Read(filename string) string {
	oauthToken := ""
	file, err := os.Open(filename)
	if err != nil {
		Error.Fatalln("could not read file", err)
		panic(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			oauthToken = text
		}
	}
	return oauthToken
}

func Write(foldername string, answer Answer) (string, error) {
	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		err = os.Mkdir(foldername, os.ModePerm)
		if err != nil {
			Error.Fatalln(err)
			return "Folder not created", err
		}
	}

	b, err := json.Marshal(answer)
	if err != nil {
		Error.Fatalln(err)
		return "JSON could not be created", err
	}
	var filename strings.Builder
	filename.WriteString(foldername)
	filename.WriteString("/")
	filename.WriteString(answer.ID)
	filename.WriteString(".json")

	err = ioutil.WriteFile(filename.String(), b, 0644)
	return filename.String(), err
}

func SendToAria(respond chan<- Answer, request Request, answer Answer, url string) {
	result, err := Send(request, url)
	if err != nil {
		Error.Fatalln(result, err)
		return
	} else {
		answer.AriaID = result
		filename, err := Write("jsons", answer)
		if err != nil {
			Warning.Println("File for: ", answer.Name, "\tFilename: ", filename)
		} else {
			Info.Println(answer.Name, " send to aria and writen to informatione writen to file: ", filename)
		}
	}
	respond <- answer
}

func main() {
	Init(os.Stdout, os.Stdout, os.Stderr)

	url := "http://localhost:6800/jsonrpc"
	answers := cmap.New()
	var conf Configuration
	argsWithProg := os.Args

	i := len(argsWithProg)
	switch i {
	case 2:
		Info.Println("reading config file: ", argsWithProg[1])
		conf = CreateConfiguration(Read(argsWithProg[1]), "")
		Info.Println("Running without filter")
	case 3:
		Info.Println("reading config file: ", argsWithProg[1])
		conf = CreateConfiguration(Read(argsWithProg[1]), argsWithProg[2])
		Info.Println("Running with filter", conf.filter)
	default:
		Error.Fatalln("script was used wrong:\n putio-go-aria2 oauth.secret\nputio-go-aria2 oauth.secret filter")
		panic("Dying horribly")
	}

	AddLink(conf, 0, answers)
	respond := make(chan Answer)
	for item := range answers.IterBuffered() {
		v := item.Val.(Answer)
		go SendToAria(respond, v.Request, v, url)
	}

	for item := range answers.IterBuffered() {
		v := item.Val.(Answer)
		Info.Println("File: ", v.Name)
		Info.Println("Respond: ", v.AriaID)
	}
}
