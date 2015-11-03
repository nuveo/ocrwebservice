package ocrws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	ocrURL = "http://www.ocrwebservice.com/restservices/processDocument?gettext=true"
)

// Config struct with LicenseCode and UserName
type Config struct {
	LicenseCode string `json:"license_code"`
	UserName    string `json:"username"`
}

// Response is struct representing response of ocrwebservice.
type Response struct {
	ErrorMessage    string     `json:"ErrorMessage"`
	AvailablePages  int        `json:"AvailablePages"`
	OCRText         [][]string `json:"OCRText"`
	OutputFileURL   string     `json:"OutputFileUrl"`
	TaskDescription string     `json:"TaskDescription"`
	Reserved        []string   `json:"Reserved"`
}

// Setup get LICENSE CODE and Username vars of environment.
func (c *Config) Setup() error {
	license := os.Getenv("LICENSE_CODE")
	username := os.Getenv("USERNAME")

	if license != "" && username != "" {
		c.LicenseCode = license
		c.UserName = username
		return nil
	}
	return errors.New("Export LICENSE_CODE and USERNAME environ vars")
}

func newfileUploadRequest(uri, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		log.Println("Creating POST Request")
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request, err
}

// OcrWs its main function of package
func OcrWs(pathFile string, language string) (Response, error) {
	conf := Config{}
	err := conf.Setup()
	if err != nil {
		log.Println("Make setup function")
		return Response{}, err
	}

	url := fmt.Sprintf("%s&language=%s&outputformat=txt", ocrURL, language)
	req, err := newfileUploadRequest(url, pathFile)
	req.SetBasicAuth(conf.UserName, conf.LicenseCode)

	if err != nil {
		return Response{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Fetch Result")
		return Response{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return Response{}, err
	}
	defer resp.Body.Close()

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
		return Response{}, err
	}

	if resp.StatusCode != 200 {
		return Response{}, errors.New(r.ErrorMessage)
	}

	return r, err
}

// Text method return Text Recognized coming Response
func (r *Response) Text() string {
	var text []string

	for indZone := range r.OCRText {
		for _, page := range r.OCRText[indZone] {
			text = append(text, page)
		}
	}
	return strings.Join(text, " ")
}
