package templates

import (
	"html/template"
	"time"

	"github.com/SakuraBurst/gitlab-bot/models"
)

const CAN_BE_MERGED = "can_be_merged"

func GetRightTemplate(isNewMrMessage, withDiffs bool) *template.Template {
	if isNewMrMessage && withDiffs {
		return TelegramMessageTemplateNewMrWithDiffs
	} else if isNewMrMessage && !withDiffs {
		return TelegramMessageTemplateNewMrWithoutDiffs
	} else if !isNewMrMessage && withDiffs {
		return TelegramMessageTemplateWithDiffs
	}
	return TelegramMessageTemplateWithoutDiffs
}

var TelegramMessageTemplateWithoutDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"humanTime":         humanTime,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper}).Parse(`
Текущее количество открытых мрок на {{.On | humanTime}} - {{.Length}}
{{range .MergeRequests}}-------------------------------------
Название: {{.Title}}
Описание: {{.Description}}
Автор: {{.Author.Name}}
Создан: {{.CreatedAt | humanTime}}
Обновлен: {{.UpdatedAt | humanTime}}
ТаргетБранч: {{.TargetBranch}}
СоурсБранч: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}
<a href="{{.WebUrl}}">Ссылка на мр</a>
{{end}}
`))

var TelegramMessageTemplateWithDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"humanTime":         humanTime,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper}).Parse(`
Текущее количество открытых мрок на {{.On | humanTime}} - {{.Length}}
{{range .MergeRequests}}-------------------------------------
Название: {{.Title}}
Описание: {{.Description}}
Автор: {{.Author.Name}}
Создан: {{.CreatedAt | humanTime}}
Обновлен: {{.UpdatedAt | humanTime}}
ТаргетБранч: {{.TargetBranch}}
СоурсБранч: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}
<a href="{{.WebUrl}}">Ссылка на мр</a>
Список изменений:
{{range .Changes}}
{{.OldPath}}
{{end}}
{{end}}
`))

var TelegramMessageTemplateNewMrWithoutDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"humanTime":         humanTime,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper,
	"newMrTitle":        newMrTitle,
}).Parse(`
{{.MergeRequests | newMrTitle}}
{{range .MergeRequests}}-------------------------------------
Название: {{.Title}}
Описание: {{.Description}}
Автор: {{.Author.Name}}
Создан: {{.CreatedAt | humanTime}}
Обновлен: {{.UpdatedAt | humanTime}}
ТаргетБранч: {{.TargetBranch}}
СоурсБранч: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}
<a href="{{.WebUrl}}">Ссылка на мр</a>
{{end}}
`))

var TelegramMessageTemplateNewMrWithDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"humanTime":         humanTime,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper,
	"newMrTitle":        newMrTitle,
}).Parse(`
{{.MergeRequests | newMrTitle}}
{{range .MergeRequests}}-------------------------------------
Название: {{.Title}}
Описание: {{.Description}}
Автор: {{.Author.Name}}
Создан: {{.CreatedAt | humanTime}}
Обновлен: {{.UpdatedAt | humanTime}}
ТаргетБранч: {{.TargetBranch}}
СоурсБранч: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}
<a href="{{.WebUrl}}">Ссылка на мр</a>
Список изменений:
{{range .Changes}}
{{.OldPath}}
{{end}}
{{end}}
`))

func newMrTitle(mrs []models.MergeRequestFileChanges) string {
	if len(mrs) == 1 {
		return "Новый МР"
	} else {
		return "Новые МРы"
	}
}

func humanTime(t time.Time) string {
	return t.Format(time.RFC822)
}

func humanBool(b bool) string {
	if b {
		return "Да ✅"
	}
	return "Нет ❌"
}

func humanBoolReverse(b bool) string {
	if b {
		return "Да ❌"
	}
	return "Нет ✅"
}

func mergeStatusHelper(s string) string {
	return humanBool(s == CAN_BE_MERGED)
}
