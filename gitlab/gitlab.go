package gitlab

type Gitlab struct {
	repo      string
	token     string
	WithDiffs bool
}

func NewGitlabConn(withDiffs bool, repo, token string) Gitlab {
	return Gitlab{repo: repo, token: token, WithDiffs: withDiffs}
}
