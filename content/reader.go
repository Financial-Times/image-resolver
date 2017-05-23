package content

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/pkg/errors"
)

type Reader interface {
	Get(c UUIDBatch) (map[string]Content, error)
}

type ContentReader struct {
	client         *http.Client
	contentAppName string
	contentAppURL  string
	path           string
}

func NewContentReader(appName string, URL string, path string, client *http.Client) *ContentReader {
	return &ContentReader{
		client:         client,
		contentAppName: appName,
		contentAppURL:  URL,
		path:           path,
	}
}

func (cr *ContentReader) Get(c UUIDBatch) (map[string]Content, error) {
	var result = make(map[string]Content)
	req, err := http.NewRequest(http.MethodGet, cr.contentAppURL+cr.path, nil)
	if err != nil {
		return result, errors.Wrapf(err, "Error connecting to %v", cr.contentAppName)
	}
	q := req.URL.Query()
	for _, uuid := range c.toArray() {
		if "" == uuid {
			continue
		}
		q.Add("uuid", uuid)

	}
	req.URL.RawQuery = q.Encode()
	req.Host = cr.contentAppName

	res, err := cr.client.Do(req)
	if err != nil {
		return result, errors.Wrapf(err, "Request to %v failed.", cr.contentAppName)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return result, errors.Errorf("Request to %v failed with status code %d", cr.contentAppName, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, errors.Wrapf(err, "Error reading response received from %v", cr.contentAppName)
	}

	cb := []Content{}
	err = json.Unmarshal(body, &cb)
	if err != nil {
		return result, errors.Wrapf(err, "Error unmarshalling response from %v", cr.contentAppName)
	}

	for _, i := range cb {
		id, ok := i[id].(string)
		if !ok {
			continue
		}
		uuid := extractUUIDFromURL(id)
		result[uuid] = i
	}

	return result, nil
}
