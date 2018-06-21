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

func (hh *Handler) GetContent(w http.ResponseWriter, r *http.Request) {
	tid := transactionidutils.GetTransactionIDFromRequest(r)
	event, err := createUnrollEvent(r, tid)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
	}

	if !validateContentImages(event.c) {
		handleError(r, tid, event.uuid, w, errors.New("Invalid content"), http.StatusBadRequest)
		return
	}

	logger.TransactionStartedEvent(r.RequestURI, tid, event.uuid)

	//unrolling images
	res := hh.Service.UnrollImages(event)
	if res.err != nil {
		handleError(r, tid, event.uuid, w, res.err, http.StatusInternalServerError)
		return
	}
	jsonRes, err := json.Marshal(res.uc)
	if err != nil {
		handleError(r, tid, event.uuid, w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusOK, event.uuid, "success")
	w.Write(jsonRes)
}

func (hh *Handler) GetInternalContent(w http.ResponseWriter, r *http.Request) {
	tid := transactionidutils.GetTransactionIDFromRequest(r)
	event, err := createUnrollEvent(r, tid)
	if err != nil {
		handleError(r, tid, "", w, err, http.StatusBadRequest)
	}

	if !validateInternalContentImages(event.c) {
		handleError(r, tid, event.uuid, w, errors.New("Invalid content"), http.StatusBadRequest)
		return
	}

	logger.TransactionStartedEvent(r.RequestURI, tid, event.uuid)

	//unrolling lead images
	res := hh.Service.UnrollInternalContent(event)
	if res.err != nil {
		handleError(r, tid, event.uuid, w, res.err, http.StatusInternalServerError)
		return
	}
	jsonRes, err := json.Marshal(res.uc)
	if err != nil {
		handleError(r, tid, event.uuid, w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusOK, event.uuid, "success")
	w.Write(jsonRes)
}

func (hh *Handler) GetContentPreview(w http.ResponseWriter, r *http.Request) {
}

func (hh *Handler) GetInternalContentPreview(w http.ResponseWriter, r *http.Request) {
}

func createUnrollEvent(r *http.Request, tid string) (UnrollEvent, error) {
	var unrollEvent UnrollEvent
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return unrollEvent, err
	}

	var article Content
	err = json.Unmarshal(b, &article)
	if err != nil {
		return unrollEvent, err
	}

	id, ok := article[id].(string)
	if !ok {
		return unrollEvent, errors.New("Missing or invalid id field")
	}
	uuid, err := extractUUIDFromString(id)
	if err != nil {
		return unrollEvent, err
	}
	unrollEvent = UnrollEvent{article, tid, uuid}

	return unrollEvent, nil
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
	_, hasAltImg := article[altImages].(map[string]interface{})

	return hasMainImage || hasBody || hasAltImg
}

func validateInternalContentImages(article Content) bool {
	_, hasLeadImages := article[leadImages]
	return hasLeadImages
}
