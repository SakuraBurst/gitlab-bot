package gitlab

import (
	"bytes"
	"encoding/json"
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
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

var mergeRequestURLMock1 = "https://test.com/api/v4/projects/test%2Ftest/merge_requests/1/changes"

var invalidMergeRequestsUrlErrorString = "parse \"_394e2904dfsfs234)Ue)R(UEWR#$%@%/api/v4/projects/test/merge_requests\": invalid URL escape \"%@%\""

var invalidMergeRequestUrlErrorString = "parse \"_394e2904dfsfs234)Ue)R(UEWR#$%@%/api/v4/projects/test/merge_requests/0/changes\": invalid URL escape \"%@%\""

var invalidArgumentErrorString = "invalid argument"

func falseFuncMock() bool {
	return false
}

func TestNewGitlabConn(t *testing.T) {
	testString := "test"
	glConn := NewGitlabConn(true, testString, testString, testString)
	assert.NotNil(t, glConn)
	assert.True(t, glConn.WithDiffs)
	assert.Equal(t, testString, glConn.repo)
	assert.Equal(t, testString, glConn.token)
	assert.Equal(t, testString, glConn.url)

}

func TestGetAllOpenedMrsWithDiffs_ZeroMergeRequests(t *testing.T) {
	mri := MergeRequestsInfo{
		MergeRequests: make([]models.MergeRequest, 0),
	}
	mrsWithDiffs, err := getAllOpenedMrsWithDiffs(gitLabMock, &mri)
	assert.Nil(t, err)
	assert.Equal(t, &mri, mrsWithDiffs)
}

func TestGetAllOpenedMrsWithDiffs_Error(t *testing.T) {
	clients.EnableMock()
	mergeRequest := models.MergeRequest{}
	mergeRequestBytes, err := json.Marshal(mergeRequest)
	require.Nil(t, err)
	validResponseMrBody := bytes.NewReader(mergeRequestBytes)
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(validResponseMrBody),
		},
		Err: nil,
	})
	gle := models.GitlabError{Message: "Unauthorized"}
	gleBytes, err := json.Marshal(gle)
	require.Nil(t, err)
	errorBody := bytes.NewReader(gleBytes)
	clients.Mocks.AddMock(mergeRequestURLMock1, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(errorBody),
		},
		Err: nil,
	})
	mri := MergeRequestsInfo{
		Length:        2,
		MergeRequests: []models.MergeRequest{{}, {Iid: 1}},
	}
	mrsWithDiffs, err := getAllOpenedMrsWithDiffs(gitLabMock, &mri)
	assert.Nil(t, mrsWithDiffs)
	assert.NotNil(t, err)
	assert.EqualError(t, err, gle.Error())
	clients.DisableMock()
}

func TestGetAllOpenedMrsWithDiffs(t *testing.T) {
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
	mri := MergeRequestsInfo{
		Length:        1,
		MergeRequests: []models.MergeRequest{mergeRequest},
	}
	mrsWithDiffs, err := getAllOpenedMrsWithDiffs(gitLabMock, &mri)
	assert.NotNil(t, mrsWithDiffs)
	assert.Nil(t, err)
	assert.Equal(t, &mri, mrsWithDiffs)
	clients.DisableMock()
}

func TestGitlab_MergeRequests_Error(t *testing.T) {
	mri, err := gitLabBadURLMock.MergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidMergeRequestsUrlErrorString)
}

func TestGitlab_MergeRequests_OKWithoutDiffs(t *testing.T) {
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
	mri, err := gitLabMock.MergeRequests()
	assert.NotNil(t, mri)
	assert.Nil(t, err)
	assert.Equal(t, 2, mri.Length)
	assert.Equal(t, glMergeRequests, mri.MergeRequests)
	clients.DisableMock()
}

// наверное это уже функциональные тесты, но мне пофиг

func TestGitlab_MergeRequests_ErrorWithDiffs(t *testing.T) {
	clients.EnableMock()
	glMergeRequestsBytes, err := json.Marshal([]models.MergeRequest{{}})
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
	mri, err := gitLabMockWithDiffs.MergeRequests()
	assert.Nil(t, mri)
	assert.NotNil(t, err)
	assert.EqualError(t, err, gle.Error())
	clients.DisableMock()
}

func TestGitlab_MergeRequests_OKWithDiffs(t *testing.T) {
	clients.EnableMock()
	glMergeRequests := []models.MergeRequest{{}}
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
	mergeRequest := models.MergeRequest{}
	mergeRequestBytes, err := json.Marshal(mergeRequest)
	require.Nil(t, err)
	validResponseMrBody := bytes.NewReader(mergeRequestBytes)
	clients.Mocks.AddMock(mergeRequestURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(validResponseMrBody),
		},
		Err: nil,
	})
	mrsWithDiffs, err := gitLabMockWithDiffs.MergeRequests()
	assert.NotNil(t, mrsWithDiffs)
	assert.Nil(t, err)
	assert.NotNil(t, mrsWithDiffs.MergeRequests)
	assert.Equal(t, mergeRequest, mrsWithDiffs.MergeRequests[0])
	clients.DisableMock()
}
