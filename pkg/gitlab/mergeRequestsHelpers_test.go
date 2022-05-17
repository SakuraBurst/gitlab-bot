package gitlab

import (
	"bytes"
	"encoding/json"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, OPENED, "opened")
}

func TestGetMergeRequestURL_URL_Error(t *testing.T) {
	url, header, err := gitLabBadURLMock.getMergeRequestURL()
	assert.Nil(t, url)
	assert.Nil(t, header)
	assert.NotNil(t, err)
}

func TestGetMergeRequestURL_OK(t *testing.T) {
	url, headers, err := gitLabMock.getMergeRequestURL()
	require.NotNil(t, url)
	require.NotNil(t, headers)
	assert.Equal(t, url.String(), mergeRequestsURLMock)
	assert.Contains(t, headers, "Private-Token")
	assert.Equal(t, headers.Get("PRIVATE-TOKEN"), tokenMock)
	assert.Nil(t, err)
}

func TestDecodeMergeRequestInfoNilRequest(t *testing.T) {
	mri, err := decodeMergeRequestsInfo(nil)
	assert.Nil(t, mri)
	require.NotNil(t, err)
	assert.Equal(t, err.Error(), "request is nil")
}

func TestDecodeMergeRequestInfoErr_BodyError(t *testing.T) {
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	response := &http.Response{
		Status:     "Invalid body test",
		StatusCode: http.StatusInternalServerError,
		Body:       invalidBody,
	}
	mri, err := decodeMergeRequestsInfo(response)
	assert.Nil(t, mri)
	assert.NotNil(t, err)
}

func TestDecodeMergeRequestInfoErr_GitlabError(t *testing.T) {
	gle := models.GitlabError{Message: "Unauthorized"}
	gleBytes, err := json.Marshal(gle)
	require.Nil(t, err)
	errorBody := bytes.NewReader(gleBytes)
	response := &http.Response{
		Status:     "401 Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(errorBody),
	}
	mri, err := decodeMergeRequestsInfo(response)
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.Equal(t, gle, err)
}

func TestDecodeMergeRequestInfoOk_BodyError(t *testing.T) {
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	response := &http.Response{
		Status:     "Invalid body test",
		StatusCode: http.StatusOK,
		Body:       invalidBody,
	}
	mri, err := decodeMergeRequestsInfo(response)
	assert.Nil(t, mri)
	assert.NotNil(t, err)
}

func TestDecodeMergeRequestInfo(t *testing.T) {
	glMergeRequests := []models.MergeRequest{{}, {}}
	glMergeRequestsBytes, err := json.Marshal(glMergeRequests)
	require.Nil(t, err)
	errorBody := bytes.NewReader(glMergeRequestsBytes)
	response := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(errorBody),
	}
	mri, err := decodeMergeRequestsInfo(response)
	assert.NotNil(t, mri)
	assert.Nil(t, err)
	assert.Equal(t, mri.Length, 2)
	assert.Equal(t, mri.MergeRequests, glMergeRequests)
}

func TestGetSingleMergeRequestWithChangesURL_URL_Error(t *testing.T) {
	url, header, err := gitLabBadURLMock.getSingleMergeRequestWithChangesURL(0)
	assert.Nil(t, url)
	assert.Nil(t, header)
	assert.NotNil(t, err)
}

func TestGetSingleMergeRequestWithChangesURL_OK(t *testing.T) {
	tokenMock := "test"
	url, headers, err := gitLabMock.getSingleMergeRequestWithChangesURL(0)
	require.NotNil(t, url)
	require.NotNil(t, headers)
	assert.Equal(t, url.String(), mergeRequestURLMock)
	assert.Contains(t, headers, "Private-Token")
	assert.Equal(t, headers.Get("PRIVATE-TOKEN"), tokenMock)
	assert.Nil(t, err)
}

func TestDecodeSingleMergeRequestItemErr_BodyError(t *testing.T) {
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	response := &http.Response{
		Status:     "Invalid body test",
		StatusCode: http.StatusInternalServerError,
		Body:       invalidBody,
	}
	mri, err := decodeSingleMergeRequestItem(response)
	assert.Nil(t, mri)
	assert.NotNil(t, err)
}

func TestDecodeSingleMergeRequestItemErr_GitlabError(t *testing.T) {
	gle := models.GitlabError{Message: "Unauthorized"}
	gleBytes, err := json.Marshal(gle)
	require.Nil(t, err)
	errorBody := bytes.NewReader(gleBytes)
	response := &http.Response{
		Status:     "401 Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(errorBody),
	}
	mri, err := decodeSingleMergeRequestItem(response)
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.Equal(t, gle, err)
}

func TestDecodeSingleMergeRequestItemOk_BodyError(t *testing.T) {
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	response := &http.Response{
		Status:     "Invalid body test",
		StatusCode: http.StatusOK,
		Body:       invalidBody,
	}
	mri, err := decodeSingleMergeRequestItem(response)
	assert.Nil(t, mri)
	assert.NotNil(t, err)
}

func TestDecodeSingleMergeRequestItem(t *testing.T) {
	glMergeRequest := &models.MergeRequest{}
	glMergeRequestsBytes, err := json.Marshal(glMergeRequest)
	require.Nil(t, err)
	mrBody := bytes.NewReader(glMergeRequestsBytes)
	response := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(mrBody),
	}
	mr, err := decodeSingleMergeRequestItem(response)
	assert.NotNil(t, mr)
	assert.Nil(t, err)
	assert.Equal(t, mr, glMergeRequest)
}
