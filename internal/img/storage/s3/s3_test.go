package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestS3Store_Upload(t *testing.T) {
	region := "us-east-1"
	bucket := "testies-document"

	awsAccessKey := ""
	awsSecret := ""
	token := ""

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, token)

	_, err := creds.Get()
	require.NoError(t, err)

	cfg := aws.NewConfig().
		WithRegion(region).
		WithCredentials(creds)

	sess, err := session.NewSession(cfg)
	require.NoError(t, err)

	imgfile, err := os.Open("./screenshot.png")
	require.NoError(t, err)

	defer imgfile.Close()

	b, err := ioutil.ReadAll(imgfile)
	require.NoError(t, err)

	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)

	store := New(*downloader, *uploader, bucket)

	res, err := store.Upload(b)
	require.NoError(t, err)

	log.Println(res)
}
