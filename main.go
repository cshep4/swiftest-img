package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cshep4/swiftest-img/internal/img"
	"github.com/cshep4/swiftest-img/internal/img/service"
	"github.com/cshep4/swiftest-img/internal/img/storage/s3"
	"github.com/cshep4/swiftest-img/internal/recognition"
	"github.com/cshep4/swiftest-img/internal/score"
	"github.com/cshep4/swiftest-img/internal/transport"
	"google.golang.org/grpc/codes"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

func main() {
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		log.Fatal("aws region not set")
	}
	bucket, ok := os.LookupEnv("S3_BUCKET")
	if !ok {
		log.Fatal("s3 bucket not set")
	}
	awsAccessKey, ok := os.LookupEnv("AWS_ACCESS_KEY")
	if !ok {
		log.Fatal("aws access key not set")
	}
	awsSecret, ok := os.LookupEnv("AWS_SECRET_KEY")
	if !ok {
		log.Fatal("s3 secret key not set")
	}

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, "")

	_, err := creds.Get()
	if err != nil {
		log.Fatalf("failed to authenticate aws: %v", err)
	}

	cfg := aws.NewConfig().
		WithRegion(region).
		WithCredentials(creds)

	sess, err := session.NewSession(cfg)
	if err != nil {
		log.Fatalf("failed to create aws session: %v", err)
	}

	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 200 * 1024 * 1024 // 64MB per part
	})

	scorer := score.New()

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "creds.json")
	if err != nil {
		log.Fatalf("failed to set GCP credentials: %v", err)
	}

	client, err := vision.NewImageAnnotatorClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create image annotator client: %v", err)
	}

	recogniser := recognition.New(*client)
	store := s3.New(*downloader, *uploader, bucket)
	svc := service.New(store, recogniser, scorer)

	var exitCode codes.Code
	var httpServer *http.Server

	sigs := make(chan os.Signal)

	go func() {
		httpServer = startHttpServer(svc)
		sigs <- syscall.SIGQUIT
	}()

	switch sig := <-sigs; sig {
	case os.Interrupt, syscall.SIGINT, syscall.SIGQUIT:
		log.Print("Shutting down")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Printf("Error shutting down http server: %v\n", err)
		}

		exitCode = codes.Aborted
	case syscall.SIGTERM:
		exitCode = codes.OK
	}

	os.Exit(int(exitCode))
}

func startHttpServer(service img.Servicer) *http.Server {
	h, err := transport.NewHttpHandler(service)
	if err != nil {
		log.Fatalf("failed to create http handler: %v", err)
	}

	path := ":8080"

	http := &http.Server{
		Addr:         path,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      h.Route(),
	}

	log.Printf("Http server listening on %s", path)

	err = http.ListenAndServe()
	if err != nil {
		log.Printf("Failed to start http server: %v\n", err)
	}

	return http
}
