package content

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Financial-Times/transactionid-utils-go"
	"github.com/Financial-Times/uuid-utils-go"
	"github.com/pkg/errors"
)

const (
	userAgent      = "User-Agent"
	userAgentValue = "UPP_content-unroller"
)

type Reader interface {
	Get([]string, string) (map[string]Content, error)
	GetInternal(uuids []string, tid string) (map[string]Content, error)
}

type ContentReader struct {
	client         *http.Client
	contentAppName string
	contentAppURL  string
}

func NewContentReader(appName string, URL string, client *http.Client) *ContentReader {
	return &ContentReader{
		client:         client,
		contentAppName: appName,
		contentAppURL:  URL,
	}
}

func (cr *ContentReader) Get(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)

	imgBatch, err := cr.doGet(uuids, tid)
	if err != nil {
		return cm, err
	}

	var imgModelUUIDs []string
	for _, i := range imgBatch {
		cr.addItemToMap(i, cm)
		if _, foundMembers := i[members]; foundMembers {
			imgModelUUIDs = append(imgModelUUIDs, i.getMembersUUID()...)
		}
	}

	if len(imgModelUUIDs) == 0 {
		return cm, nil
	}

	imgModelsList, err := cr.doGet(imgModelUUIDs, tid)
	if err != nil {
		return cm, err
	}

	for _, i := range imgModelsList {
		cr.addItemToMap(i, cm)
	}

	return cm, nil
}

func (cr *ContentReader) GetInternal(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)

	internalContent, err := cr.doGet(uuids, tid, "http://document-store-api:8080/internalcomponents")
	if err != nil {
		return cm, err
	}

	for _, i := range internalContent {
		cr.addItemToMap(i, cm)
	}

	return cm, nil
}

func (cr *ContentReader) addItemToMap(c Content, cm map[string]Content) {
	id, ok := c[id].(string)
	if !ok {
		return
	}
	uuid, err := extractUUIDFromString(id)
	if err != nil {
		return
	}
	cm[uuid] = c
}

func (cr *ContentReader) doGet(uuids []string, tid string, url ...string) ([]Content, error) {
	var cb []Content
	contentAppURL := cr.contentAppURL

	if url != nil {
		contentAppURL = url[1]
	}

	req, err := http.NewRequest(http.MethodGet, contentAppURL, nil)
	if err != nil {
		return cb, errors.Wrapf(err, "Error connecting to %v", cr.contentAppName)
	}
	req.Header.Add(transactionidutils.TransactionIDHeader, tid)
	req.Header.Set(userAgent, userAgentValue)
	q := req.URL.Query()
	for _, uuid := range uuids {
		if err = uuidutils.ValidateUUID(uuid); err == nil {
			q.Add("uuid", uuid)
		}
	}
	req.URL.RawQuery = q.Encode()

	res, err := cr.client.Do(req)
	if err != nil {
		return cb, errors.Wrapf(err, "Request to %v failed.", cr.contentAppName)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cb, errors.Errorf("Request to %v failed with status code %d", cr.contentAppName, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cb, errors.Wrapf(err, "Error reading response received from %v", cr.contentAppName)
	}

	err = json.Unmarshal(body, &cb)
	if err != nil {
		return cb, errors.Wrapf(err, "Error unmarshalling response from %v", cr.contentAppName)
	}
	return cb, nil
}
