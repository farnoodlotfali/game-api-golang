package objS3

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/game-api/utils"
	"github.com/joho/godotenv"
)

// S3 configuration (ParsPack)
var S3Client *s3.S3

func InitS3() {
	// Load .env in development; ignore error in prod if you use real env vars
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, relying on environment variables")
	}

	endpoint := os.Getenv("S3_ENDPOINT")
	region := os.Getenv("S3_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		panic(fmt.Errorf("failed to init S3 session: %v", err))
	}
	S3Client = s3.New(sess)
}

func UploadFileToS3(fileHeader *multipart.FileHeader, prefix string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
	bucket := os.Getenv("S3_BUCKET")

	randomString, err := utils.RandomString(10)

	if err != nil {
		return "", err
	}

	objectKey := prefix + "-" + randomString + "-" + fileHeader.Filename

	// Upload to S3
	_, err = S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		Body:        file,
		ACL:         aws.String("private"),
		ContentType: aws.String(DetectContentType(fileHeader.Filename)),
	})
	if err != nil {
		return "", fmt.Errorf("S3 PutObject failed: %w", err)
	}

	url := fmt.Sprintf("/%s", objectKey)
	return url, nil
}

func DetectContentType(path string) string {
	switch ext := filepath.Ext(path); ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

func GetS3Endpoint() string {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		return "http://localhost:9000" // Default for local MinIO
	}
	return endpoint
}
