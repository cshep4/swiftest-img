package img

import "errors"

type UploadRequest struct {
	Image []byte `json:"image"`
}

type GetRequest struct {
	Uri string `json:"uri"`
}

type UploadResult struct {
	Uri    string         `json:"uri"`
	Result MarkedDocument `json:"result"`
}

type Document struct {
	Image []byte `json:"image"`
}

type MarkedDocument struct {
	Questions []Question `json:"questions"`
	Total     int        `json:"count"`
	Score     int        `json:"score"`
}

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Correct  bool   `json:"correct"`
}

var ErrDocumentNotFound = errors.New("document not found")
