package service

import (
	"github.com/cshep4/swiftest-img/internal/img"
	"github.com/cshep4/swiftest-img/internal/recognition"
	"github.com/cshep4/swiftest-img/internal/score"
	"log"
)

type service struct {
	store      img.Storer
	recogniser recognition.Recogniser
	scorer     score.Scorer
}

func New(store img.Storer, recogniser recognition.Recogniser, scorer score.Scorer) img.Servicer {
	return &service{
		store:      store,
		recogniser: recogniser,
		scorer:     scorer,
	}
}

func (s *service) Upload(req img.UploadRequest) (*img.UploadResult, error) {
	uri, err := s.store.Upload(req.Image)
	if err != nil {
		return nil, err
	}

	docText, err := s.recogniser.Extract(uri)
	if err != nil {
		return nil, err
	}

	log.Printf("Recognised text: %s", docText)

	marks, err := s.scorer.Grade(docText)
	if err != nil {
		return nil, err
	}

	return &img.UploadResult{
		Uri:    uri,
		Result: *marks,
	}, nil
}

func (s *service) GetDocument(uri string) (*img.Document, error) {
	image, err := s.store.Get(uri)
	if err != nil {
		return nil, err
	}

	return &img.Document{
		Image: image,
	}, nil
}
