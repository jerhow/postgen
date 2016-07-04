package main

import (
	"fmt" // A package in the Go standard library.
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
)

func main() {
	input, err := ioutil.ReadFile("test.md")
	check(err)

	unsafe := blackfriday.MarkdownCommon(input)
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	err = ioutil.WriteFile("test.html", safe, 0777)
	check(err)

	fmt.Println("Program finished, check result.")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
