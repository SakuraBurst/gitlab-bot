package models

type FileChanges struct {
	OldPath       string `json:"old_path"`
	NewPath       string `json:"new_path"`
	IsNewFile     bool   `json:"new_file"`
	IsRenamedFile bool   `json:"renamed_file"`
	IsDeletedFile bool   `json:"deleted_file"`
}

type MergeRequestFileChanges struct {
	MergeRequestListItem
	Changes []FileChanges `json:"changes"`
}
