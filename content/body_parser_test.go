package content

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestShouldReturnImages(t *testing.T) {
	var reader Reader
	var parser Parser
	var ir ImageResolver
	var expectedOutput = []string{"http://api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f", "http://api.ft.com/content/71231d3a-13c7-11e7-2ea7-a07ecd9ac73f", "http://api.ft.com/content/0261ea4a-1474-11e7-1e92-847abda1ac65", "http://api.ft.com/content/da0e3d5d-ccf0-3b40-b865-f648189fb849"}
	reader = NewContentReader("", "")
	parser = NewBodyParser()
	ir = *NewImageResolver(&reader, &parser)
	fileBytes, err := ioutil.ReadFile("../test-resources/bodyXml.xml")
	if err != nil {
		assert.Fail(t, "Cannot read test file")
	}
	str := string(fileBytes)
	emImagesUUIDs, err := ir.parser.GetEmbedded(str)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, expectedOutput, emImagesUUIDs, "Response image ids shoud be equal to expected images")
}

func TestBodyEmpty(t *testing.T) {
	var reader Reader
	var parser Parser
	var ir ImageResolver
	var str string
	reader = NewContentReader("", "")
	parser =NewBodyParser()
	ir = *NewImageResolver(&reader, &parser)
	emImagesUUIDs, _ := ir.parser.GetEmbedded(str)
	assert.Equal(t, 0, len(emImagesUUIDs), "Response should return zero images")
}
