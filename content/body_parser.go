package content

import (
	"errors"
	"fmt"
	"strings"
	"golang.org/x/net/html"
)

type Parser interface {
	GetEmbedded(content Content) ([]string, error)
}

type BodyParser struct{}

func (bp BodyParser) GetEmbedded(content Content) ([]string, error) {
	var ids []string

	if content.BodyXML == "" {
		return ids, errors.New(fmt.Sprintf("Cannot parse empty body of content [%s]", content.Uuid))
	}

	ids, err := parseXMLBody(content.BodyXML)
	if err != nil {
		return ids, err
	}
	return ids, nil
}

func parseXMLBody(body string) ([]string, error) {
	embedsImg := []string{}
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return embedsImg, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "ft-content" {
			isEmbedded := false
			var uuid string
			for _, a := range n.Attr {
				if a.Key == "data-embedded" && a.Val == "true" {
					isEmbedded = true
				}
				if a.Key == "url" {
					splitUrl := strings.Split(a.Val, "/")
					uuid = splitUrl[len(splitUrl)-1]
				}
			}
			if isEmbedded {
				embedsImg = append(embedsImg, uuid)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return embedsImg, nil
}
