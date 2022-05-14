package models

type GitlabError struct {
	Message string `json:"message"`
}

func (g GitlabError) Error() string {
	return g.Message
}
