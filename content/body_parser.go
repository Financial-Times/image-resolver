package content

import (
	"golang.org/x/net/html"
	"strings"
	"regexp"
)

type Parser interface {
	GetEmbedded(content string) ([]string, error)
}

type BodyParser struct {
	embedsType string
}

func NewBodyParser(embedsType string) *BodyParser {
	return &BodyParser{
		embedsType: embedsType,
	}
}

func (bp *BodyParser) GetEmbedded(body string) ([]string, error) {
	embedsImg := []string{}
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return embedsImg, err
	}

	re, _ := regexp.Compile(bp.embedsType)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "ft-content" {
			isEmbedded := false
			isImageSet := false
			var uuid string
			for _, a := range n.Attr {
				if a.Key == "data-embedded" && a.Val == "true" {
					isEmbedded = true
				}
				if a.Key == "type" {
					values := re.FindStringSubmatch(a.Val)
					if len(values) == 1 {
						isImageSet = true
					}
				}
				if a.Key == "url" {
					uuid = a.Val
				}
			}
			if isEmbedded && isImageSet{
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

