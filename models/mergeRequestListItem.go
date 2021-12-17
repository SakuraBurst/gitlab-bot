package models

import "time"

type Author struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UserName  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type MergeRequestListItem struct {
	Id           int       `json:"id"`
	Iid          int       `json:"iid"`
	ProjectId    int       `json:"project_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TargetBranch string    `json:"target_branch"`
	SourceBranch string    `json:"source_branch"`
	Author       Author    `json:"author"`
	IsWip        bool      `json:"work_in_progress"`
	MergeStatus  string    `json:"merge_status"`
	Sha          string    `json:"sha"`
	WebUrl       string    `json:"web_url"`
	HasConflicts bool      `json:"has_conflicts"`
}
