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
)

var tokenMock = "test"

var gitLabBadURLMock = Gitlab{
	url:   "_394e2904dfsfs234)Ue)R(UEWR#$%@%",
	repo:  "test",
	token: tokenMock,
}

var gitLabMock = Gitlab{
	url:   "https://test.com",
	repo:  "test/test",
	token: tokenMock,
}

var gitLabMockWithDiffs = Gitlab{
	url:       "https://test.com",
	repo:      "test/test",
	token:     tokenMock,
	WithDiffs: true,
}

var mergeRequestsURLMock = "https://test.com/api/v4/projects/test%2Ftest/merge_requests?state=opened&with_merge_status_recheck=true"

var mergeRequestURLMock = "https://test.com/api/v4/projects/test%2Ftest/merge_requests/0/changes"

func TestGetMRDiffs_RequestError(t *testing.T) {
	panic("implement me")
}

func TestGetMrsWithDiffs_EmptyMergeRequestsInfo(t *testing.T) {
	panic("implement me")
}

func TestMergeRequestsRequest_URLError(t *testing.T) {
	mri, err := gitLabBadURLMock.MergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
}

func TestMergeRequests_RequestError(t *testing.T) {
	clients.EnableMock()
	mockErr := errors.New("test error")
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: nil,
		Err:      mockErr,
	})
	mri, err := gitLabMock.MergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.Equal(t, err, mockErr)
	clients.Mocks.ClearMocks()
}

func TestMergeRequests_ErrorWithBodyError(t *testing.T) {
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     "Invalid body test",
			StatusCode: http.StatusInternalServerError,
			Body:       invalidBody,
		},
		Err: nil,
	})
	mri, err := gitLabMock.MergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	clients.Mocks.ClearMocks()
}

func TestMergeRequests_ErrorWithGitlabError(t *testing.T) {
	clients.EnableMock()
	gle := models.GitlabError{Message: "Unauthorized"}
	gleBytes, err := json.Marshal(gle)
	require.Nil(t, err)
	errorBody := bytes.NewReader(gleBytes)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     "401 Unauthorized",
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(errorBody),
		},
		Err: nil,
	})
	mri, err := gitLabMock.MergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.Equal(t, err, gle)
	clients.Mocks.ClearMocks()
}

func TestMergeRequests_OKWithBodyError(t *testing.T) {
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     "Invalid body test",
			StatusCode: http.StatusOK,
			Body:       invalidBody,
		},
		Err: nil,
	})
	mri, err := gitLabMock.MergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	clients.Mocks.ClearMocks()
}

func TestMergeRequests(t *testing.T) {
	clients.EnableMock()
	glMergeRequests := []models.MergeRequest{{}, {}}
	glMergeRequestsBytes, err := json.Marshal(glMergeRequests)
	require.Nil(t, err)
	mrBody := bytes.NewReader(glMergeRequestsBytes)
	clients.Mocks.AddMock(mergeRequestsURLMock, clients.Mock{
		Response: &http.Response{
			Status:     "Invalid body test",
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(mrBody),
		},
		Err: nil,
	})
	mri, err := gitLabMock.MergeRequests()
	assert.NotNil(t, mri)
	assert.Nil(t, err)
	assert.Equal(t, mri.Length, 2)
	assert.Equal(t, mri.MergeRequests, glMergeRequests)
	clients.Mocks.ClearMocks()
}
