package content

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func getEmbedded(body string, embedsType string, tid string, uuid string) ([]string, error) {
	embedsResult := []string{}
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return embedsResult, err
	}

	re, err := regexp.Compile(embedsType)
	if err != nil {
		return embedsResult, errors.Wrap(err, "Error while compiling whitelist")
	}

	parse(doc, re, &embedsResult, tid, uuid)
	return embedsResult, nil
}

func parse(n *html.Node, re *regexp.Regexp, embedsResult *[]string, tid string, uuid string) {
	if n.Data == "ft-content" {
		isEmbedded := false
		isTypeMaching := false
		var id string
		for _, a := range n.Attr {
			if a.Key == "data-embedded" && a.Val == "true" {
				isEmbedded = true
			} else if a.Key == "type" {
				values := re.FindStringSubmatch(a.Val)
				if len(values) > 0 {
					isTypeMaching = true
				}
			} else if a.Key == "url" {
				id = a.Val
			}
		}

		if isEmbedded && isTypeMaching {
			u, err := extractUUIDFromString(id)
			if err != nil {
				logger.Infof(tid, uuid, "Cannot extract UUID: %v", err.Error())
			} else {
				*embedsResult = append(*embedsResult, u)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, re, embedsResult, tid, uuid)
	}
}
