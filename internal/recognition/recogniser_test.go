package recognition

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestRecogniser_Extract(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../creds.json")
	require.NoError(t, err)

	client, err := vision.NewImageAnnotatorClient(context.Background())
	require.NoError(t, err)

	recogniser := New(*client)

	text, err := recogniser.Extract("https://testies-document.s3.amazonaws.com/testies.jpg")
	require.NoError(t, err)

	log.Println(text)
}
