package content

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
)

type Reader interface {
	Get(uuid string) (Content, error)
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
