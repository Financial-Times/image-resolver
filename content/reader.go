package content

import (
	"encoding/json"
	"fmt"
	"io"
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
	GetInternal([]string, string) (map[string]Content, error)
	GetNative([]string, string) (map[string]Content, error)
	GetNativeInternal([]string, string) (map[string]Content, error)
}

type ReaderFunc func([]string, string) (map[string]Content, error)

type ReaderConfig struct {
	ContentSourceAppName                  string
	ContentSourceAppURL                   string
	ContentSourceInternalURL              string
	NativeContentSourceAppName            string
	NativeContentSourceAppURL             string
	NativeContentSourceAppAuth            string
	TransformContentSourceURL             string
	TransformContentSourceAppName         string
	TransformInternalContentSourceURL     string
	TransformInternalContentSourceAppName string
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

	contentBatch, err := cr.doGet(uuids, tid, cr.config.ContentSourceAppURL, cr.config.ContentSourceAppName)
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

	internalContent, err := cr.doGet(uuids, tid, cr.config.ContentSourceInternalURL, cr.config.ContentSourceAppName)
	if err != nil {
		return cm, err
	}

	for _, c := range internalContent {
		cr.addItemToMap(c, cm)
	}

	return cm, nil
}

// GetNative from Methode API
func (cr *ContentReader) GetNative(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)

	for _, uuid := range uuids {
		nativeResponse, err := cr.doGetNative(uuid, tid)
		if err != nil {
			return nil, err
		}

		transformedBody, err := cr.doGetTransformedContent(nativeResponse.Body, cr.config.TransformContentSourceURL, cr.config.TransformContentSourceAppName, uuid, tid)
		if err != nil {
			return nil, err
		}

		defer nativeResponse.Body.Close()

		cm[uuid] = transformedBody
	}

	return cm, nil
}

// GetNativeInternal reads internalcomponents from Methode API
func (cr *ContentReader) GetNativeInternal(uuids []string, tid string) (map[string]Content, error) {
	var cm = make(map[string]Content)

	for _, uuid := range uuids {
		nativeResponse, err := cr.doGetNative(uuid, tid)
		if err != nil {
			return nil, err
		}

		transformedBody, err := cr.doGetTransformedContent(nativeResponse.Body, cr.config.TransformInternalContentSourceURL, cr.config.TransformInternalContentSourceAppName, uuid, tid)
		if err != nil {
			return nil, err
		}

		defer nativeResponse.Body.Close()

		cm[uuid] = transformedBody
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

func (cr *ContentReader) doGetNative(uuid string, tid string) (*http.Response, error) {
	requestURL := fmt.Sprintf("%s%s", cr.config.NativeContentSourceAppURL, uuid)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Error connecting to %v for uuid: %s", cr.config.NativeContentSourceAppName, uuid)
	}

	req.Header.Set(transactionidutils.TransactionIDHeader, tid)
	req.Header.Set("Authorization", "Basic "+cr.config.NativeContentSourceAppAuth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgentValue)

	res, err := cr.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Request to %v failed for uuid: %s", cr.config.NativeContentSourceAppName, uuid)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Request to %v failed for uuid: %s with status code %d", cr.config.NativeContentSourceAppName, uuid, res.StatusCode)
	}

	return res, nil
}

func (cr *ContentReader) doGetTransformedContent(nativeContent io.Reader, transformerURL string, trasformerSourceApp string, uuid string, tid string) (Content, error) {
	req, err := http.NewRequest(http.MethodPost, transformerURL, nativeContent)
	if err != nil {
		return nil, errors.Wrapf(err, "Error connecting to %v for uuid: %s", trasformerSourceApp, uuid)
	}

	req.Header.Add(transactionidutils.TransactionIDHeader, tid)
	req.Header.Set(userAgent, userAgentValue)
	req.Header.Set("Content-Type", "application/json")

	res, err := cr.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Request to %v failed for uuid: %s", trasformerSourceApp, uuid)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Request to %v failed for uuid: %s with status code %d", trasformerSourceApp, uuid, res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading response received from %v", trasformerSourceApp)
	}

	var c Content
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, errors.Wrapf(err, "Error unmarshalling response from %v for uuid %s", trasformerSourceApp, uuid)
	}

	return c, nil
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
