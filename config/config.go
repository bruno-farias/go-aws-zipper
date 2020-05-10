package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func SetEnvConfig() {
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
