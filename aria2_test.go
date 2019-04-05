package main

import (
	"encoding/json"
	"fmt"
	"github.com/qri-io/jsonschema"
	"io/ioutil"
	"strings"
	"testing"
)

func TestAddUri(t *testing.T) {

	link := "testLink"
	request := AddURI(link)
	b, err := json.Marshal(request)

	if err != nil {
		t.Error(err)
	}
	requestString := string(b)

	if !strings.Contains(requestString, "params\":[[\"testLink\"]]") {
		fmt.Println("generated JSON\t:", requestString)
		t.Error("Link part is not in the JSON")
	}

	if !strings.Contains(requestString, "aria2.addUri") {
		fmt.Println("generated JSON\t:", requestString)
		t.Error("method name not in JSON")
	}

	schemaData, err := ioutil.ReadFile("RequestSchema.json")
	if err != nil {
		t.Error(err)
	}
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		t.Error("unmarshal schema: " + err.Error())
	}
	if errors, _ := rs.ValidateBytes(b); len(errors) > 0 {
		t.Error(errors)
	}
}

func TestTellStatus(t *testing.T) {

	link := "testLink"
	request := TellStatus(link)
	b, err := json.Marshal(request)

	if err != nil {
		t.Error(err)
	}
	requestString := string(b)
	fmt.Println(requestString)

	if !strings.Contains(requestString, "params\":[\"testLink\"]") {
		fmt.Println("generated JSON\t:", requestString)
		t.Error("Link part is not in the JSON")
	}

	if !strings.Contains(requestString, "aria2.tellStatus") {
		fmt.Println("generated JSON\t:", requestString)
		t.Error("method name not in JSON")
	}

	schemaData, err := ioutil.ReadFile("RequestSchema.json")
	if err != nil {
		t.Error(err)
	}
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		t.Error("unmarshal schema: " + err.Error())
	}
	if errors, _ := rs.ValidateBytes(b); len(errors) > 0 {
		t.Error(errors)
	}
}
