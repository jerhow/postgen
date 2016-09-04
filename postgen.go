package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// For holding our post's configuration settings
type Config struct {
	Title     string `json:"title`
	Date      string `json:"date`
	Permalink string `json:"string`
}

var archiveList []Config // a slice of Config structs, for the 'Archive' page, etc

func main() {

	// Note: index 0 contains the program path,
	// so I'm excluding it from what gets passed in
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

	generateArchiveList() // exists in memory at this point
	writeArchiveList()    // write to file

	fmt.Println("Program finished, check result.")
}

func writeArchiveList() {

	var buffer bytes.Buffer
	// tmp := ""

	for _, postConfigStruct := range archiveList {
		fmt.Println("\nConfigs for file:")
		fmt.Println(postConfigStruct.Title)
		fmt.Println(postConfigStruct.Date)
		fmt.Println(postConfigStruct.Permalink)

		buffer.WriteString(
			"<div class=\"archive_row\">\n" +
				"<p>" + postConfigStruct.Date + "</p>\n" +
				"<p><a href=\"" + postConfigStruct.Permalink + "\">" +
				postConfigStruct.Title + "</a></p>\n" +
				"</div>\n\n")
	}

	// fmt.Println(buffer.String())
	topHTML, bottomHTML := getSharedMarkup()
	title := "Archive"
	date := ""
	permalink := "archive"
	content := buffer.Bytes() // we need a byte array to pass along
	combinedOutput := buildCombinedOutput(topHTML, content, bottomHTML)
	finalOutput := interpolateConfigVals(combinedOutput, &title, &date, &permalink)
	outputFile := "archive.html"
	writeOutputFile(finalOutput, &outputFile)
}

func generateArchiveList() {
	root := "./content"
	err := filepath.Walk(root, pushArchiveConfigs)
	check(err)
}

// This is a callback for filepath.Walk(), called from generateArchiveList()
func pushArchiveConfigs(path string, f os.FileInfo, err error) error {

	if filepath.Ext(path) == ".json" {
		// fmt.Printf("Visited: %s\n", path)
		configJson, err := ioutil.ReadFile(path)
		check(err)

		var config Config
		err = json.Unmarshal(configJson, &config)
		check(err)
		archiveList = append(archiveList, config)
	}

	return nil
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
func interpolateConfigVals(combinedOutput []byte, title *string, date *string,
	permalink *string) []byte {

	encodedSquigglyOpen := url.QueryEscape("{")
	encodedSquigglyClose := url.QueryEscape("}")
	permalinkPlaceHolder := encodedSquigglyOpen + encodedSquigglyOpen +
		"permalink" + encodedSquigglyClose + encodedSquigglyClose

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

// func visit(path string, f os.FileInfo, err error) error {
// 	fmt.Printf("Visited: %s\n", path)
// 	return nil
// }

func check(e error) {
	if e != nil {
		panic(e)
	}
}
