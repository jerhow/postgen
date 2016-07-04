package main

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"reflect"
)

func main() {

	args := os.Args[1:] // Reason: index 0 contains the program path
	if len(args) != 2 {
		fmt.Println("ERROR: Wrong number of arguments provided. We're expecting:")
		fmt.Println("1. Input file")
		fmt.Println("2. Output file")
		os.Exit(1)
	}
	inputFile := args[0]
	outputFile := args[1]

	input, err := ioutil.ReadFile(inputFile)
	check(err)

	unsafe := blackfriday.MarkdownCommon(input)
	content := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	fmt.Println(reflect.TypeOf(content))

	// Get the shared markup pieces
	topHTML, err := ioutil.ReadFile("./shared_markup/top.html")
	check(err)
	fmt.Println(reflect.TypeOf(topHTML))
	// fmt.Print(string(topHTML))
	// fmt.Println(topHTML)

	bottomHTML, err := ioutil.ReadFile("./shared_markup/bottom.html")
	check(err)
	fmt.Println(reflect.TypeOf(bottomHTML))
	fmt.Println(len(bottomHTML))

	fmt.Println("Trying our thing here...")
	finalOutput := append(topHTML[:], content[:]...)
	finalOutput = append(finalOutput, bottomHTML[:]...)
	// fmt.Println(string(testOutput))
	// fmt.Println(append(topHTML[:], bottomHTML[:]...))

	// capacity := len(topHTML) + len(content) + len(bottomHTML)
	// fmt.Println(capacity)

	// topHTML_str := string(topHTML)
	// content_str := string(content)
	// bottomHTML_str := string(bottomHTML)

	// finalOutput_str := topHTML_str + content_str + bottomHTML_str
	// finalOutput := []byte(finalOutput_str)

	err = ioutil.WriteFile(outputFile, finalOutput, 0644)
	check(err)

	fmt.Println("Program finished, check result.")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
