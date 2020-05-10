package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bruno-farias/go-aws-zipper/zipper"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	AwsAccessKeyId, _ := os.LookupEnv("AWS_ACCESS_KEY_ID")
	AwsSecretAccessKey, _ := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	AwsRegion, _ := os.LookupEnv("AWS_REGION")
	_ = os.Setenv("AWS_ACCESS_KEY_ID", AwsAccessKeyId)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", AwsSecretAccessKey)
	_ = os.Setenv("AWS_REGION", AwsRegion)
}

var requestJson map[string]interface{}
var dirname string
var zipName string
var output string

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func responseParser(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	buf, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(buf, &requestJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return requestJson, err
}

func download(w http.ResponseWriter, r *http.Request) {
	requestJson, _ := responseParser(w, r)
	bucket := requestJson["bucket"].(string)
	zipName = requestJson["zip_name"].(string)
	items := requestJson["items"]
	dirname = "downloads/" + zipName

	for _, item := range items.([]interface{}) {
		_ = os.Mkdir(dirname, 0755)
		file, err := os.Create(dirname + "/" + item.(string))
		if err != nil {
			exitErrorf("Unable to open file %q, %v", item, err)
		}

		defer file.Close()

		sess, _ := session.NewSession()
		downloader := s3manager.NewDownloader(sess)

		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(item.(string)),
			})
		if err != nil {
			exitErrorf("Unable to download item %q, %v", item, err)
		}

		fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	}
	createZipFile()
	file, _ := os.Open(output)
	//remove zip
	defer os.Remove(output)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=" + zipName + ".zip")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")

	http.ServeContent(w, r, zipName + ".zip", time.Now(), file)
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func createZipFile() {
	zipName := requestJson["zip_name"].(string)
	output = "zip/" + zipName + ".zip"
	var files []string

	err := filepath.Walk(dirname, visit(&files))
	if err != nil {
		panic(err)
	}
	// remove dir and files
	defer os.RemoveAll(dirname)

	// removes folder from slice
	files = files[1:]

	if err := zipper.ZipFiles(output, files); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", download)
	log.Println("Running...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
