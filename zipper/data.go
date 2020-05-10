package zipper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Bucket  string   `json:"bucket"`
	ZipName string   `json:"zip_name"`
	Items   []string `json:"items"`
}

func ResponseParser(w http.ResponseWriter, r *http.Request) (Request, error) {
	var request Request
	buf, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(buf, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return request, err
	}
	return request, nil
}
