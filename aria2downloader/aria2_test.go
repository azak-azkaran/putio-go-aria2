package aria2downloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/qri-io/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddUri(t *testing.T) {
	fmt.Println("Running TestAddUri")
	link := "testLink"
	request := AddURI(link)
	b, err := json.Marshal(request)
	require.NoError(t, err)

	requestString := string(b)
	fmt.Println(requestString)

	assert.Contains(t, requestString, "params\":[[\"testLink\"]]")
	assert.Contains(t, requestString, "aria2.addUri")

	schemaData, err := ioutil.ReadFile("../schema/AddAriaRequest.json")
	assert.NoError(t, err)
	rs := &jsonschema.RootSchema{}
	err = json.Unmarshal(schemaData, rs)
	assert.NoError(t, err)
	errors, _ := rs.ValidateBytes(b)
	assert.Empty(t, errors)

}

func TestTellStatus(t *testing.T) {
	fmt.Println("Running TestTellStatus")
	link := "testLink"
	request := TellStatus(link)
	b, err := json.Marshal(request)
	require.NoError(t, err)

	requestString := string(b)
	fmt.Println(requestString)

	assert.Contains(t, requestString, "params\":[\"testLink\"]")
	assert.Contains(t, requestString, "aria2.tellStatus")
	schemaData, err := ioutil.ReadFile("../schema/TellStatusRequest.json")
	assert.NoError(t, err)

	rs := &jsonschema.RootSchema{}
	err = json.Unmarshal(schemaData, rs)
	assert.NoError(t, err)
	errors, _ := rs.ValidateBytes(b)
	assert.Empty(t, errors)
}

func TestPurgeDownloadResult(t *testing.T) {
	fmt.Println("Running TestPurgeDownloadResult")
	request := PurgeDownload()
	b, err := json.Marshal(request)
	require.NoError(t, err)

	requestString := string(b)
	fmt.Println(requestString)

	assert.Contains(t, requestString, "params\":[]")
	assert.Contains(t, requestString, "aria2.purgeDownloadResult")
	schemaData, err := ioutil.ReadFile("../schema/PurgeDownload.json")
	assert.NoError(t, err)

	rs := &jsonschema.RootSchema{}
	err = json.Unmarshal(schemaData, rs)
	assert.NoError(t, err)
	errors, _ := rs.ValidateBytes(b)
	assert.Empty(t, errors)
}
