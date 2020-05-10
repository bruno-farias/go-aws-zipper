package main

import (
	"github.com/bruno-farias/go-aws-zipper/config"
	"github.com/bruno-farias/go-aws-zipper/zipper"
	"log"
	"net/http"
)

func init() {
	config.SetEnvConfig()
}

func main() {
	http.HandleFunc("/", zipper.Download)
	log.Println("Running...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
