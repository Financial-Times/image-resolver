package content

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

type ReaderConfig struct {
	ContentSourceAppName         string
	ContentSourceAppURL          string
	InternalContentSourceAppName string
	InternalContentSourceAppURL  string
	NativeContentSourceAppName   string
	NativeContentSourceAppURL    string
	NativeContentSourceAppAuth   string
}

type ContentReader struct {
	client *http.Client
	config ReaderConfig
}

func NewContentReader(rConfig ReaderConfig, client *http.Client) *ContentReader {
	return &ContentReader{
		client: client,
		config: rConfig,
	}
}

// Get content from content-public-read
func (cr *ContentReader) Get(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)

	imgBatch, err := cr.doGet(uuids, tid, cr.config.ContentSourceAppURL, cr.config.ContentSourceAppName)
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

	imgModelsList, err := cr.doGet(imgModelUUIDs, tid, cr.config.ContentSourceAppURL, cr.config.ContentSourceAppName)
	if err != nil {
		return cm, err
	}

	for _, i := range imgModelsList {
		cr.addItemToMap(i, cm)
	}

	return cm, nil
}

// GetInternal internal components from document-store-api
func (cr *ContentReader) GetInternal(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)

	internalContent, err := cr.doGet(uuids, tid, cr.config.InternalContentSourceAppURL, cr.config.InternalContentSourceAppName)
	if err != nil {
		return cm, err
	}

	for _, c := range internalContent {
		uuid, ok := c["uuid"].(string)
		if !ok {
			log.Printf("Cannot extract uuid for content: %v", c)
			continue
		}
		cm[uuid] = c
	}

	return cm, nil
}

func (cr *ContentReader) doGet(uuids []string, tid string, url string, appName string) ([]Content, error) {
	var cb []Content

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return cb, errors.Wrapf(err, "Error connecting to %v", appName)
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
		return cb, errors.Wrapf(err, "Request to %v failed.", appName)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cb, errors.Errorf("Request to %v failed with status code %d", appName, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cb, errors.Wrapf(err, "Error reading response received from %v", appName)
	}

	err = json.Unmarshal(body, &cb)
	if err != nil {
		return cb, errors.Wrapf(err, "Error unmarshalling response from %v", appName)
	}
	return cb, nil
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
