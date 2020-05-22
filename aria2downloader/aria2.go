package aria2downloader

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/azak-azkaran/putio-go-aria2/utils"
)

type PurgeDownloadResult struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      string   `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type Answer struct {
	ID      string            `json:"id"`
	Link    string            `json:"link"`
	Name    string            `json:"name"`
	AriaID  string            `json:"ariaId"`
	Request AddUriRequest     `json:"request"`
	Status  TellStatusRequest `json:"status"`
}

type AddUriRequest struct {
	Jsonrpc string     `json:"jsonrpc"`
	ID      string     `json:"id"`
	Method  string     `json:"method"`
	Params  [][]string `json:"params"`
}

type TellStatusRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      string   `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}
type Response struct {
	ID      string
	Jsonrpc string
	Result  string
}

func AddURI(link string) AddUriRequest {
	request := AddUriRequest{}
	request.Jsonrpc = "2.0"
	request.ID = "qwer"
	request.Method = "aria2.addUri"
	var nested []string
	nested = append(nested, link)
	request.Params = append(request.Params, nested)
	return request
}

func TellStatus(link string) TellStatusRequest {
	request := TellStatusRequest{}
	request.Jsonrpc = "2.0"
	request.ID = "qwer"
	request.Method = "aria2.tellStatus"
	request.Params = append(request.Params, link)
	return request
}

func PurgeDownload() PurgeDownloadResult {
	request := PurgeDownloadResult{}
	request.Jsonrpc = "2.0"
	request.ID = "qwer"
	request.Method = "aria2.purgeDownloadResult"
	return request
}

func Send(request []byte, url string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	if err != nil {
		return "Error while converting json to http request:", err
	}
	req.Header.Set("Content-Type", "application/json")

	aria := http.Client{}
	resp, err := aria.Do(req)
	if err != nil {
		return "Error while sending to Aria2:", err
	}
	defer resp.Body.Close()
	utils.Info.Println("Status: ", resp.Status)

	decoder := json.NewDecoder(resp.Body)

	var result string
	for decoder.More() {
		var m Response
		err = decoder.Decode(&m)
		if err != nil {
			return "Error while decoding json:", err
		}
		utils.Info.Println("result: ", m.Result)
		result = m.Result
	}
	return result, nil
}
