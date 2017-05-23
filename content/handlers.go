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

func (hh *Handler) GetContentImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	tid := transactionidutils.GetTransactionIDFromRequest(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(r, tid, w, err)
		return
	}

	var article Content
	err = json.Unmarshal(b, &article)
	if err != nil {
		handleBadRequest(r, tid, w, err)
		return
	}

	id, ok := article[id].(string)
	if !ok {
		handleBadRequest(r, tid, w, err)
		return
	}
	uuid := extractUUIDFromURL(id)
	logger.TransactionStartedEvent(r.RequestURI, tid, uuid)

	//unrolling images
	c, err := hh.Service.UnrollImages(article)
	if err != nil {
		handleInternalServerError(r, tid, uuid, w, err)
		return
	}
	jsonRes, err := json.Marshal(c)
	if err != nil {
		handleInternalServerError(r, tid, uuid, w, err)
		return
	}

	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusOK, uuid, "success")
	w.Write(jsonRes)
}

func (hh *Handler) GetLeadImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	tid := transactionidutils.GetTransactionIDFromRequest(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(r, tid, w, err)
		return
	}

	var article Content
	err = json.Unmarshal(b, &article)
	if err != nil {
		handleBadRequest(r, tid, w, err)
		return
	}

	id, ok := article[id].(string)
	if !ok {
		handleBadRequest(r, tid, w, err)
		return
	}
	uuid := extractUUIDFromURL(id)
	logger.TransactionStartedEvent(r.RequestURI, tid, uuid)

	//unrolling lead images
	c, err := hh.Service.UnrollLeadImages(article)
	if err != nil {
		handleInternalServerError(r, tid, uuid, w, err)
		return
	}
	jsonRes, err := json.Marshal(c)
	if err != nil {
		handleInternalServerError(r, tid, uuid, w, err)
		return
	}

	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusOK, uuid, "success")
	w.Write(jsonRes)
}

func handleBadRequest(r *http.Request, tid string, w http.ResponseWriter, err error) {
	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusBadRequest, "", err.Error())
	w.WriteHeader(http.StatusBadRequest)
	msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("Error reading content, err=%v", err)})
	if errm != nil {
		w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
	} else {
		w.Write([]byte(msg))
	}
}

func handleInternalServerError(r *http.Request, tid string, uuid string, w http.ResponseWriter, err error) {
	logger.TransactionFinishedEvent(r.RequestURI, tid, http.StatusInternalServerError, uuid, err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("Error parsing result for content with id %s, err=%v", uuid, err)})
	if errm != nil {
		w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
	} else {
		w.Write([]byte(msg))
	}
}
