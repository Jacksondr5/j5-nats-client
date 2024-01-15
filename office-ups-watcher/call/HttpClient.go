package call

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/jacksondr5/go-monorepo/logger"
	log "github.com/sirupsen/logrus"
)

type HttpClientImpl struct {}

type logFields struct {
	bodyData string
	requestName string
	responseBody string
	responseCode int
	url string
}

func getLogFields(fields logFields) log.Fields {
	return log.Fields{
		"bodyData": fields.bodyData,
		"requestName": fields.requestName,
		"responseBody": fields.responseBody,
		"responseCode": fields.responseCode,
		"url": fields.url,
	}
}

func (h HttpClientImpl) Post(url string, bodyData string, requestName string, accessToken string) error {
	logFields := logFields{
		bodyData: bodyData,
		requestName: requestName,
		url: url,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(bodyData))
	if err != nil {
		logger.ErrorWithFields("Error creating HTTP request", err, getLogFields(logFields))
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	
	response, err := client.Do(request)
	logFields.responseCode = response.StatusCode
	if err != nil {
		logger.ErrorWithFields("Error making HTTP request", err, getLogFields(logFields))
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		logger.ErrorWithFields("Error reading HTTP response body", err, getLogFields(logFields))
	}
	responseBody := string(b)
	logFields.responseBody = responseBody

	if response.StatusCode != 200 {
		logger.InfoWithFields("Non-200 code returned when calling %s", getLogFields(logFields))
	}
	logger.DebugWithFields("HTTP request completed successfully", getLogFields(logFields))
	return err
}