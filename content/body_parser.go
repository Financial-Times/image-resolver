package content

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func getEmbedded(body string, embedsType string, tid string, uuid string) ([]imageSetUUID, error) {
	embedsImg := []imageSetUUID{}
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

func parse(n *html.Node, re *regexp.Regexp, embedsImg *[]imageSetUUID, tid string, uuid string) {
	if n.Data == "ft-content" {
		isEmbedded := false
		isImageSet := false
		var id string
		for _, a := range n.Attr {
			if a.Key == "data-embedded" && a.Val == "true" {
				isEmbedded = true
			} else if a.Key == "type" {
				values := re.FindStringSubmatch(a.Val)
				if len(values) == 1 {
					isImageSet = true
				}
			} else if a.Key == "url" {
				id = a.Val
			}
		}
		var err error
		if isEmbedded && isImageSet {
			var emb imageSetUUID
			emb.uuid = extractUUIDFromURL(id)
			emb.imageModelUUID, err = getImageModelUUID(emb.uuid)
			if err != nil {
				logger.Infof(tid, uuid, "Cannot get image model UUID from image set UUID %s", emb.uuid)
			}
			*embedsImg = append(*embedsImg, emb)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, re, embedsImg, tid, uuid)
	}
}
