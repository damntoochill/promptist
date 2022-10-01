package image

import (
	"bytes"
	"image"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/rs/zerolog/log"
)

var awsConfig = &aws.Config{
	Credentials: credentials.NewStaticCredentials(os.Getenv("S3_KEY"), os.Getenv("S3_SECRET"), ""),
	Endpoint:    aws.String(os.Getenv("S3_ENDPOINT")),
	Region:      aws.String(os.Getenv("S3_REGION")),
}

var bucket = os.Getenv("S3_BUCKET")

// Upload an image to the cloud
func Upload(filename string, objectKey string, crop bool, sizes ...int) bool {

	img, err := imaging.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("failed to open image")
		return false
	}

	// Send original over
	buf := &bytes.Buffer{}
	err = imaging.Encode(buf, img, imaging.JPEG)
	if err != nil {
		log.Error().Err(err).Msg("failed to encode image")
		return false
	}

	ok := send(buf.Bytes(), objectKey)
	if !ok {
		return false
	}

	// Send resized
	for _, size := range sizes {
		var newImg *image.NRGBA
		sizeStr := strconv.Itoa(size)
		var objectKeyRenamed string
		if crop {
			newImg = imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)
			objectKeyRenamed = objectKey + "-" + sizeStr + "-c"
		} else {
			newImg = imaging.Resize(img, size, 0, imaging.Lanczos)
			objectKeyRenamed = objectKey + "-" + sizeStr

		}
		buf = &bytes.Buffer{}
		err = imaging.Encode(buf, newImg, imaging.JPEG)
		if err != nil {
			log.Error().Err(err).Msg("failed to encode image")
			return false
		}

		ok = send(buf.Bytes(), objectKeyRenamed)
		if !ok {
			return false
		}
	}
	return true
}

func send(b []byte, objectKey string) bool {

	// reader := bytes.NewReader(buf.Bytes())
	reader := bytes.NewReader(b)

	session := session.New(awsConfig)
	s3Client := s3.New(session)
	object := s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		Body:        reader,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/jpeg"),
		Metadata: map[string]*string{
			"x-amz-meta-my-key": aws.String("your-value"), //required
		},
	}
	_, err := s3Client.PutObject(&object)
	if err != nil {
		log.Error().Err(err).Msg("unable to upload to Spaces ")
		return false
	}
	return true
}
