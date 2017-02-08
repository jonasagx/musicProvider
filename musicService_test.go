package main

import (
	"testing"
)

var (
	NotFilteredFilenamesInput = []string{"Pixies  - Where Is My Mind-5iC0YXspJRM.webm", "a.mkv", "kj3lk2j4.mp4", "a.txt", "b.out"} 
	
	validFilenamesInput = []string{"Pixies  - Where Is My Mind-5iC0YXspJRM.webm", "a.mkv", "kj3lk2j4.mp4"}
	validFilenamesOutput = []string{"Pixies  - Where Is My Mind-5iC0YXspJRM.mp3", "a.mp3", "kj3lk2j4.mp3"}
)

func TestRaplceFunc(t *testing.T){
	for index,filename := range validFilenamesInput {
		testableOutput := StringParserToMp3(filename)
		if testableOutput != validFilenamesOutput[index] {
			t.Error("Expected", validFilenamesOutput[index], "but got", testableOutput, "instead")
		}
	}
}

func TestFilterFilenames(t *testing.T) {
	testableOutput := FilterFilenames(NotFilteredFilenamesInput, `\.txt`)
	validOutput := []string{"a.txt"}

	if len(validOutput) != len(testableOutput) {
		t.Error("Expected", len(validOutput), " elements but got", (testableOutput), "instead")
	}

	for index,value := range testableOutput {
		if validOutput[index] != value{
			t.Error("Expected", validOutput[index], "but got", value, "instead")
		}
	}
}