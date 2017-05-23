package content

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	ID         = "http://www.ft.com/thing/22c0d426-1466-11e7-b0c1-37e417ee6c76"
	expectedId = "22c0d426-1466-11e7-b0c1-37e417ee6c76"
)

type ReaderMock struct {
	mockGet func(c UUIDBatch) (map[string]Content, error)
}

func (rm *ReaderMock) Get(c UUIDBatch) (map[string]Content, error) {
	return rm.mockGet(c)
}

type ParserMock struct {
	mockGetEmbedded func(body string) ([]imageSetUUID, error)
}

func (pm *ParserMock) GetEmbedded(body string) ([]imageSetUUID, error) {
	return pm.mockGetEmbedded(body)
}

func TestUnrollImages(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c UUIDBatch) (map[string]Content, error) {
				b, err := ioutil.ReadFile(testResourcesRoot + "valid-content-reader-response.json")
				assert.NoError(t, err, "Cannot open file necessary for test case")
				var res map[string]Content
				err = json.Unmarshal(b, &res)
				assert.NoError(t, err, "Cannot return valid response")
				return res, nil
			},
		},
		parser: &ParserMock{
			mockGetEmbedded: func(body string) ([]imageSetUUID, error) {
				return []imageSetUUID{
					{uuid: "639cd952-149f-11e7-2ea7-a07ecd9ac73f", imageModelUUID: "639cd952-149f-11e7-b0c1-37e417ee6c76"},
					{uuid: "71231d3a-13c7-11e7-2ea7-a07ecd9ac73f", imageModelUUID: "71231d3a-13c7-11e7-b0c1-37e417ee6c76"},
					{uuid: "0261ea4a-1474-11e7-1e92-847abda1ac65", imageModelUUID: "0261ea4a-1474-11e7-80f4-13e067d5072c"},
				}, nil
			},
		},
		apiHost: "test.api.ft.com",
	}

	expected, err := ioutil.ReadFile("../test-resources/valid-expanded-content-response.json")
	assert.NoError(t, err, "Cannot read necessary test file")

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read necessary test file")
	err = json.Unmarshal(fileBytes, &c)
	assert.NoError(t, err, "Cannot build json body")

	actual, err := ir.UnrollImages(c)
	assert.NoError(t, err, "Should not get an error when expanding images")

	actualJson, err := json.Marshal(actual)
	assert.JSONEq(t, string(actualJson), string(expected))
}

func TestImageResolver_UnrollImages_ErrorWhenReaderReturnsError(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c UUIDBatch) (map[string]Content, error) {
				return nil, errors.New("Cannot retrieve content")
			},
		},
		parser: &ParserMock{
			mockGetEmbedded: func(body string) ([]imageSetUUID, error) {
				return []imageSetUUID{
					{uuid: "639cd952-149f-11e7-2ea7-a07ecd9ac73f", imageModelUUID: "639cd952-149f-11e7-b0c1-37e417ee6c76"},
					{uuid: "71231d3a-13c7-11e7-2ea7-a07ecd9ac73f", imageModelUUID: "71231d3a-13c7-11e7-b0c1-37e417ee6c76"},
					{uuid: "0261ea4a-1474-11e7-1e92-847abda1ac65", imageModelUUID: "0261ea4a-1474-11e7-80f4-13e067d5072c"},
				}, nil
			},
		},
		apiHost: "test.api.ft.com",
	}

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	err = json.Unmarshal(fileBytes, &c)

	_, err = ir.UnrollImages(c)
	assert.Error(t, err)
}

func TestImageResolver_UnrollImages_EmbeddedImagesSkippedWhenParserReturnsError(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c UUIDBatch) (map[string]Content, error) {
				b, err := ioutil.ReadFile(testResourcesRoot + "valid-content-reader-response-no-body.json")
				assert.NoError(t, err, "Cannot open file necessary for test case")
				var res map[string]Content
				err = json.Unmarshal(b, &res)
				assert.NoError(t, err, "Cannot return valid response")
				return res, nil
			},
		},
		parser: &ParserMock{
			mockGetEmbedded: func(body string) ([]imageSetUUID, error) {
				return nil, errors.New("Cannot parse body")
			},
		},
		apiHost: "test.api.ft.com",
	}

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	err = json.Unmarshal(fileBytes, &c)

	result, err := ir.UnrollImages(c)
	assert.NoError(t, err, "Should not receive error when body cannot be parsed.")
	assert.Nil(t, result["embeds"], "Response should not contain embeds field")
}

func TestImageResolver_UnrollLeadImages(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c UUIDBatch) (map[string]Content, error) {
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

	actual, err := ir.UnrollLeadImages(c)
	assert.NoError(t, err, "Should not receive error for expanding lead images")
	actualJson, err := json.Marshal(actual)
	assert.JSONEq(t, string(actualJson), string(expected))
}

func TestImageResolver_UnrollLeadImages_ErrorWhenReaderFails(t *testing.T) {
	ir := ImageResolver{
		reader: &ReaderMock{
			mockGet: func(c UUIDBatch) (map[string]Content, error) {
				return nil, errors.New("Cannot read content")
			},
		},
		apiHost: "test.api.ft.com",
	}

	var c Content
	fileBytes, err := ioutil.ReadFile("../test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "Cannot read test file")
	err = json.Unmarshal(fileBytes, &c)

	_, err = ir.UnrollLeadImages(c)
	assert.Error(t, err)
}

func TestExtractIDFromURL(t *testing.T) {
	actual := extractUUIDFromURL(ID)
	assert.Equal(t, expectedId, actual, "Response id should be equal")
}
