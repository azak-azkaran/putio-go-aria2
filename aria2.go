package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Answer Comment
type Answer struct {
	ID      string  `json:"id"`
	Link    string  `json:"link"`
	Name    string  `json:"name"`
	AriaID  string  `json:"ariaId"`
	Request Request `json:"request"`
}

type Request struct {
	Jsonrpc string     `json:"jsonrpc"`
	ID      string     `json:"id"`
	Method  string     `json:"method"`
	Params  [][]string `json:"params"`
}

type Response struct {
	ID      string
	Jsonrpc string
	Result  string
}

func AddURI(link string) Request {
	request := Request{}
	request.Jsonrpc = "2.0"
	request.ID = "qwer"
	request.Method = "aria2.addUri"
	var nested []string
	nested = append(nested, link)
	request.Params = append(request.Params, nested)
	return request
}

func TellStatus(link string) Request {
	request := Request{}
	request.Jsonrpc = "2.0"
	request.ID = "qwer"
	request.Method = "aria2.tellStatus"
	var nested []string
	nested = append(nested, link)
	request.Params = append(request.Params, nested)
	return request
}

func Send(answer Answer, url string) bool {
	b, err := json.Marshal(answer.Request)
	if err != nil {
		Error.Fatalln(err)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		Error.Fatalln("Error while converting json to http request:", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	aria := http.Client{}
	resp, err := aria.Do(req)
	Info.Println("Status: ", resp.Status)
	if err != nil {
		Error.Fatalln("Error while sending to Aria2:", err)
		return false
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var m Response
		err = decoder.Decode(&m)
		if err != nil {
			Error.Fatalln("Error while decoding json:", err)
			Error.Fatalln(err)
		}
		Info.Println("result: ", m.Result)
		answer.AriaID = m.Result
	}

	return true
}
