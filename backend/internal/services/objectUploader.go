package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type ObjectUploader struct {
	s3Client          *s3.Client
	bucketName        string
	bucketRegion      string
	errObjectTooLarge error
	errObjectNotFound error
}

// Creates a new Object Uploader with default config using credentials from .env
// and creates a named bucket from .env if it does not exist
func NewObjectUploader() (*ObjectUploader, error) {
	dc, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	ou := ObjectUploader{
		s3Client: s3.NewFromConfig(
			aws.Config{
				Region:       *aws.String(os.Getenv("AWS_BUCKET_REGION")),
				BaseEndpoint: aws.String(os.Getenv("AWS_BUCKET_ENDPOINT")),
				Credentials:  dc.Credentials,
			},
		),
		bucketName:        os.Getenv("AWS_BUCKET_NAME"),
		bucketRegion:      os.Getenv("AWS_BUCKET_REGION"),
		errObjectTooLarge: errors.New("The object for upload is too large!"),
		errObjectNotFound: errors.New("Specified object does not exist or was not created"),
	}

	exists, err := ou.checkExists(ou.bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = ou.createBucket()
		if err != nil {
			return nil, err
		}
	}

	return &ou, nil
}

func (u *ObjectUploader) checkExists(bucketName string) (bool, error) {
	var exists bool = true
	_, err := u.s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				exists = false
				err = nil
			default:
				return false, errors.New(
					fmt.Sprintf("Either you don't have access to bucket %v or another error occurred. "+
						"Here's what happened: %v\n", bucketName, err))
			}
		}
	}
	return exists, err
}

func (u *ObjectUploader) createBucket() error {
	_, err := u.s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(u.bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(u.bucketRegion),
		},
	})
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			return errors.New(fmt.Sprintf("You already own bucket %s.\n", u.bucketName))
		} else if errors.As(err, &exists) {
			return errors.New(fmt.Sprintf("Bucket %s already exists.\n", u.bucketName))
		}
	} else {
		err = s3.NewBucketExistsWaiter(u.s3Client).Wait(
			context.Background(), &s3.HeadBucketInput{Bucket: aws.String(u.bucketName)}, time.Minute)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *ObjectUploader) PutObject(objectKey string, file io.Reader) (string, error) {
	var link string

	_, err := u.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(u.bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	err = s3.NewObjectExistsWaiter(u.s3Client).Wait(
		context.Background(), &s3.HeadObjectInput{Bucket: aws.String(u.bucketName), Key: aws.String(objectKey)}, time.Minute)
	if err != nil {
		return "", u.errObjectNotFound
	}
	// TODO: Get a link for the object
	return link, nil
}

// Checks the provided file size compared to the max quota for the user tier
func (u *ObjectUploader) CheckFileSize(size int64, userTier string) error {
	var sizeCap int64

	switch userTier {
	case "freetier":
		sizeCap = 50 * 1024 * 1024 // 50M
	case "paidtier":
		sizeCap = 150 * 1024 * 1024 // 150M
	case "picture":
		sizeCap = 12 * 1024 * 1024 // 12M
	}

	if size > sizeCap {
		return u.errObjectTooLarge
	}
	return nil
}

func (u *ObjectUploader) DeleteObject(objectKey string) error {
	_, err := u.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &u.bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		var noKey *types.NoSuchKey
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &noKey) {
			return u.errObjectNotFound
		} else if errors.As(err, &apiErr) {
			return err
		}
	}
	err = s3.NewObjectNotExistsWaiter(u.s3Client).Wait(
		context.Background(), &s3.HeadObjectInput{Bucket: aws.String(u.bucketName), Key: aws.String(objectKey)}, time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func (u *ObjectUploader) ObjectTooLargeErr() error { return u.errObjectTooLarge }
func (u *ObjectUploader) ObjectNotFoundErr() error { return u.errObjectNotFound }
