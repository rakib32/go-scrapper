package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-scrapper/http"
	"go-scrapper/parser"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

var (
	webUrl = flag.String("url", "", "the web url")
)

func init() {
	flag.Parse()
}

func main() {

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	htmlData, err := http.Get(*webUrl, http.CustomClient())

	if err != nil {
		fmt.Printf("Error: %v", err)
		panic(err)
	}
	desiredTokens := []string{
		"Doctype", "title", "a", "h1", "h2", "h3", "h4", "h5", "form",
	}
	//r := strings.NewReader(htmlData.Body)
	z := html.NewTokenizer(htmlData.Body)

	ch := make(chan parser.Node, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		errCh <- parse(z, ch, desiredTokens)
	}()

	// NOTE: this is a blocking call, when this returns, the parser will be done.
	result, err := aggregate(ch)

	if err := <-errCh; err != nil {
		// fail
		fmt.Printf("error: %v", err)
	}

	if err != nil {
		// fail
		fmt.Printf("error: %v", err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)

}

type Result struct {
	Version           string
	Title             string
	H1                int
	H2                int
	H3                int
	H4                int
	H5                int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	IsLoginForm       bool
}

func parse(root *html.Tokenizer, ch chan<- parser.Node, desiredTokens []string) error {
	// this is a producer
	defer close(ch)

	for {
		tt := root.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return nil
		case tt == html.StartTagToken || tt == html.DoctypeToken:
			token := root.Token()
			for _, name := range desiredTokens {
				if token.Data == name || token.Type.String() == name {
					ch <- parser.Node{Type: name, Token: token, Doc: *root}
				}
				continue
			}
		}
	}

}

func aggregate(ch <-chan parser.Node) (Result, error) {
	result := Result{}
	for node := range ch {
		if node.Type == "Doctype" {
			matchedData := parser.FindMatchedString("X?HTML \\d.\\d*", node.Token.String())
			if matchedData != "" {
				result.Version = matchedData
			}
		} else if node.Type == "title" {
			result.Title = parser.FindText(&node)
		} else if node.Type == "a" {
			href := parser.Attributes(&node.Token, "href")

			if href != "" {
				isInternal := false
				if !strings.Contains(href, "http") && (!strings.Contains(href, "https")) {
					result.InternalLinks++
					isInternal = true
				} else {
					result.ExternalLinks++
				}
				url := href

				if isInternal {
					url = *webUrl + href
				}

				resp, err := http.Get(url, http.CustomClient())
				if (resp != nil && resp.StatusCode != 200) || err != nil {
					result.InaccessibleLinks++
				}
			}
		} else if node.Type == "h1" {
			result.H1++
		} else if node.Type == "h2" {
			result.H2++
		} else if node.Type == "h3" {
			result.H3++
		} else if node.Type == "h4" {
			result.H4++
		} else if node.Type == "h5" {
			result.H5++
		} else if node.Type == "form" {
			result.IsLoginForm = parser.CheckForLoginForm(&node)

		}
	}

	return result, nil
}
