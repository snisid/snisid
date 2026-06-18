package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type ImageRepo struct {
	client *minio.Client
	bucket string
}

func NewImageRepo(client *minio.Client, bucket string) *ImageRepo {
	return &ImageRepo{client: client, bucket: bucket}
}

func (r *ImageRepo) Upload(ctx context.Context, objectName string, data []byte, contentType string) error {
	_, err := r.client.PutObject(ctx, r.bucket, objectName, bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return fmt.Errorf("minio upload: %w", err)
	}
	return nil
}

func (r *ImageRepo) Download(ctx context.Context, objectName string) ([]byte, error) {
	obj, err := r.client.GetObject(ctx, r.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio get: %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("minio read: %w", err)
	}
	return data, nil
}

func (r *ImageRepo) Delete(ctx context.Context, objectName string) error {
	return r.client.RemoveObject(ctx, r.bucket, objectName, minio.RemoveObjectOptions{})
}
