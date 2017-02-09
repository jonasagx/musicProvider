package main

import (
	"testing"
	"io/ioutil"
)

var (
	NotFilteredFilenamesInput = []string{"Pixies  - Where Is My Mind-5iC0YXspJRM.webm", "a.mkv", "kj3lk2j4.mp4", "a.txt", "b.out"} 
	
	validFilenamesInput = []string{"Pixies  - Where Is My Mind-5iC0YXspJRM.webm", "a.mkv", "kj3lk2j4.mp4"}
	validFilenamesOutput = []string{"Pixies  - Where Is My Mind-5iC0YXspJRM.mp3", "a.mp3", "kj3lk2j4.mp3"}
)

func isValidList(validList, testingList []string, t *testing.T) {
	if len(validList) != len(testingList) {
		t.Error("Expected", len(validList), " elements but got", (testingList), "instead")
	}

	for index,value := range testingList {
		if validList[index] != value{
			t.Error("Expected", validList[index], "but got", value, "instead")
		}
	}
}

func TestRaplceFunc(t *testing.T){
	for index,filename := range validFilenamesInput {
		testableOutput := StringParserToMp3(filename)
		if testableOutput != validFilenamesOutput[index] {
			t.Error("Expected", validFilenamesOutput[index], "but got", testableOutput, "instead")
		}
	}
}

func TestFilterFilenames(t *testing.T) {
	testableOutput := FilterFilenames(NotFilteredFilenamesInput, `\.txt|\.out`)
	validOutput := []string{"a.txt", "b.out"}

	isValidList(validOutput, testableOutput, t)
}

func TestGetFilesList(t *testing.T) {
	var path = "."
	var validFilenames []string
	var testableOutput = GetFilesList(path)

	infoFiles, _ := ioutil.ReadDir(path)
    
    for _,filename := range infoFiles {
    	validFilenames = append(validFilenames, filename.Name())
    }

    isValidList(validFilenames, testableOutput, t)
}