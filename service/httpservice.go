package service

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func attachmentUploadRequest(url string, fileType string, fileName string, filePath string,
	username string, secret string) []byte {
	file, fileErr := os.Open(filePath)
	if fileErr != nil {
		log.Fatalf("invalid file %s", fileErr)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, formErr := writer.CreateFormFile("file", filepath.Base(filePath))
	if formErr != nil {
		log.Fatalf("invalid form data %s", formErr)
	}
	_, formErr = io.Copy(part, file)
	_ = writer.WriteField("type", fileType)
	_ = writer.WriteField("name", fileName)
	formErr = writer.Close()
	if formErr != nil {
		log.Fatalf("invalid form data %s", formErr)
	}

	req, requestErr := http.NewRequest("POST", url, body)
	if requestErr != nil {
		log.Fatalf("invalid request %s", requestErr)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(username, secret)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, responseErr := client.Do(req)
	if responseErr != nil {
		log.Fatalf("invalid response %s", responseErr)
	}

	responseBody, responseBodyErr := ioutil.ReadAll(resp.Body)
	if responseBodyErr != nil {
		log.Fatalf("invalid responsebody %s", responseBodyErr)
	}
	return responseBody

}
