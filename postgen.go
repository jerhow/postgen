package main

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
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
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	err = ioutil.WriteFile(outputFile, safe, 0777)
	check(err)

	fmt.Println("Program finished, check result.")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
