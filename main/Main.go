package main

import (
	"github.com/gin-gonic/gin"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/apex/gateway"
	"log"
)

const BucketName = "audio-streaming-s3"
const Name = "name"
const FileName = "fileName"


func main() {
	addr := ":" + os.Getenv("PORT")
	log.Fatal(gateway.ListenAndServe(addr, routerEngine()))
}

func routerEngine() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("audio-streaming/ping", pingHandler)
	r.POST("audio-streaming/upload", uploadHandler)
	return r
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func uploadHandler(c *gin.Context) {
	isSuccess, code := upload(c)
	if isSuccess {
		c.JSON(code, gin.H{
			"message": "upload",
		})
	} else {
		c.JSON(code, gin.H{
			"message": "error",
		})
	}
}

func upload(c *gin.Context) (success bool, code int) {

	formFile, err := c.FormFile(Name)
	fileName, exist := c.GetPostForm(FileName)
	if !exist {
		log.Println("Erreur de récupération du fileName")
		return false, 401
	}
	if err != nil {
		log.Println("Erreur durant la récupération du fichier")
		log.Println(err.Error())
		return false, 500
	}

	file, err := formFile.Open()
	if err != nil {
		log.Println("Erreur durant l'ouverture du fichier")
		log.Println(err.Error())
		return false, 500
	}

	log.Println("Début de l'upload")


	// Create an uploader with the session and default options
	uploader := getAwsUploader()

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		log.Println("Erreur 2 ")
		log.Println(err.Error())
		return
	}
	log.Println("File uploaded to "+ result.Location)
	return true, 200
}

func getAwsConfig() (*aws.Config) {
	return aws.NewConfig().WithRegion("eu-west-3").WithCredentialsChainVerboseErrors(true).WithCredentials(credentials.AnonymousCredentials)
}

func getAwsSession() (*session.Session) {
	return session.Must(session.NewSession(getAwsConfig()))
}

func getAwsUploader() (*s3manager.Uploader) {
	return s3manager.NewUploader(getAwsSession())
}