package main

import (
	"fmt"
	"github.com/bxcodec/faker"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	fmt.Println("Running TestWrite")
	foldername := "testFolder"
	var answer Answer
	err := faker.FakeData(&answer)
	if err != nil {
		t.Error("could not create Fake data")
	} else {
		fmt.Println("created fake data:\nID:", answer.ID, "\tAriaID:", answer.AriaID)
	}

	filename, err := Write(foldername, answer)
	if err != nil {
		t.Error("testfile could not be created\nFilename: ", filename)
	}

	if strings.Compare(filename, foldername+"/"+answer.ID+".json") != 0 {
		t.Error("filename is not correct")
	}

	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		t.Error("folder was not created")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	err = os.Remove(filename)
	if err != nil {
		t.Error("File could not be removed")
	}

	err = os.Remove(foldername)
	if err != nil {
		t.Error("Folder could not be removed")
	}
}

func TestRead(t *testing.T) {
	fmt.Println("Running TestRead")
	foldername := "testFolder"
	filename := "testsecret"
	testsecret := "test"

	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		err = os.Mkdir(foldername, os.ModePerm)
		if err != nil {
			Error.Fatalln(err)
			t.Error("Folder not created", err)
		}
	}
	data := []byte(testsecret)
	err := ioutil.WriteFile(foldername+"/"+filename, data, 0644)
	if err != nil {
		t.Error("File was not created")
	}
	read := Read(foldername + "/" + filename)
	if strings.Compare(read, testsecret) != 0 {
		t.Error("Content of file not correct: ", read)
	}
	err = os.Remove(foldername + "/" + filename)
	if err != nil {
		t.Error("File could not be removed")
	}

	err = os.Remove(foldername)
	if err != nil {
		t.Error("Folder could not be removed")
	}
}
