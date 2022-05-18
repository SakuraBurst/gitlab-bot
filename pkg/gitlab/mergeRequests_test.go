package gitlab

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// SINGLE MERGE REQUEST
func TestGetMRDiffs_BadURLClosedGoroutine(t *testing.T) {
	var isClosed = func() bool {
		return true
	}
	resChan := make(chan MergeRequestTransfer, 1)
	go gitLabBadURLMock.getMRWithDiffs(0, resChan, isClosed)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, resChan, 0)
}

func TestGetMRDiffs_BadRequestClosedGoroutine(t *testing.T) {
	var isClosed = func() bool {
		return true
	}
	resChan := make(chan MergeRequestTransfer, 1)
	clients.EnableMock()
	mockErr := errors.New("test error")
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: nil,
		Err:      mockErr,
	})
	go gitLabMock.getMRWithDiffs(0, resChan, isClosed)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, resChan, 0)
}

func TestGetMRDiffs_BadBodyClosedGoroutine(t *testing.T) {
	var isClosed = func() bool {
		return true
	}
	resChan := make(chan MergeRequestTransfer, 1)
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
			Body:       invalidBody,
		},
		Err: nil,
	})
	go gitLabMock.getMRWithDiffs(0, resChan, isClosed)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, resChan, 0)
}

func TestGetMRDiffs_URLError(t *testing.T) {
	resChan := make(chan MergeRequestTransfer)
	go gitLabBadURLMock.getMRWithDiffs(0, resChan, falseFuncMock)
	res := <-resChan
	assert.Nil(t, res.mergeRequest)
	assert.NotNil(t, res.error)
	assert.EqualError(t, res.error, invalidMergeRequestUrlErrorString)
}

func TestGetMRDiffs_RequestError(t *testing.T) {
	clients.EnableMock()
	mockErr := errors.New("test error")
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: nil,
		Err:      mockErr,
	})
	resChan := make(chan MergeRequestTransfer)
	go gitLabMock.getMRWithDiffs(0, resChan, falseFuncMock)
	res := <-resChan
	assert.Nil(t, res.mergeRequest)
	assert.NotNil(t, res.error)
	assert.EqualError(t, res.error, mockErr.Error())
	clients.Mocks.ClearMocks()
}

func TestGetMRDiffs_ErrorWithBodyError(t *testing.T) {
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
			Body:       invalidBody,
		},
		Err: nil,
	})
	resChan := make(chan MergeRequestTransfer)
	go gitLabMock.getMRWithDiffs(0, resChan, falseFuncMock)
	res := <-resChan
	assert.Nil(t, res.mergeRequest)
	assert.NotNil(t, res.error)
	assert.EqualError(t, res.error, invalidArgumentErrorString)
	clients.Mocks.ClearMocks()
}

func TestGetMRDiffs_ErrorWithGitlabError(t *testing.T) {
	clients.EnableMock()
	gle := models.GitlabError{Message: "Unauthorized"}
	gleBytes, err := json.Marshal(gle)
	require.Nil(t, err)
	errorBody := bytes.NewReader(gleBytes)
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(errorBody),
		},
		Err: nil,
	})
	resChan := make(chan MergeRequestTransfer)
	go gitLabMock.getMRWithDiffs(0, resChan, falseFuncMock)
	res := <-resChan
	assert.Nil(t, res.mergeRequest)
	assert.NotNil(t, res.error)
	assert.EqualError(t, res.error, gle.Error())
	clients.Mocks.ClearMocks()
}

func TestGetMRDiffs_OKWithBodyError(t *testing.T) {
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       invalidBody,
		},
		Err: nil,
	})
	resChan := make(chan MergeRequestTransfer)
	go gitLabMock.getMRWithDiffs(0, resChan, falseFuncMock)
	res := <-resChan
	assert.Nil(t, res.mergeRequest)
	assert.NotNil(t, res.error)
	assert.EqualError(t, res.error, invalidArgumentErrorString)
	clients.Mocks.ClearMocks()
}

func TestGetMRDiffs(t *testing.T) {
	clients.EnableMock()
	mergeRequest := models.MergeRequest{}
	mergeRequestBytes, err := json.Marshal(mergeRequest)
	require.Nil(t, err)
	validResponseBody := bytes.NewReader(mergeRequestBytes)
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(validResponseBody),
		},
		Err: nil,
	})
	resChan := make(chan MergeRequestTransfer)
	go gitLabMock.getMRWithDiffs(0, resChan, falseFuncMock)
	res := <-resChan
	assert.NotNil(t, res.mergeRequest)
	assert.Nil(t, res.error)
	assert.Equal(t, res.mergeRequest, &mergeRequest)
	clients.Mocks.ClearMocks()
}

// ALL MERGE REQUESTS

func TestGetAllOpenedMergeRequests_URLError(t *testing.T) {
	mri, err := gitLabBadURLMock.getAllOpenedMergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidMergeRequestsUrlErrorString)
}

func TestGetAllOpenedMergeRequests_RequestError(t *testing.T) {
	clients.EnableMock()
	mockErr := errors.New("test error")
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: nil,
		Err:      mockErr,
	})
	mri, err := gitLabMock.getAllOpenedMergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.EqualError(t, err, mockErr.Error())
	clients.Mocks.ClearMocks()
}

func TestGetAllOpenedMergeRequests_ErrorWithBodyError(t *testing.T) {
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
			Body:       invalidBody,
		},
		Err: nil,
	})
	mri, err := gitLabMock.getAllOpenedMergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidArgumentErrorString)
	clients.Mocks.ClearMocks()
}

func TestGetAllOpenedMergeRequests_ErrorWithGitlabError(t *testing.T) {
	clients.EnableMock()
	gle := models.GitlabError{Message: "Unauthorized"}
	gleBytes, err := json.Marshal(gle)
	require.Nil(t, err)
	errorBody := bytes.NewReader(gleBytes)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(errorBody),
		},
		Err: nil,
	})
	mri, err := gitLabMock.getAllOpenedMergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.EqualError(t, err, gle.Error())
	clients.Mocks.ClearMocks()
}

func TestGetAllOpenedMergeRequests_OKWithBodyError(t *testing.T) {
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       invalidBody,
		},
		Err: nil,
	})
	mri, err := gitLabMock.getAllOpenedMergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidArgumentErrorString)
	clients.Mocks.ClearMocks()
}

func TestGetAllOpenedMergeRequests(t *testing.T) {
	clients.EnableMock()
	glMergeRequests := []models.MergeRequest{{}, {}}
	glMergeRequestsBytes, err := json.Marshal(glMergeRequests)
	require.Nil(t, err)
	validResponseBody := bytes.NewReader(glMergeRequestsBytes)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(validResponseBody),
		},
		Err: nil,
	})
	mri, err := gitLabMock.getAllOpenedMergeRequests()
	assert.NotNil(t, mri)
	assert.Nil(t, err)
	assert.Equal(t, mri.Length, 2)
	assert.Equal(t, mri.MergeRequests, glMergeRequests)
	clients.Mocks.ClearMocks()
}
