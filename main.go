package main

import (
  "io"
  "log"
  "os"
  "os/exec"
  "bufio"
  // "compress/gzip"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
)


awsKey    := "Your Amazon s3 key"
awsSecret := "Your Amazon s3 Secret"

func main() {
  // Set credentials
  os.Setenv("AWS_ACCESS_KEY_ID", awsKey)
  os.Setenv("AWS_SECRET_ACCESS_KEY", awsSecret)

  // listAllBuckets()
  uploadFile()
}

// Upload a image to amazon S3
func uploadFile() {

  file, err := os.Open("./images/oso.jpeg")
  if err != nil {
    log.Fatal("Failed to open file :", err)
  }

  reader, writer := io.Pipe()
  go func () {
    content := bufio.NewWriter(writer)
    // gw := gzip.NewWriter(writer)
    // io.Copy(gw, file)
    io.Copy(content, file)

    file.Close()
    // gw.Close()
    //content.Close()
    writer.Close()
  }()

  uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
  key, err := exec.Command("uuidgen").Output()

  if err != nil {
    log.Fatal("Error generating UDID", err)
  }

  result, err := uploader.Upload(&s3manager.UploadInput{
    Body : reader,
    Bucket: aws.String("daemongear"),
    Key : aws.String(string(key)),
  })

  if err != nil {
    log.Fatalln("Failed to upload ", err)
  }
  log.Println("Success Upload to ", result.Location)
}



/* List all buckets*/
func listAllBuckets() {
  svc := s3.New(session.New(&aws.Config{Region: aws.String("us-west-2")}))
  result, err := svc.ListBuckets(&s3.ListBucketsInput{})
  if err != nil {
      log.Println("Failed to list buckets", err)
      return
  }
  log.Println("Buckets:")
  for _, bucket := range result.Buckets {
      log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
  }
}
