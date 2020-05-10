package zipper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"os"
	"time"
)

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Download(w http.ResponseWriter, r *http.Request) {
	request, _ := ResponseParser(w, r)
	dirname := "downloads/" + request.ZipName
	output := "zip/" + request.ZipName + ".zip"

	for _, item := range request.Items {
		_ = os.Mkdir(dirname, 0755)
		file, err := os.Create(dirname + "/" + item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			exitErrorf("Unable to open file %q, %v", item, err)
		}

		defer file.Close()

		sess, _ := session.NewSession()
		downloader := s3manager.NewDownloader(sess)

		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(request.Bucket),
				Key:    aws.String(item),
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			exitErrorf("Unable to download item %q, %v", item, err)
		}

		fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	}
	CreateZipFile(dirname, output)
	file, _ := os.Open(output)

	//remove zip
	defer os.Remove(output)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+request.ZipName+".zip")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")

	http.ServeContent(w, r, request.ZipName+".zip", time.Now(), file)
}
