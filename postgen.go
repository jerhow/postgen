package main

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	// "path/filepath" // NOTE: Holding this for something I'll be adding soon.
	"encoding/json"
	"net/url"
	"strings"
)

// For holding our post's configuration settings
type Config struct {
	Title     string `json:"title`
	Date      string `json:"date`
	Permalink string `json:"string`
}

func main() {

	// NOTE: Holding this for something I'll be adding soon.
	//
	// root := "./content"
	// err := filepath.Walk(root, visit)
	// fmt.Printf("filepath.Walk() returned %v\n", err)
	// os.Exit(1)

	// Note: index 0 contains the program path, so I'm excluding it from what gets passed in
	inputFile, outputFile := dealWithArgs(os.Args[1:])

	title, date, permalink := getPostConfigsJson(&inputFile)

	fmt.Println("postTitle: " + title)
	fmt.Println("postDate: " + date)
	fmt.Println("permalink: " + permalink)
	// fmt.Println("Leaving off here")
	// os.Exit(1)

	content := getContent(&inputFile)
	topHTML, bottomHTML := getSharedMarkup()
	combinedOutput := buildCombinedOutput(topHTML, content, bottomHTML)
	finalOutput := interpolateConfigVals(combinedOutput, &title, &date, &permalink)
	writeOutputFile(finalOutput, &outputFile)
	fmt.Println("Program finished, check result.")
}

func dealWithArgs(args []string) (string, string) {
	if len(args) != 2 {
		fmt.Println("ERROR: Wrong number of arguments provided. We're expecting:")
		fmt.Println("1. Input file")
		fmt.Println("2. Output file")
		os.Exit(1)
	}
	inputFile := args[0]
	outputFile := args[1]

	return inputFile, outputFile
}

func writeOutputFile(finalOutput []byte, outputFile *string) bool {
	// Write our output to an HTML file
	err := ioutil.WriteFile("./content/"+*outputFile, finalOutput, 0644)
	check(err)
	return true
}

// Interpolate config values into the output
func interpolateConfigVals(combinedOutput []byte, title *string, date *string, permalink *string) []byte {

	encodedSquigglyOpen := url.QueryEscape("{")
	encodedSquigglyClose := url.QueryEscape("}")
	permalinkPlaceHolder := encodedSquigglyOpen + encodedSquigglyOpen + "permalink" + encodedSquigglyClose + encodedSquigglyClose

	str := string(combinedOutput[:])
	str = strings.Replace(str, "{{title}}", *title, -1)
	str = strings.Replace(str, "{{date}}", *date, -1)
	str = strings.Replace(str, permalinkPlaceHolder, *permalink, -1)

	return []byte(str)
}

func buildCombinedOutput(topHTML []byte, content []byte, bottomHTML []byte) []byte {
	// Incrementally building a singular byte array via append()
	combinedOutput := append(topHTML[:], content[:]...)
	combinedOutput = append(combinedOutput, bottomHTML[:]...)
	return combinedOutput
}

func getSharedMarkup() ([]byte, []byte) {

	// Get the shared markup pieces
	topHTML, err := ioutil.ReadFile("./shared_markup/top.html")
	check(err)
	bottomHTML, err := ioutil.ReadFile("./shared_markup/bottom.html")
	check(err)

	return topHTML, bottomHTML
}

func getContent(inputFile *string) []byte {

	// Get the Markdown input
	rawInput, err := ioutil.ReadFile("./content/" + *inputFile)
	check(err)

	// Convert and sanitize our content
	unsafe := blackfriday.MarkdownCommon(rawInput)
	policy := bluemonday.UGCPolicy()
	policy.RequireNoFollowOnLinks(false)
	// This allows us to arbitrarily add "rel" attributes to our links
	policy.AllowAttrs("rel").OnElements("a", "area")

	content := policy.SanitizeBytes(unsafe)
	return content
}

func getPostConfigsJson(inputFile *string) (string, string, string) {
	// Get the configs for this page:
	// First read the post's corresponding .json file
	configPath := "./content/" + strings.Replace(*inputFile, ".md", "", 1) + ".json"
	configJson, err := ioutil.ReadFile(configPath)
	check(err)

	// Then parse out our relevant config values for later use
	var config Config
	err = json.Unmarshal(configJson, &config)
	check(err)

	return config.Title, config.Date, config.Permalink
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
