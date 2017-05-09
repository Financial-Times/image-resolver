package content

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
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
		return result, fmt.Errorf("Error connecting to content-public-read for uuid=%s, err=%v", uuid, err)
	}
	req.Host = cr.contentHost

	res, err := cr.client.Do(req)
	if err != nil {
		return result, fmt.Errorf("Error connecting to content-public-read for uuid=%s, err=%v", uuid, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return result, fmt.Errorf("Cannot read from content-public-read for uuid=%s, err=%v", uuid, err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("Error reading response from content-public-read for uuid=%s, err=%v", uuid, err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("Error unmarshalling response from content-public-read for uuid=%s, err=%v", uuid, err)
	}

	return result, nil
}
