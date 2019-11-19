package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cshep4/swiftest-img/internal/img"
	uuid "github.com/kevinburke/go.uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type store struct {
	downloader s3manager.Downloader
	uploader   s3manager.Uploader
	bucket     string
}

func New(downloader s3manager.Downloader, uploader s3manager.Uploader, bucket string) img.Storer {
	return &store{
		downloader: downloader,
		uploader:   uploader,
		bucket:     bucket,
	}
}

func (s *store) Upload(img []byte) (string, error) {
	fileName := fmt.Sprintf("%s.jpg", uuid.NewV4().String())

	input := &s3manager.UploadInput{
		Body:        bytes.NewReader(img),
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileName),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(http.DetectContentType(img)),
	}
	res, err := s.uploader.Upload(input)
	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func (s *store) Get(uri string) ([]byte, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	path := strings.SplitN(u.Path, "/", 3)
	key := path[len(path)-1]

	file, err := os.Create(key)
	if err != nil {
		return nil, err
	}

	numBytes, err := s.downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)

	if err != nil || numBytes == 0 {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return b, nil
}
