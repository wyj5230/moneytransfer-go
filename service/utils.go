package service

import (
	"os"
	"strconv"
	"time"
)

func GetExternalId() string {
	now := time.Now()
	sec := now.Unix()
	return strconv.FormatInt(sec, 10)
}

func GetApiKey() string {
	ApiKey := os.Getenv("MONEYTRANSFER_API_KEY")
	return ApiKey
}

func GetApiSecret() string {
	ApiSecret := os.Getenv("MONEYTRANSFER_API_SECRET")
	return ApiSecret
}

func FloatToString(float float32) string {
	return strconv.FormatFloat(float64(float), 'f', 0, 64)
}

func IntToString(int int) string {
	return strconv.Itoa(int)
}
