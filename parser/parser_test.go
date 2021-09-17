package parser

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const testHTML = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
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

func TestFindMatchedString(t *testing.T) {
	r := strings.NewReader(testHTML)
	z := html.NewTokenizer(r)

	for z.Next() != html.ErrorToken {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			break
		case tt == html.DoctypeToken:
			token := z.Token()
			result := FindMatchedString("X?HTML \\d.\\d*", token.String())

			if result != "HTML 4.01" {
				t.Error("Expected HTML 4.01 but found", result)
			}
			break
		}
	}

}

func TestCheckForLoginForm(t *testing.T) {
	r := strings.NewReader(testHTML)
	z := html.NewTokenizer(r)
	result := false
	for z.Next() != html.ErrorToken {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			break
		case tt == html.StartTagToken:
			token := z.Token()
			if token.Data == "form" {
				result = CheckForLoginForm(&Node{Type: "form", Token: token, Doc: *z})

				if !result {
					t.Error("Expected true but found", result)
				}
				break
			}
			continue
		}
	}

}
