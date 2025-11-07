package s3

import (
	"bytes"
	"context"
	"errors"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/shpdwx/web/conf"
)

func Put(ctx context.Context, m conf.Minio, image *ImageMeta) (*minio.UploadInfo, error) {

	if image.Name == "" || len(image.Byte) < 1 {
		return nil, errors.New("上传时图片信息异常")
	}

	var (
		endpoint  = m.Endpoint
		bucket    = m.Bucket
		ak        = m.AK
		sk        = m.SK
		secondDir = m.Dir
	)

	// minio client
	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ak, sk, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	if c.IsOffline() {
		return nil, errors.New("the endpoint is offline.")
	}

	// bucket info
	exist, err := c.BucketExists(context.Background(), bucket)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("the bucket do not exist.")
	}

	// image
	filename := path.Join(secondDir, image.Name)
	r := bytes.NewReader(image.Byte)

	// put image
	info, err := c.PutObject(ctx, bucket, filename, r, image.Size, minio.PutObjectOptions{
		ContentType:  image.ContentType,
		UserMetadata: map[string]string{"origin": image.Origin, "request-id": image.RequestId, "desc": image.Desc},
	})

	if err != nil {
		return nil, err
	}

	return &info, nil
}
