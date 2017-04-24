package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

type ServiceConfig struct {
	AppPort             string
	Content_public_read string
	RouterAddress       string
	GraphiteTCPAddress  string
	GraphitePrefix      string
	HttpClient          *http.Client
}

type ContentHandler struct {
	ServiceConfig *ServiceConfig
	Service       Resolver
}

func (hh *ContentHandler) GetContentImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	contentUUID := vars["uuid"]
	err := validateUuid(contentUUID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("The given uuid is not valid, err=%v", err)})
		if errm != nil {
			w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
		} else {
			w.Write([]byte(msg))
		}
		return
	}

	unrolledContent, found, err := hh.Service.UnrollImages(contentUUID)

	if !found {
		w.WriteHeader(http.StatusNotFound)
		msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("Requested item does not exist %s", contentUUID)})
		if errm != nil {
			w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
		} else {
			w.Write([]byte(msg))
		}
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("Error retrieving images for %s, err=%v", contentUUID, err)})
		if errm != nil {
			w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
		} else {
			w.Write([]byte(msg))
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(unrolledContent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorMessage{fmt.Sprintf("Error parsing result for content with uuid %s, err=%v", contentUUID, err)})
		w.Write([]byte(msg))
	}
}


func validateUuid(contentUUID string) error {
	parsedUUID, err := uuid.FromString(contentUUID)
	if err != nil {
		return err
	}
	if contentUUID != parsedUUID.String() {
		return fmt.Errorf("Parsed UUID (%v) is different than the given uuid (%v).", parsedUUID, contentUUID)
	}
	return nil
}
