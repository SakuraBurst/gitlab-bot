package gitlab

type Gitlab struct {
	url       string
	repo      string
	token     string
	WithDiffs bool
}

func NewGitlabConn(withDiffs bool, repo, token, url string) Gitlab {
	return Gitlab{repo: repo, token: token, WithDiffs: withDiffs, url: url}
}
