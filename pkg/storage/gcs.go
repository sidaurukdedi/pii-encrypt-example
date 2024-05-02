package storage

import (
	"context"
	"io"
	"time"

	gcstorage "cloud.google.com/go/storage"
)

type gcsAdapter struct {
	gcpAccessID                   string
	gcpPrivateKey                 string
	client                        *gcstorage.Client
	defaultExpiresTimeOfSignedURL time.Duration
}

func NewGCSAdapter(client *gcstorage.Client, accessID, privateKey string) Storage {
	return &gcsAdapter{
		gcpAccessID:                   accessID,
		gcpPrivateKey:                 privateKey,
		client:                        client,
		defaultExpiresTimeOfSignedURL: time.Second * 60,
	}
}

func (gcs *gcsAdapter) PutObject(ctx context.Context, bucketName string, filepath string, file io.Reader, contentType string, meta map[string]string) (err error) {
	w := gcs.client.Bucket(bucketName).Object(filepath).NewWriter(ctx)
	w.ContentType = contentType
	w.Metadata = meta

	io.Copy(w, file)

	return w.Close()
}

func (gcs *gcsAdapter) SignURL(ctx context.Context, method, bucketName, filepath string, expiresIn time.Duration) (url string, err error) {
	url, err = gcs.client.Bucket(bucketName).SignedURL(filepath, &gcstorage.SignedURLOptions{
		GoogleAccessID: gcs.gcpAccessID,
		PrivateKey:     []byte(gcs.gcpPrivateKey),
		Expires:        time.Now().Add(expiresIn),
		Scheme:         gcstorage.SigningSchemeV4,
		Method:         method,
	})

	return
}
