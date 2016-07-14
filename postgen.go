package main

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	// "path/filepath"
	"encoding/json"
	"strings"
)

// For holding our post's configuration settings
type Config struct {
	Title string `json:"title`
	Date  string `json:"date`
}

func main() {

	// root := "./content"
	// err := filepath.Walk(root, visit)
	// fmt.Printf("filepath.Walk() returned %v\n", err)
	// os.Exit(1)

	args := os.Args[1:] // Reason: index 0 contains the program path
	if len(args) != 2 {
		fmt.Println("ERROR: Wrong number of arguments provided. We're expecting:")
		fmt.Println("1. Input file")
		fmt.Println("2. Output file")
		os.Exit(1)
	}
	inputFile := args[0]
	outputFile := args[1]

	// Get the configs for this page:
	// First read the post's corresponding .json file
	configFile := "./content/" + strings.Replace(inputFile, ".md", "", 1) + ".json"
	configJson, err := ioutil.ReadFile(configFile)
	check(err)

	// Then parse out our relevant config values for later use
	var config Config
	err = json.Unmarshal(configJson, &config)
	check(err)
	fmt.Println("postTitle: " + config.Title)
	fmt.Println("postDate: " + config.Date)
	fmt.Println("Leaving off here")
	os.Exit(1)

	// Get the Markdown input
	rawInput, err := ioutil.ReadFile(inputFile)
	check(err)

	// Convert and sanitize our content
	unsafe := blackfriday.MarkdownCommon(rawInput)
	content := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	// Get the shared markup pieces
	topHTML, err := ioutil.ReadFile("./shared_markup/top.html")
	check(err)
	bottomHTML, err := ioutil.ReadFile("./shared_markup/bottom.html")
	check(err)

	// Incrementally building a singular byte array via append()
	finalOutput := append(topHTML[:], content[:]...)
	finalOutput = append(finalOutput, bottomHTML[:]...)

	// Write our output to an HTML file
	err = ioutil.WriteFile(outputFile, finalOutput, 0644)
	check(err)

	fmt.Println("Program finished, check result.")
}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
