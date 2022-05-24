package templates

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/SakuraBurst/gitlab-bot/pkg/services/gitlab"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateStringFromTemplatePanic(t *testing.T) {
	assert.Panics(t, func() {
		CreateStringFromTemplate(TelegramMessageTemplateWithoutDiffs, nil)
	})
}

func TestCreateStringFromTemplate(t *testing.T) {
	templateStringMock := "\nТекущее количество открытых MR на 01.01.2020 00:00 - 1\n------------------------------------\n<b></b>, 01.01.0001 03:00, \n<i></i>\n\nАвтор: \nВ ветку: \nИз ветки: \nЕсть ли конфликты: Нет ✅\nМожно ли мержить: Нет ❌\n\n<a href=\"\">Ссылка на MR</a>\n\n"
	mri := gitlab.MergeRequestsInfo{
		Length:        1,
		On:            time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		MergeRequests: []models.MergeRequest{{}},
	}
	s := CreateStringFromTemplate(TelegramMessageTemplateWithoutDiffs, mri)
	assert.NotEmpty(t, s)
	assert.Equal(t, templateStringMock, s)
}
