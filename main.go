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
	HTMLString := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
  <html itemscope itemtype="http://schema.org/QAPage">
  <head>
     <title>go - Golang parse HTML, extract all content with &lt;body&gt; &lt;/body&gt; tags - Stack Overflow</title>
    <link rel="shortcut icon" href="//cdn.sstatic.net/Sites/stackoverflow/img/favicon.ico?v=4f32ecc8f43d">
    <link rel="apple-touch-icon image_src" href="//cdn.sstatic.net/Sites/stackoverflow/img/apple-touch-icon.png?v=c78bd457575a">
    <link rel="search" type="application/opensearchdescription+xml" title="Stack Overflow" href="/opensearch.xml">
    <meta name="twitter:card" content="summary">
    <meta name="twitter:domain" content="stackoverflow.com"/>
    <meta property="og:type" content="website" />
    </head>
<body class="template-blog">
<nav class="navigation">
<div class="navigation__container container">
<a class="navigation__logo" href="/page1">
<a class="navigation__logo" href="http://www.google.com">
<h1>Foobar</h1>
</a>
<ul class="navigation__menu">
<li><a href="/tags/">Topics</a></li>
<li><a href="/about">About</a></li>
</ul>
<form><div><input name="user">Username</input>
<input name="password">Password</input></div></form>
</div>`
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	//htmlData, err := http.Get(*webUrl, http.CustomClient())

	//bodyBytes, _ := ioutil.ReadAll(htmlData)
	//bodyString := string(bodyBytes)
	//fmt.Println(bodyString)

	/*if err != nil {
		fmt.Printf("Error: %v", err)
		panic(err)
	}*/
	desiredTokens := []string{
		"Doctype", "title", "a", "h1", "h2", "h3", "h4", "form",
	}
	r := strings.NewReader(HTMLString)
	z := html.NewTokenizer(r)

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
				if (resp != nil && resp.Status != "200") || err != nil {
					result.InaccessibleLinks++
				}
			}
		} else if node.Type == "H1" {
			result.H1++
		} else if node.Type == "H2" {
			result.H2++
		} else if node.Type == "H3" {
			result.H3++
		} else if node.Type == "H4" {
			result.H4++
		} else if node.Type == "H5" {
			result.H5++
		} else if node.Type == "form" {
			result.IsLoginForm = parser.CheckForLoginForm(&node)

		}
	}

	return result, nil
}
