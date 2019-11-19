package recognition

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	vision2 "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type Recogniser interface {
	Extract(uri string) (string, error)
}

type recogniser struct {
	client vision.ImageAnnotatorClient
}

func New(client vision.ImageAnnotatorClient) Recogniser {
	return &recogniser{
		client: client,
	}
}

func (r *recogniser) Extract(uri string) (string, error) {
	annotations, err := r.client.DetectDocumentText(
		context.Background(),
		vision.NewImageFromURI(uri),
		&vision2.ImageContext{
			LanguageHints: []string{"en-t-i0-handwrit"},
		},
	)
	if err != nil {
		return "", err
	}

	return annotations.Text, nil
}
