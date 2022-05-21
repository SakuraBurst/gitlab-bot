package templates

import (
	"bytes"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/SakuraBurst/gitlab-bot/pkg/services/gitlab"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, "can_be_merged", CanBeMerged)
	assert.Equal(t, "02.01.2006 15:04", HumanTimeFormat)
}

func TestGetRightTemplate_TelegramMessageTemplateNewMrWithDiffs(t *testing.T) {
	template := GetRightTemplate(true, true)
	assert.Equal(t, "TelegramMessageTemplateNewMrWithDiffs", template.Name())
}

func TestGetRightTemplate_TelegramMessageTemplateNewMrWithoutDiffs(t *testing.T) {
	template := GetRightTemplate(true, false)
	assert.Equal(t, "TelegramMessageTemplateNewMrWithoutDiffs", template.Name())
}

func TestGetRightTemplate_TelegramMessageTemplateWithDiffs(t *testing.T) {
	template := GetRightTemplate(false, true)
	assert.Equal(t, "TelegramMessageTemplateWithDiffs", template.Name())
}

func TestGetRightTemplate_TelegramMessageTemplateWithoutDiffs(t *testing.T) {
	template := GetRightTemplate(false, false)
	assert.Equal(t, "TelegramMessageTemplateWithoutDiffs", template.Name())
}

func TestNewMrTitleZeroLength(t *testing.T) {
	title := newMrTitle(nil)
	assert.Equal(t, "Новый MR", title)
}

func TestNewMrTitleLengthOne(t *testing.T) {
	title := newMrTitle([]models.MergeRequest{{}})
	assert.Equal(t, "Новый MR", title)
}

func TestNewMrTitleLengthTwo(t *testing.T) {
	title := newMrTitle([]models.MergeRequest{{}, {}})
	assert.Equal(t, "Новые MR'ы", title)
}

func TestLastUpdateUpdatedAfterCreated(t *testing.T) {
	created := time.Now()
	updated := time.Now().Add(time.Hour)
	lastUpdate := lastUpdate(created, updated)
	assert.Greater(t, lastUpdate.Unix(), updated.Unix())
	assert.Equal(t, updated.Add(time.Hour*3).Unix(), lastUpdate.Unix())
	assert.Greater(t, lastUpdate.Unix(), created.Add(time.Hour*3).Unix())
}

func TestLastUpdateCreatedEqualsUpdated(t *testing.T) {
	timeNow := time.Now()
	lastUpdate := lastUpdate(timeNow, timeNow)
	assert.Greater(t, lastUpdate.Unix(), timeNow.Unix())
	assert.Equal(t, timeNow.Add(time.Hour*3).Unix(), lastUpdate.Unix())
}

func TestHumanTime(t *testing.T) {
	testTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	timeString := "01.01.2000 00:00"
	ht := humanTime(testTime)
	assert.Equal(t, timeString, ht)
}

func TestHumanBoolFalse(t *testing.T) {
	s := humanBool(false)
	assert.Equal(t, "Нет ❌", s)
}

func TestHumanBoolTrue(t *testing.T) {
	s := humanBool(true)
	assert.Equal(t, "Да ✅", s)
}

func TestHumanBoolReverseFalse(t *testing.T) {
	s := humanBoolReverse(false)
	assert.Equal(t, "Нет ✅", s)
}

func TestHumanBoolReverseTrue(t *testing.T) {
	s := humanBoolReverse(true)
	assert.Equal(t, "Да ❌", s)
}

func TestMergeStatusHelperCanBeMerged(t *testing.T) {
	s := mergeStatusHelper("can_be_merged")
	assert.Equal(t, "Да ✅", s)
}

func TestMergeStatusHelperOther(t *testing.T) {
	s := mergeStatusHelper("unchecked")
	assert.Equal(t, "Нет ❌", s)
}

func TestTelegramMessageTemplateNewMrWithDiffs(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	mri := gitlab.MergeRequestsInfo{
		MergeRequests: []models.MergeRequest{{
			ID:           0,
			Iid:          0,
			ProjectID:    0,
			Title:        "Test",
			Description:  "Test",
			State:        "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			TargetBranch: "test",
			SourceBranch: "test",
			Author: models.Author{
				Name: "Test",
			},
			MergeStatus:  "can_be_merged",
			WebURL:       "",
			HasConflicts: false,
			Changes: []models.FileChanges{
				{
					OldPath:       "",
					NewPath:       "test",
					IsNewFile:     true,
					IsRenamedFile: false,
					IsDeletedFile: false,
				},
			},
		}},
	}
	var err error
	assert.NotPanics(t, func() {
		err = TelegramMessageTemplateNewMrWithDiffs.Execute(buffer, mri)
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, buffer.Bytes())
}
func TestTelegramMessageTemplateNewMrWithoutDiffs(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	mri := gitlab.MergeRequestsInfo{
		MergeRequests: []models.MergeRequest{{
			ID:           0,
			Iid:          0,
			ProjectID:    0,
			Title:        "Test",
			Description:  "Test",
			State:        "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			TargetBranch: "test",
			SourceBranch: "untest",
			Author: models.Author{
				Name: "Test",
			},
			MergeStatus:  "can_be_merged",
			WebURL:       "",
			HasConflicts: false,
		}},
	}
	var err error
	assert.NotPanics(t, func() {
		err = TelegramMessageTemplateNewMrWithDiffs.Execute(buffer, mri)
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, buffer.Bytes())
}

func TestTelegramMessageTemplateWithDiffs(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	mri := gitlab.MergeRequestsInfo{
		MergeRequests: []models.MergeRequest{{
			ID:           0,
			Iid:          0,
			ProjectID:    0,
			Title:        "Test",
			Description:  "Test",
			State:        "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			TargetBranch: "test",
			SourceBranch: "untest",
			Author: models.Author{
				Name: "Test",
			},
			MergeStatus:  "can_be_merged",
			WebURL:       "",
			HasConflicts: false,
			Changes: []models.FileChanges{
				{
					OldPath:       "",
					NewPath:       "test",
					IsNewFile:     true,
					IsRenamedFile: false,
					IsDeletedFile: false,
				},
			},
		}},
	}
	var err error
	assert.NotPanics(t, func() {
		err = TelegramMessageTemplateNewMrWithDiffs.Execute(buffer, mri)
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, buffer.Bytes())
}

func TestTelegramMessageTemplateWithoutDiffs(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	mri := gitlab.MergeRequestsInfo{
		MergeRequests: []models.MergeRequest{{
			ID:           0,
			Iid:          0,
			ProjectID:    0,
			Title:        "Test",
			Description:  "Test",
			State:        "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			TargetBranch: "test",
			SourceBranch: "untest",
			Author: models.Author{
				Name: "Test",
			},
			MergeStatus:  "can_be_merged",
			WebURL:       "",
			HasConflicts: false,
		}},
	}
	var err error
	assert.NotPanics(t, func() {
		err = TelegramMessageTemplateNewMrWithDiffs.Execute(buffer, mri)
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, buffer.Bytes())
}
