package parser

import (
	"golang.org/x/net/html"
	"regexp"
)

type Node struct {
	Type  string
	Token html.Token
	Doc   html.Tokenizer
}

// FindMatchedString should return matched string when  found.
func FindMatchedString(expression, token string) string {
	result := ""
	rx, _ := regexp.Compile(expression)
	matchedData := rx.FindString(token)

	if matchedData != "" {
		result = matchedData
	}
	return result
}

// FindText should return text of that token.
func FindText(node *Node) string {
	result := ""
	nextToken := node.Doc.Next()

	if nextToken == html.TextToken {
		txtNode := node.Doc.Token()
		result = txtNode.Data
	}

	return result
}

// CheckForLoginForm should return true if found.
func CheckForLoginForm(node *Node) bool {
	isFound := false

	for node.Doc.Next() != html.ErrorToken {
		nextToken := node.Doc.Next()

		switch {
		case nextToken == html.ErrorToken:
			// End of the document, we're done
			break
		case nextToken == html.StartTagToken:
			token := node.Doc.Token()
			if token.Data == "input" {
				name := Attributes(&token, "name")
				id := Attributes(&token, "id")

				if name == "password" || id == "passwor" {
					isFound = true
					break
				}
			}
			continue
		}
	}

	return isFound
}

func Attributes(node *html.Token, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
