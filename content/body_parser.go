package content

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func getEmbedded(body string, embedsType string, tid string, uuid string) ([]string, error) {
	embedsImg := []string{}
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return embedsImg, err
	}

	re, err := regexp.Compile(embedsType)
	if err != nil {
		return embedsImg, errors.Wrap(err, "Error while compiling whitelist")
	}

	parse(doc, re, &embedsImg, tid, uuid)
	return embedsImg, nil
}

func parse(n *html.Node, re *regexp.Regexp, embedsImg *[]string, tid string, uuid string) {
	if n.Data == "ft-content" {
		isEmbedded := false
		isImageSet := false
		var id string
		for _, a := range n.Attr {
			if a.Key == "data-embedded" && a.Val == "true" {
				isEmbedded = true
			} else if a.Key == "type" {
				values := re.FindStringSubmatch(a.Val)
				if len(values) > 0 {
					isImageSet = true
				}
			} else if a.Key == "url" {
				id = a.Val
			}
		}
		if isEmbedded && isImageSet {
			u, err := extractUUIDFromString(id)
			if err != nil {
				logger.Infof(tid, uuid, "Cannot extract UUID: %v", err.Error())
			} else {
				*embedsImg = append(*embedsImg, u)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, re, embedsImg, tid, uuid)
	}
}
