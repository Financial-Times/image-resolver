package content

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Financial-Times/transactionid-utils-go"
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
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}
	uuid := extractUUIDFromString(id)
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
		handleError(r, tid, "", w, err, http.StatusBadRequest)
		return
	}
	uuid := extractUUIDFromString(id)
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
		errMsg = fmt.Sprintf("Error getting expanding images because supplied content is invalid: %v", err.Error())
	} else if statusCode >= 500 {
		errMsg = fmt.Sprintf("Error getting expanding images for: %v: %v", uuid, err.Error())
	}
	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusInternalServerError, uuid, err.Error())
	w.WriteHeader(statusCode)
	w.Write([]byte(errMsg))
}
