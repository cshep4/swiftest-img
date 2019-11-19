package transport

import (
	"encoding/json"
	"github.com/cshep4/swiftest-img/internal/img"
	"github.com/cshep4/swiftest-img/internal/score"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type httpHandler struct {
	service img.Servicer
}

func NewHttpHandler(service img.Servicer) (*httpHandler, error) {
	return &httpHandler{
		service: service,
	}, nil
}

func (h *httpHandler) Route() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/document", h.getDocumentResult).
		Methods(http.MethodGet)
	router.HandleFunc("/document", h.uploadDocument).
		Methods(http.MethodPost)

	return router
}

func (h *httpHandler) getDocumentResult(w http.ResponseWriter, r *http.Request) {
	var req img.GetRequest
	defer r.Body.Close()

	if err := h.decodeRequest(r.Body, req); err != nil {
		log.Printf("cannot decode request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.GetDocument(req.Uri)
	h.sendResponse(res, err, w)
}

func (h *httpHandler) uploadDocument(w http.ResponseWriter, r *http.Request) {
	var req img.UploadRequest
	defer r.Body.Close()

	if err := h.decodeRequest(r.Body, req); err != nil {
		log.Printf("cannot decode request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.Upload(req)
	h.sendResponse(res, err, w)
}

func (h *httpHandler) decodeRequest(body io.ReadCloser, req interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &req); err != nil {
		return err
	}

	return nil
}

func (h *httpHandler) sendResponse(data interface{}, err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err.Error())

		switch err {
		case img.ErrDocumentNotFound:
			w.WriteHeader(http.StatusNotFound)
		case score.ErrUnsupportedOperator:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		log.Println("cannot encode response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
