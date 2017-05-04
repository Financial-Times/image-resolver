package content

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Reader interface {
	Get(uuid string) (Content, error, int)
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

func (cr *ContentReader) Get(uuid string) (Content, error, int) {
	var result Content

	url := "http://" + cr.routingAddr + "/content/" + uuid
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Host = cr.contentHost

	res, err := cr.client.Do(req)

	if res != nil {
		if err != nil {
			return result, err, res.StatusCode
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return result, err, res.StatusCode
		}

		err = json.Unmarshal(body, &result)

		if err != nil {
			return result, err, res.StatusCode
		}
		return result, nil, res.StatusCode
	}
	return result, nil, http.StatusInternalServerError
}
