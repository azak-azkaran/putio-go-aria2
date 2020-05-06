package aria2downloader

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/azak-azkaran/putio-go-aria2/utils"
)

func Read(filename string) string {
	oauthToken := ""
	file, err := os.Open(filename)
	if err != nil {
		utils.Error.Fatalln("could not read file", err)
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

func Write(folder string, answer Answer) (string, error) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			utils.Error.Fatalln(err)
			return "Folder not created", err
		}
	}

	b, err := json.Marshal(answer)
	if err != nil {
		utils.Error.Fatalln(err)
		return "JSON could not be created", err
	}
	var filename strings.Builder
	filename.WriteString(folder)
	filename.WriteString("/")
	filename.WriteString(answer.ID)
	filename.WriteString(".json")

	err = ioutil.WriteFile(filename.String(), b, 0644)
	return filename.String(), err
}

func AddUriToAria(request AddUriRequest, answer Answer, url string) {
	b, err := json.Marshal(request)
	if err != nil {
		utils.Error.Fatalln("Error while Marshaling the Request: ", err)
		return
	}
	result, err := Send(b, url)
	if err != nil {
		utils.Error.Fatalln(result, err)
		return
	} else {
		answer.AriaID = result
		utils.Info.Println("Successfully sent to aria: ", answer.Name, " - ", answer.AriaID)
		//filename, err := Write("jsons", answer)
		//if err != nil {
		//	utils.Warning.Println("File for: ", answer.Name, "\tFilename: ", filename)
		//} else {
		//	utils.Info.Println(answer.Name, " send to aria and written to information written to file: ", filename)
		//}
	}
}
