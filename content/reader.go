package content

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

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
	GetInternal([]string, string) (map[string]Content, error)
	GetPreview([]string, string) (map[string]Content, error)
	GetInternalPreview([]string, string) (map[string]Content, error)
}

type ReaderFunc func([]string, string) (map[string]Content, error)

type ReaderConfig struct {
	ContentStoreAppName         string
	ContentStoreHost            string
	ContentPreviewAppName       string
	ContentPreviewHost          string
	ContentPathEndpoint         string
	InternalContentPathEndpoint string
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

// Get reads content from content-public-read
func (cr *ContentReader) Get(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)
	requestURL := fmt.Sprintf("%s%s", cr.config.ContentStoreHost, cr.config.ContentPathEndpoint)

	contentBatch, err := cr.doGet(uuids, tid, requestURL, cr.config.ContentStoreAppName)
	if err != nil {
		return cm, err
	}

	var imgModelUUIDs []string
	for _, c := range contentBatch {
		cr.addItemToMap(c, cm)
		if _, foundMembers := c[members]; foundMembers {
			imgModelUUIDs = append(imgModelUUIDs, c.getMembersUUID()...)
		}
	}

	if len(imgModelUUIDs) == 0 {
		return cm, nil
	}

	imgModelsList, err := cr.doGet(imgModelUUIDs, tid, requestURL, cr.config.ContentStoreAppName)
	if err != nil {
		return cm, err
	}

	for _, i := range imgModelsList {
		cr.addItemToMap(i, cm)
	}

	return cm, nil
}

// GetInternal reads internal components from content-public-read
func (cr *ContentReader) GetInternal(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)
	requestURL := fmt.Sprintf("%s%s", cr.config.ContentStoreHost, cr.config.InternalContentPathEndpoint)

	internalContent, err := cr.doGet(uuids, tid, requestURL, cr.config.ContentStoreAppName)
	if err != nil {
		return cm, err
	}

	for _, c := range internalContent {
		cr.addItemToMap(c, cm)
	}

	return cm, nil
}

// GetPreview reads content from Content-Preview API
func (cr *ContentReader) GetPreview(uuids []string, tid string) (map[string]Content, error) {
	return cr.getPreviewAsync(uuids, tid, false)
}

// GetInternalPreview reads internalcomponents from Internal-Content-Preview API
func (cr *ContentReader) GetInternalPreview(uuids []string, tid string) (map[string]Content, error) {
	return cr.getPreviewAsync(uuids, tid, true)
}

func (cr *ContentReader) getPreviewAsync(uuids []string, tid string, isInternalPreview bool) (map[string]Content, error) {
	cm := make(map[string]Content)
	ch := make(chan Content, len(uuids))

	var wg sync.WaitGroup
	wg.Add(len(uuids))

	for _, uuid := range uuids {
		go func(uuid string, tid string, cr *ContentReader) {
			requestURL := cr.createPreviewRequestURL(uuid, isInternalPreview)
			content, err := cr.doGetPreview(uuid, tid, requestURL, cr.config.ContentPreviewAppName)

			if err != nil {
				logger.Errorf(tid, "Error while expanding content %s", err.Error())
			} else {
				ch <- content
			}
			wg.Done()
		}(uuid, tid, cr)
	}

	wg.Wait()
	close(ch)

	for {
		content, ok := <-ch
		if !ok {
			break
		}
		cr.addItemToMap(content, cm)
	}

	return cm, nil
}

func (cr *ContentReader) doGet(uuids []string, tid string, reqURL string, appName string) ([]Content, error) {
	var cb []Content

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
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

func (cr *ContentReader) doGetPreview(uuid string, tid string, reqURL string, appName string) (Content, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Error connecting to %v for uuid: %s", appName, uuid)
	}

	req.Header.Set(transactionidutils.TransactionIDHeader, tid)
	req.Header.Set("User-Agent", userAgentValue)

	res, err := cr.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Request to %v failed for uuid: %s", appName, uuid)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Request to %v failed for uuid: %s with status code %d", appName, uuid, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading response received from %v", appName)
	}

	var content Content
	err = json.Unmarshal(body, &content)
	if err != nil {
		return content, errors.Wrapf(err, "Error unmarshalling response from %v", appName)
	}
	return content, nil
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

func (cr *ContentReader) createPreviewRequestURL(uuid string, isInternalPreview bool) string {
	if isInternalPreview {
		return fmt.Sprintf("%s%s/%s", cr.config.ContentPreviewHost, cr.config.InternalContentPathEndpoint, uuid)
	}

	return fmt.Sprintf("%s%s/%s", cr.config.ContentPreviewHost, cr.config.ContentPathEndpoint, uuid)
}
