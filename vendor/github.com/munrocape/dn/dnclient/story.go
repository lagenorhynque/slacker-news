package dnclient

type Stories struct {
	Stories []Story
}

type Story struct {
	Comment         string    `json:"comment"`
	Comments        []Comment `json:"comments"`
	CommentCount    int       `json:"comment_count"`
	Created_at      string    `json:"created_at"`
	Id              int       `json:"id"`
	SiteUrl         string    `json:"site_url"`
	Title           string    `json:"title"`
	Url             string    `json:"url"`
	VoteCount       int       `json:"vote_count"`
	UserDisplayName string    `json:"user_display_name"`
	UserPortraitUrl string    `json:"user_portrait_url"`
	Sponsored       bool      `json:"sponsored"`
}

type Comment struct {
	Id              int       `json:"id"`
	Body            string    `json:"body"`
	BodyHtml        string    `json:"body_html"`
	CreatedAt       string    `json:"created_at"`
	Depth           int       `json:"depth"`
	VoteCount       int       `json:"vote_count"`
	Url             string    `json:"url"`
	UserUrl         string    `json:"user_url"`
	UserId          int       `json:"user_id"`
	UserDisplayName string    `json:"user_display_name"`
	UserPortraitUrl string    `json:"user_portrait_url"`
	UserJob         string    `json:"user_job"`
	Comments        []Comment `json:"comments"`
}
