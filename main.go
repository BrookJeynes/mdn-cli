package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"golang.org/x/net/html"
)

var version string = "mdn v0.1.0\n"
var help_menu string = `Usage: mdn [search-terms]

search mdn docs from the terminal

Flags:
  -h, --help                     Show help information and exit.
  -v, --version                  Print version information and exit.

Example
  mdn ".join()"
  mdn accept header
`
var main_content_class string = "main-page-content"
var url_base string = "https://api.duckduckgo.com?format=json&q=! site:developer.mozilla.org "

func main() {
	if len(os.Args) == 1 {
		fmt.Fprint(os.Stdout, help_menu)
		os.Exit(2)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprint(os.Stdout, help_menu)
		os.Exit(0)
	}

	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Fprint(os.Stdout, version)
		os.Exit(0)
	}

	url := url_base + strings.Join(os.Args[1:], " ")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to retrieve docs - %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse html body - %v\n", err)
		os.Exit(1)
	}

	var markdown []byte = nil
	for child_node := range node.Descendants() {
		attributes := child_node.Attr
		if len(attributes) == 0 {
			continue
		}

		for _, attr := range attributes {
			if attr.Key != "class" {
				continue
			}

			if attr.Val != main_content_class {
				continue
			}

			markdown, err = htmltomarkdown.ConvertNode(child_node)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to convert html to markdown - %v\n", err)
				os.Exit(1)
			}

			break
		}
	}

	if markdown == nil {
		fmt.Fprint(os.Stderr, "error: unable to find content for provided search terms.\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(markdown))
	os.Exit(0)
}
