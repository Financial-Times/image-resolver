package content

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Financial-Times/transactionid-utils-go"
	"github.com/pkg/errors"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

var logger = NewAppLogger()

type Handler struct {
	Service Resolver
}

type UnrollEvent struct {
	c    Content
	tid  string
	uuid string
}

type UnrollResult struct {
	uc  Content
	err error
}

func (hh *Handler) GetContentImages(w http.ResponseWriter, r *http.Request) {
	tid := transactionidutils.GetTransactionIDFromRequest(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}

	var article Content
	err = json.Unmarshal(b, &article)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}

	id, ok := article[id].(string)
	if !ok {
		handleError(r, tid, "", w, errors.New("Missing or invalid id field"), http.StatusBadRequest)
		return
	}
	uuid, err := extractUUIDFromString(id)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}

	if !validateContentImages(article) {
		handleError(r, tid, uuid, w, errors.New("Invalid content"), http.StatusBadRequest)
		return
	}

	logger.TransactionStartedEvent(r.RequestURI, tid, uuid)

	//unrolling images
	req := UnrollEvent{article, tid, uuid}
	res := hh.Service.UnrollImages(req)
	if res.err != nil {
		handleError(r, tid, uuid, w, res.err, http.StatusInternalServerError)
		return
	}
	jsonRes, err := json.Marshal(res.uc)
	if err != nil {
		handleError(r, tid, uuid, w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusOK, uuid, "success")
	w.Write(jsonRes)
}

func (hh *Handler) GetLeadImages(w http.ResponseWriter, r *http.Request) {
	tid := transactionidutils.GetTransactionIDFromRequest(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}

	var article Content
	err = json.Unmarshal(b, &article)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}

	id, ok := article[id].(string)
	if !ok {
		handleError(r, tid, "", w, errors.New("Missing or invalid id field"), http.StatusBadRequest)
		return
	}
	uuid, err := extractUUIDFromString(id)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}

	if !validateInternalContentImages(article) {
		handleError(r, tid, uuid, w, errors.New("Invalid content"), http.StatusBadRequest)
		return
	}

	logger.TransactionStartedEvent(r.RequestURI, tid, uuid)

	//unrolling lead images
	req := UnrollEvent{article, tid, uuid}
	res := hh.Service.UnrollLeadImages(req)
	if res.err != nil {
		handleError(r, tid, uuid, w, res.err, http.StatusInternalServerError)
		return
	}
	jsonRes, err := json.Marshal(res.uc)
	if err != nil {
		handleError(r, tid, uuid, w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusOK, uuid, "success")
	w.Write(jsonRes)
}

func handleError(r *http.Request, tid string, uuid string, w http.ResponseWriter, err error, statusCode int) {
	var errMsg string
	if statusCode >= 400 && statusCode < 500 {
		errMsg = fmt.Sprintf("Error getting expanding images because supplied content is invalid: %s", err.Error())
		logger.Errorf(tid, "Error getting expanding images because supplied content is invalid: %s", err.Error())
	} else if statusCode >= 500 {
		errMsg = fmt.Sprintf("Error getting expanding images for: %v: %v", uuid, err.Error())
		logger.TransactionFinishedEvent(r.RequestURI, tid, statusCode, uuid, err.Error())
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(errMsg))
}

func validateContentImages(article Content) bool {
	_, hasMainImage := article[mainImage]
	_, hasBody := article[bodyXML]
	altImg, hasAltImg := article[altImages].(map[string]interface{})
	var hasPromImg bool
	if hasAltImg {
		_, hasPromImg = altImg[promotionalImage]
	}

	return hasMainImage || hasBody || hasPromImg
}

func validateInternalContentImages(article Content) bool {
	_, hasLeadImages := article[leadImages]
	return hasLeadImages
}
