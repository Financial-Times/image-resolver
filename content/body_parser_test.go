package content

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnImages(t *testing.T){
	var result Content
	var reader Reader
	var parser Parser
	var ir ImageResolver
	var expectedOutput = []string{"639cd952-149f-11e7-2ea7-a07ecd9ac73f", "71231d3a-13c7-11e7-2ea7-a07ecd9ac73f", "0261ea4a-1474-11e7-1e92-847abda1ac65", "da0e3d5d-ccf0-3b40-b865-f648189fb849"}
	reader = NewContentReader("", "")
	parser = BodyParser{}
	ir = *NewImageResolver(&reader, &parser)
	fileBytes, err := ioutil.ReadFile("../test-resources/bodyXml.xml")
	if err != nil {
		assert.Fail(t, "Cannot read test file")
	}
	result.BodyXML = string(fileBytes)
	emImagesUUIDs, err := ir.parser.GetEmbedded(result)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, expectedOutput, emImagesUUIDs, "Response image ids shoud be equal to expected images")
}


func TestBodyEmpty(t *testing.T){
	var result Content
	var reader Reader
	var parser Parser
	var ir ImageResolver
	reader = NewContentReader("", "")
	parser = BodyParser{}
	ir = *NewImageResolver(&reader, &parser)
	result.BodyXML = ""
	_, err := ir.parser.GetEmbedded(result)
	assert.Equal(t, "Cannot parse empty body of content []", err.Error(), "Response should return empty body error")
}