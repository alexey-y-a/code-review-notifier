package model

type User struct {
	TelegramID  int64
	GitHubLogin string
}

type PullRequestEvent struct {
	Action     string
	Assignee   string
	Title      string
	HTMLURL    string
	Repository string
}
