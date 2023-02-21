package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	uploader *s3manager.Uploader
)

func init() {
	godotenv.Load(".env")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	log.Printf("accessKey: %s", accessKey)
	s3Region := os.Getenv("S3_REGION")
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	config := &aws.Config{
		Credentials:      creds,
		Region:           aws.String(s3Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		MaxRetries:       aws.Int(3),
	}

	sess := session.Must(session.NewSession(config))
	uploader = s3manager.NewUploader(sess)
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))
	r.POST("/upload", saveFileHandler)
	log.Fatal(r.Run(":8080"))
}

func saveFileHandler(c *gin.Context) {
	godotenv.Load()
	s3Bucket := os.Getenv("S3_BUCKET")
	s3Url := os.Getenv("S3_URL")
	fileHeader, err := c.FormFile("file")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "fail to open file",
		})
		return
	}

	filename := fileHeader.Filename
	newFileName := uuid.New().String() + filename
	if err := putFileToS3(c.Request.Context(), s3Bucket, newFileName, f); err != nil {
		log.Printf("fail to upload to s3: %v", err)
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"message": "fail to upload to s3",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"url":     s3Url + newFileName,
		"message": "file uploaded: %s" + newFileName,
	})
}

func putFileToS3(ctx context.Context, bucket, fileName string, f io.Reader) error {
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   f,
	})
	if err != nil {
		return err
	}
	return nil
}
