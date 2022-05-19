package templates

import (
	"bytes"
	"github.com/SakuraBurst/gitlab-bot/pkg/gitlab"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, CanBeMerged, "can_be_merged")
	assert.Equal(t, HumanTimeFormat, "02.01.2006 15:04")
}

func TestGetRightTemplate_TelegramMessageTemplateNewMrWithDiffs(t *testing.T) {
	template := GetRightTemplate(true, true)
	assert.Equal(t, template.Name(), "TelegramMessageTemplateNewMrWithDiffs")
}

func TestGetRightTemplate_TelegramMessageTemplateNewMrWithoutDiffs(t *testing.T) {
	template := GetRightTemplate(true, false)
	assert.Equal(t, template.Name(), "TelegramMessageTemplateNewMrWithoutDiffs")
}

func TestGetRightTemplate_TelegramMessageTemplateWithDiffs(t *testing.T) {
	template := GetRightTemplate(false, true)
	assert.Equal(t, template.Name(), "TelegramMessageTemplateWithDiffs")
}

func TestGetRightTemplate_TelegramMessageTemplateWithoutDiffs(t *testing.T) {
	template := GetRightTemplate(false, false)
	assert.Equal(t, template.Name(), "TelegramMessageTemplateWithoutDiffs")
}

func TestNewMrTitleZeroLength(t *testing.T) {
	title := newMrTitle(nil)
	assert.Equal(t, title, "Новый MR")
}

func TestNewMrTitleLengthOne(t *testing.T) {
	title := newMrTitle([]models.MergeRequest{{}})
	assert.Equal(t, title, "Новый MR")
}

func TestNewMrTitleLengthTwo(t *testing.T) {
	title := newMrTitle([]models.MergeRequest{{}, {}})
	assert.Equal(t, title, "Новые MR'ы")
}

func TestLastUpdateUpdatedAfterCreated(t *testing.T) {
	created := time.Now()
	updated := time.Now().Add(time.Hour)
	lastUpdate := lastUpdate(created, updated)
	assert.Greater(t, lastUpdate.Unix(), updated.Unix())
	assert.Equal(t, lastUpdate.Unix(), updated.Add(time.Hour*3).Unix())
	assert.Greater(t, lastUpdate.Unix(), created.Add(time.Hour*3).Unix())
}

func TestLastUpdateCreatedEqualsUpdated(t *testing.T) {
	timeNow := time.Now()
	lastUpdate := lastUpdate(timeNow, timeNow)
	assert.Greater(t, lastUpdate.Unix(), timeNow.Unix())
	assert.Equal(t, lastUpdate.Unix(), timeNow.Add(time.Hour*3).Unix())
}

func TestHumanTime(t *testing.T) {
	testTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	timeString := "01.01.2000 00:00"
	ht := humanTime(testTime)
	assert.Equal(t, ht, timeString)
}

func TestHumanBoolFalse(t *testing.T) {
	s := humanBool(false)
	assert.Equal(t, s, "Нет ❌")
}

func TestHumanBoolTrue(t *testing.T) {
	s := humanBool(true)
	assert.Equal(t, s, "Да ✅")
}

func TestHumanBoolReverseFalse(t *testing.T) {
	s := humanBoolReverse(false)
	assert.Equal(t, s, "Нет ✅")
}

func TestHumanBoolReverseTrue(t *testing.T) {
	s := humanBoolReverse(true)
	assert.Equal(t, s, "Да ❌")
}

func TestMergeStatusHelperCanBeMerged(t *testing.T) {
	s := mergeStatusHelper("can_be_merged")
	assert.Equal(t, s, "Да ✅")
}

func TestMergeStatusHelperOther(t *testing.T) {
	s := mergeStatusHelper("unchecked")
	assert.Equal(t, s, "Нет ❌")
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
