package content

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"strings"
	"io/ioutil"
)

var contentAPIMock *httptest.Server

const (
	Content_Id   = "http://www.ft.com/thing/22c0d426-1466-11e7-b0c1-37e417ee6c76"
	Type_Art     = "http://www.ft.com/ontology/content/Article"
	Image_ID     = "http://www.ft.com/thing/639cd952-149f-11e7-2ea7-a07ecd9ac73f"
	Image_Date   = "2017-03-29T19:39:00.000Z"
	Publish_Date = "publishedDate"
	Expected_Id  = "22c0d426-1466-11e7-b0c1-37e417ee6c76"
)

func startContentAPIMock(status string) {
	router := mux.NewRouter()
	var getContent http.HandlerFunc

	if status == "happy" {
		getContent = happyContentAPIMock

	} else if status == "notFound" {
		getContent = notFoundHandler

	} else {
		getContent = internalErrorHandler

	}
	router.Path("/content/{uuid}").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(getContent)})
	contentAPIMock = httptest.NewServer(router)
}

func happyContentAPIMock(writer http.ResponseWriter, request *http.Request) {
	file, err := os.Open("../test-resources/image.json")
	if err != nil {
		return
	}
	defer file.Close()
	io.Copy(writer, file)
}

func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func serviceIR() ImageResolver {
	contentAPIURI := contentAPIMock.URL + "/content/"
	router := strings.Replace(contentAPIMock.URL, "http://", "", -1)

	var reader Reader
	var parser Parser
	var ir ImageResolver

	reader = NewContentReader(contentAPIURI, router)
	parser = BodyParser{}
	ir = *NewImageResolver(&reader, &parser)
	return ir
}

func TestUnrollImages(t *testing.T) {
	var content Content
	startContentAPIMock("happy")
	defer contentAPIMock.Close()
	resp, err := http.Get(contentAPIMock.URL + "/content/22c0d426-1466-11e7-b0c1-37e417ee6c76")
	assert.NoError(t, err, "Cannot send request to content endpoint")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")

	fileBytes, err := ioutil.ReadFile("../test-resources/content.json")
	if err != nil {
		assert.Fail(t, "Cannot read test file")
	}
	err = json.Unmarshal(fileBytes, &content)

	ir := serviceIR()
	result:= ir.UnrollImages(content)

	assert.Equal(t, Content_Id, result[ID], "Response ID  shoud be equal")
	assert.Equal(t, Type_Art, result[Type], "Response Type  shoud be equal")
	assert.Equal(t, Image_ID, result[MainImage].(Content)[ID], "Response Main Image Id  shoud be equal")
	assert.Equal(t, Image_Date, result[MainImage].(Content)["publishedDate"], "Response Main image publishedDate shoud be equal")
	assert.Equal(t, 4, len(result[Embeds].([]Content)), "Response Embeds length shoud be equal 4")
	img := result[AltImages].(PromotionalImage)
	promo := img.PromotionalImage.(Content)
	assert.Equal(t, Image_ID, promo[ID], "Response Promotional Image Id  shoud be equal")
	assert.Equal(t, Image_Date, promo[Publish_Date], "Response Promotional Image publishedDate shoud be equal")
	lead := result[LeadImages].([]ImageOutput)
	assert.Equal(t, 3, len(lead), "Response LeadImages length shoud be equal 3")
}

func TestExtractIdfromUrl(t *testing.T){
	ir := serviceIR()
	actual := ir.ExtractIdfromUrl(Content_Id)
	assert.Equal(t, Expected_Id, actual, "Response Embeds length shoud be equal")
}
