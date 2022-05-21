package telegram

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type SendDocumentRequestHelper struct {
	MainBuffer *bytes.Buffer
	Writer     *multipart.Writer
	File       *os.File
}

func (t Bot) createSendFileURL() (*url.URL, error) {
	rawUrl := fmt.Sprintf("%s/bot%s/sendDocument", telegramApi, t.token)
	sendMessageURL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	return sendMessageURL, nil
}

func (t Bot) makeBodyAndHeadersToSendFile(file *os.File, fileName string) (io.Reader, http.Header, error) {
	sendFileRequestHelper := prepareSendMessageRequestHelper(file)
	err := sendFileRequestHelper.Writer.WriteField("chat_id", t.mainChannel)
	if err != nil {
		return nil, nil, err
	}
	body, err := sendFileRequestHelper.Writer.CreateFormFile("document", fileName)
	if err != nil {
		return nil, nil, err
	}
	_, err = io.Copy(body, sendFileRequestHelper.File)
	if err != nil {
		return nil, nil, err
	}
	err = sendFileRequestHelper.Writer.Close()
	if err != nil {
		return nil, nil, err
	}
	headers := http.Header{}
	headers.Set("Content-Type", sendFileRequestHelper.Writer.FormDataContentType())
	return sendFileRequestHelper.MainBuffer, headers, nil
}

func prepareSendMessageRequestHelper(file *os.File) *SendDocumentRequestHelper {
	buffer := bytes.NewBuffer(nil)
	multipartFile := multipart.NewWriter(buffer)
	return &SendDocumentRequestHelper{
		MainBuffer: buffer,
		Writer:     multipartFile,
		File:       file,
	}
}
