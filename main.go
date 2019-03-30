package main

import (
	"bufio"
	"encoding/json"
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
	filename.WriteString(answer.AriaID)
	filename.WriteString(".json")

	err = ioutil.WriteFile(filename.String(), b, 0644)
	return filename.String(), err
}

func main() {
	Init(os.Stdout, os.Stdout, os.Stderr)

	url := "http://localhost:6800/jsonrpc"
	var links []string
	var answers []Answer
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

	_, answers = AddLink(conf, 0, links, answers)
	for _, v := range answers {
		v.Request = AddURI(v.Link)
		if Send(v, url) {
			filename, err := Write("jsons", v)
			if err != nil {
				Warning.Println("File for: ", v.Name, "\tFilename: ", filename)
			} else {
				Info.Println(v.Name, " send to aria and writen to informatione writen to file: ", filename)
			}
		}
	}
}
