package dnclient

type MotdWrapper struct {
	Motd Motd `json:"motd"`
}
type Motd struct {
	DownvoteCount   int    `json:"downvote_count"`
	Message         string `json:"message"`
	UpvoteCount     int    `json:"upvote_count"`
	UserDisplayName string `json:"user_display_name"`
	UserId          int    `json:"user_id"`
}
