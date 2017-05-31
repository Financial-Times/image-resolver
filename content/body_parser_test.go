package content

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const ImageSetType = "http://www.ft.com/ontology/content/ImageSet"

func TestShouldReturnImages(t *testing.T) {
	var expectedOutput = []imageSetUUID{
		{
			uuid:           "639cd952-149f-11e7-2ea7-a07ecd9ac73f",
			imageModelUUID: "639cd952-149f-11e7-b0c1-37e417ee6c76",
		},
		{
			uuid:           "71231d3a-13c7-11e7-2ea7-a07ecd9ac73f",
			imageModelUUID: "71231d3a-13c7-11e7-b0c1-37e417ee6c76",
		},
		{
			uuid:           "0261ea4a-1474-11e7-1e92-847abda1ac65",
			imageModelUUID: "0261ea4a-1474-11e7-80f4-13e067d5072c",
		},
	}
	fileBytes, err := ioutil.ReadFile("../test-resources/bodyXml.xml")
	if err != nil {
		assert.Fail(t, "Cannot read test file")
	}
	str := string(fileBytes)
	emImagesUUIDs, err := getEmbedded(str, ImageSetType, "", "")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, expectedOutput, emImagesUUIDs, "Response image ids shoud be equal to expected images")
}

func TestBodyNoEmbeddedImagesReturnsEmptyList(t *testing.T) {
	emImagesUUIDs, err := getEmbedded("<body><p>Sample body</p></body>", ImageSetType, "", "")
	assert.NoError(t, err, "Body parsing should be successful")
	assert.Len(t, emImagesUUIDs, 0, "Response image ids shoud be equal to expected images")
}

func TestMalformedBodyReturnsEmptyList(t *testing.T) {
	emImagesUUIDs, err := getEmbedded("Sample body", ImageSetType, "", "")
	assert.NoError(t, err, "Body parsing should be successful")
	assert.Len(t, emImagesUUIDs, 0, "Response image ids shoud be equal to expected images")
}

func TestEmptyBodyReturnsEmptyList(t *testing.T) {
	emImagesUUIDs, _ := getEmbedded("", ImageSetType, "", "")
	assert.Equal(t, 0, len(emImagesUUIDs), "Response should return zero images")
}
