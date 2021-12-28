package templates

import (
	"html/template"
	"time"

	"github.com/SakuraBurst/gitlab-bot/models"
)

const CanBeMerged = "can_be_merged"

func GetRightTemplate(isNewMrMessage, withDiffs bool) *template.Template {
	switch {
	case isNewMrMessage && withDiffs:
		return TelegramMessageTemplateNewMrWithDiffs
	case isNewMrMessage && !withDiffs:
		return TelegramMessageTemplateNewMrWithoutDiffs
	case !isNewMrMessage && withDiffs:
		return TelegramMessageTemplateWithDiffs
	}
	return TelegramMessageTemplateWithoutDiffs
}

var TelegramMessageTemplateWithoutDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"humanTime":         humanTime,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper,
	"lastUpdate":        lastUpdate,
}).Parse(`
Текущее количество открытых MR на {{.On | humanTime}} - {{.Length}}
{{range .MergeRequests}}------------------------------------
<b>{{.Title}}</b>, {{lastUpdate .CreatedAt .UpdatedAt }}, 
<i>{{.Description}}</i>

Автор: {{.Author.Name}}
В ветку: {{.TargetBranch}}
Из ветки: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}

<a href="{{.WebURL}}">Ссылка на MR</a>
{{end}}
`))

var TelegramMessageTemplateWithDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"humanTime":         humanTime,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper,
	"lastUpdate":        lastUpdate,
}).Parse(`
Текущее количество открытых MR на {{.On | humanTime}} - {{.Length}}
{{range .MergeRequests}}------------------------------------
<b>{{.Title}}</b>, {{lastUpdate .CreatedAt .UpdatedAt }}, 
<i>{{.Description}}</i>

Автор: {{.Author.Name}}
В ветку: {{.TargetBranch}}
Из ветки: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}

<a href="{{.WebURL}}">Ссылка на MR</a>
Список изменений:
{{range .Changes}}
{{.OldPath}}
{{end}}
{{end}}
`))

var TelegramMessageTemplateNewMrWithoutDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"lastUpdate":        lastUpdate,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper,
	"newMrTitle":        newMrTitle,
}).Parse(`
{{.MergeRequests | newMrTitle}}
{{range .MergeRequests}}------------------------------------
<b>{{.Title}}</b>, {{lastUpdate .CreatedAt .UpdatedAt }}, 
<i>{{.Description}}</i>

Автор: {{.Author.Name}}
В ветку: {{.TargetBranch}}
Из ветки: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}

<a href="{{.WebURL}}">Ссылка на MR</a>
{{end}}
`))

var TelegramMessageTemplateNewMrWithDiffs = template.Must(template.New("mr").Funcs(template.FuncMap{
	"lastUpdate":        lastUpdate,
	"humanBool":         humanBool,
	"humanBoolReverse":  humanBoolReverse,
	"mergeStatusHelper": mergeStatusHelper,
	"newMrTitle":        newMrTitle,
}).Parse(`
{{.MergeRequests | newMrTitle}}
{{range .MergeRequests}}------------------------------------
<b>{{.Title}}</b>, {{lastUpdate .CreatedAt .UpdatedAt }}, 
<i>{{.Description}}</i>

Автор: {{.Author.Name}}
В ветку: {{.TargetBranch}}
Из ветки: {{.SourceBranch}}
Есть ли конфликты: {{.HasConflicts | humanBoolReverse}}
Можно ли мержить: {{.MergeStatus | mergeStatusHelper}}

<a href="{{.WebURL}}">Ссылка на MR</a>
Список изменений:
{{range .Changes}}
{{.OldPath}}
{{end}}
{{end}}
`))

func newMrTitle(mrs []models.MergeRequestFileChanges) string {
	if len(mrs) == 1 {
		return "Новый MR"
	}
	return "Новые MR'ы"
}

func lastUpdate(created, updated time.Time) string {
	created = created.Add(time.Hour * 3)
	updated = updated.Add(time.Hour * 3)
	if updated.After(created) {
		return humanTime(updated)
	}
	return humanTime(created)
}

func humanTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
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
	return humanBool(s == CanBeMerged)
}
