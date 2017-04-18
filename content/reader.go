package content

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
)

type Reader interface {
	Get(uuid string) (Content, error)
	ConvertContentIntoOutput(result Content) (ContentOutput)
}

type ContentReader struct {
	client      *http.Client
	contentHost string
	routingAddr string
}

func NewContentReader(ch string, routingAddr string) *ContentReader {
	return &ContentReader{
		client:      http.DefaultClient,
		contentHost: ch,
		routingAddr: routingAddr,
	}
}

func (cr *ContentReader) Get(uuid string) (Content, error) {
	var result Content

	url := "http://" + cr.routingAddr + "/content/" + uuid
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Host = cr.contentHost

	res, err := cr.client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (cr *ContentReader) ConvertContentIntoOutput(result Content) (ContentOutput) {
	resultOutput := ContentOutput{}
	resultOutput.Uuid = result.Uuid
	resultOutput.Type = result.Type
	resultOutput.BodyXML = result.BodyXML
	resultOutput.OpeningXML = result.OpeningXML
	resultOutput.Title = result.Title
	resultOutput.Byline = result.Byline
	resultOutput.Description = result.Description
	resultOutput.PublishedDate = result.PublishedDate
	resultOutput.Identifiers = result.Identifiers
	resultOutput.RequestUrl = result.RequestUrl
	resultOutput.BinaryUrl = result.BinaryUrl
	resultOutput.Brands = result.Brands
	resultOutput.Annotations = result.Annotations
	resultOutput.MainImage = result.MainImage
	resultOutput.Comments = result.Comments
	resultOutput.Realtime = result.Realtime
	resultOutput.Copyright = result.Copyright
	resultOutput.PublishReference = result.PublishReference
	resultOutput.PixelWidth = result.PixelWidth
	resultOutput.PixelHeight = result.PixelHeight
	resultOutput.Stdout = result.Stdout
	resultOutput.LastModified = result.LastModified
	resultOutput.Standfirst = result.Standfirst
	resultOutput.AltTitles = result.AltTitles
	resultOutput.AltStandfirsts = result.AltStandfirsts
	resultOutput.AltImages = result.AltImages
	resultOutput.WebUrl = result.WebUrl
	resultOutput.CanBeSyndicated = result.CanBeSyndicated
	resultOutput.FirstPublishedDate = result.FirstPublishedDate
	resultOutput.AccessLevel = result.AccessLevel
	resultOutput.CanBeDistributed = result.CanBeDistributed
	return resultOutput
}
