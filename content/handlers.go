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

type ServiceConfig struct {
	AppPort             string
	Content_public_read string
	RouterAddress       string
	EmbedsType          string
	GraphiteTCPAddress  string
	GraphitePrefix      string
	HttpClient          *http.Client
}

type ContentHandler struct {
	ServiceConfig *ServiceConfig
	Service       Resolver
	Log           *appLogger
}

func (hh *ContentHandler) GetContentImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var article Content
	var id string
	tid := transactionidutils.GetTransactionIDFromRequest(r)

	hh.Log.TransactionStartedEvent(r.RequestURI, tid, id)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(hh, r, tid, id, w, err)
		return
	}

	err = json.Unmarshal(b, &article)
	if err != nil {
		handleBadRequest(hh, r, tid, id, w, err)
		return
	}

	contentUUID, ok := article[ID].(string)
	if ok {
		id = extractIdfromUrl(contentUUID)
	}

	hh.Log.RequestEvent(AppName, r.RequestURI, tid, id)
	//unrolling images
	content, err := json.Marshal(hh.Service.UnrollImages(article));
	if err != nil {
		handleInternalServerError(hh, r, tid, id, w, err)
		return
	}

	hh.Log.ResponseEvent(AppName, r.URL.String(), tid, http.StatusOK, id)
	w.Write([]byte(content))
}

func (hh *ContentHandler) GetLeadImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var article Content
	var id string
	tid := transactionidutils.GetTransactionIDFromRequest(r)

	hh.Log.TransactionStartedEvent(r.RequestURI, tid, id)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(hh, r, tid, id, w, err)
		return
	}

	err = json.Unmarshal(b, &article)
	if err != nil {
		handleBadRequest(hh, r, tid, id, w, err)
		return
	}

	contentUUID, ok := article[UUID].(string)
	if ok {
		id = extractIdfromUrl(contentUUID)
	}

	hh.Log.RequestEvent(AppName, r.RequestURI, tid, id)
	//unrolling lead images
	content, err := json.Marshal(hh.Service.UnrollLeadImages(article));
	if  err != nil {
		handleInternalServerError(hh, r, tid, id, w, err)
		return
	}

	hh.Log.ResponseEvent(AppName, r.URL.String(), tid, http.StatusOK, id)
	w.Write([]byte(content))
}


func handleBadRequest(hh *ContentHandler, r *http.Request, tid string, id string, w http.ResponseWriter, err error) {
	hh.Log.RequestFailedEvent(AppName, r.RequestURI, tid, http.StatusBadRequest, id)
	w.WriteHeader(http.StatusBadRequest)
	msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("Error reading content, err=%v", err)})
	if errm != nil {
		w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
	} else {
		w.Write([]byte(msg))
	}
}

func handleInternalServerError(hh *ContentHandler, r *http.Request, tid string, id string, w http.ResponseWriter, err error)  {
	hh.Log.RequestFailedEvent(AppName, r.RequestURI, tid, http.StatusInternalServerError, id)
	w.WriteHeader(http.StatusInternalServerError)
	msg, _ := json.Marshal(ErrorMessage{fmt.Sprintf("Error parsing result for content with id %s, err=%v", id, err)})
	w.Write([]byte(msg))
}
