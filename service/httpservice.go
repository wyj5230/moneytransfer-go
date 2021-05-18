package service

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HttpGet(url string, username string, secret string) []byte {
	return httpRequest("GET", url, nil, username, secret)
}

func HttpPost(url string, reqJson []byte, username string, secret string) []byte {
	return httpRequest("POST", url, bytes.NewBuffer(reqJson), username, secret)
}

func httpRequest(method string, url string, requestBody io.Reader, username string, secret string) []byte {
	client := &http.Client{Timeout: 10 * time.Second}
	req, requestErr := http.NewRequest(method, url, requestBody)
	if requestErr != nil {
		log.Fatalf("invalid request %s", requestErr)
	}
	req.SetBasicAuth(username, secret)
	req.Header.Add("Content-Type", "application/json")
	resp, responseErr := client.Do(req)
	if responseErr != nil {
		log.Fatalf("invalid response %s", responseErr)
	}
	body, responseBodyErr := ioutil.ReadAll(resp.Body)
	if responseBodyErr != nil {
		log.Fatalf("invalid responsebody %s", responseBodyErr)
	}
	return body
}
