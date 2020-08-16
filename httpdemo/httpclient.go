package httpdemo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type HttpClient struct {

}

func NewHttpClient() (*HttpClient, error) {
	return &HttpClient{}, nil
}

func (hc *HttpClient) UploadFile(filename string, targeturl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		log.Println("failed to CreateFormFile, err: ", err)
		return err
	}

	fh, err := os.Open(filename)
	if err != nil {
		log.Println("failed to Open file, err: ", err)
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Println("failed to io.Copy(), err: ", err)
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targeturl, contentType, bodyBuf)
	if err != nil {
		log.Println("failed to http.Post, err: ", err)
		return err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to ioutil.ReadAll(), err: ", err)
		return err
	}
	fmt.Println("resp.Status: ", resp.Status)
	fmt.Println("resp.Body: ", string(resBody))

	return nil
}

