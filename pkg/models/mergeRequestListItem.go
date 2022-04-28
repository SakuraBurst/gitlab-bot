package models

import "time"

type Author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserName  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type MergeRequestListItem struct {
	ID           int       `json:"id"`
	Iid          int       `json:"iid"`
	ProjectID    int       `json:"project_id"`
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
	WebURL       string    `json:"web_url"`
	HasConflicts bool      `json:"has_conflicts"`
}
