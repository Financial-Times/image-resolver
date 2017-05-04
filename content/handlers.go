package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/satori/go.uuid"
	"io/ioutil"
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

	var article Content
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &article)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("The given json content is not valid, err=%v", err)})
		if errm != nil {
			w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
		} else {
			w.Write([]byte(msg))
		}
		return
	}

	contentUUID := article[ID].(string)
	id := hh.Service.ExtractIdfromUrl(contentUUID)
	err = validateUuid(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg, errm := json.Marshal(ErrorMessage{fmt.Sprintf("The uuid =%s is not valid, err=%v", id, err)})
		if errm != nil {
			w.Write([]byte(fmt.Sprintf("Error message couldn't be encoded in json: , err=%s", errm.Error())))
		} else {
			w.Write([]byte(msg))
		}
		return
	}

	unrolledContent := hh.Service.UnrollImages(article)

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(unrolledContent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg, _ := json.Marshal(ErrorMessage{fmt.Sprintf("Error parsing result for content with id %s, err=%v", id, err)})
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
