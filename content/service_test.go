package content

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	ID         = "http://www.ft.com/thing/22c0d426-1466-11e7-b0c1-37e417ee6c76"
	expectedId = "22c0d426-1466-11e7-b0c1-37e417ee6c76"
)

type ReaderMock struct {
	mockGet func(c []string) (map[string]Content, error)
}

func (rm *ReaderMock) Get(c []string) (map[string]Content, error) {
	return rm.mockGet(c)
}

func TestUnrollImages(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c []string) (map[string]Content, error) {
				b, err := ioutil.ReadFile(testResourcesRoot + "valid-content-reader-response.json")
				assert.NoError(t, err, "Cannot open file necessary for test case")
				var res map[string]Content
				err = json.Unmarshal(b, &res)
				assert.NoError(t, err, "Cannot return valid response")
				return res, nil
			},
		},
		whitelist: ImageSetType,
		apiHost:   "test.api.ft.com",
	}

	expected, err := ioutil.ReadFile("../test-resources/valid-expanded-content-response.json")
	assert.NoError(t, err, "Cannot read necessary test file")

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read necessary test file")
	err = json.Unmarshal(fileBytes, &c)
	assert.NoError(t, err, "Cannot build json body")
	req := UnrollEvent{c, "tid_sample", "sample_uuid"}
	actual := ir.UnrollImages(req)
	assert.NoError(t, actual.err, "Should not get an error when expanding images")

	actualJson, err := json.Marshal(actual.uc)
	assert.JSONEq(t, string(actualJson), string(expected))
}

func TestImageResolver_UnrollImages_ErrorWhenReaderReturnsError(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c []string) (map[string]Content, error) {
				return nil, errors.New("Cannot retrieve content")
			},
		},
		whitelist: ImageSetType,
		apiHost:   "test.api.ft.com",
	}

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	err = json.Unmarshal(fileBytes, &c)

	req := UnrollEvent{c, "tid_sample", "sample_uuid"}
	actual := ir.UnrollImages(req)
	assert.Error(t, actual.err)
}

func TestImageResolver_UnrollImages_EmbeddedImagesSkippedWhenParserReturnsError(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c []string) (map[string]Content, error) {
				b, err := ioutil.ReadFile(testResourcesRoot + "valid-content-reader-response-no-body.json")
				assert.NoError(t, err, "Cannot open file necessary for test case")
				var res map[string]Content
				err = json.Unmarshal(b, &res)
				assert.NoError(t, err, "Cannot return valid response")
				return res, nil
			},
		},
		whitelist: ImageSetType,
		apiHost:   "test.api.ft.com",
	}

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	err = json.Unmarshal(fileBytes, &c)
	c[bodyXML] = "invalid body"

	req := UnrollEvent{c, "tid_sample", "sample_uuid"}
	res := ir.UnrollImages(req)
	assert.NoError(t, res.err, "Should not receive error when body cannot be parsed.")
	assert.Nil(t, res.uc["embeds"], "Response should not contain embeds field")
}

func TestImageResolver_UnrollLeadImages(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c []string) (map[string]Content, error) {
				b, err := ioutil.ReadFile(testResourcesRoot + "valid-internalcontent-reader-response.json")
				assert.NoError(t, err, "Cannot open file necessary for test case")
				var res map[string]Content
				err = json.Unmarshal(b, &res)
				assert.NoError(t, err, "Cannot return valid response")
				return res, nil
			},
		},
		apiHost: "test.api.ft.com",
	}

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "File necessary for building request body nod found")
	err = json.Unmarshal(fileBytes, &c)

	expected, err := ioutil.ReadFile("../test-resources/valid-expanded-internalcontent-response.json")
	assert.NoError(t, err, "File necessary for building expected output not found.")

	req := UnrollEvent{c, "tid_sample", "sample_uuid"}
	actual := ir.UnrollLeadImages(req)
	assert.NoError(t, actual.err, "Should not receive error for expanding lead images")
	actualJson, err := json.Marshal(actual.uc)
	assert.JSONEq(t, string(actualJson), string(expected))
}

func TestImageResolver_UnrollLeadImages_ErrorWhenReaderFails(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c []string) (map[string]Content, error) {
				return nil, errors.New("Cannot read content")
			},
		},
		apiHost: "test.api.ft.com",
	}

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "Cannot read test file")
	err = json.Unmarshal(fileBytes, &c)

	req := UnrollEvent{c, "tid_sample", "sample_uuid"}
	actual := ir.UnrollLeadImages(req)
	assert.Error(t, actual.err)
}

func TestExtractIDFromURL(t *testing.T) {
	actual, err := extractUUIDFromString(ID)
	assert.NoError(t, err, "Test should not return error")
	assert.Equal(t, expectedId, actual, "Response id should be equal")
}
