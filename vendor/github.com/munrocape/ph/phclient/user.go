package phclient

type User struct {
	CreatedAt  string            `json:"created_at"`
	Headline   string            `json:"headline"`
	Id         int               `json:"id"`
	Images     map[string]string `json:"images"`
	Name       string            `json:"name"`
	ProfileUrl string            `json:"profile_url"`
	Username   string            `json:"username"`
	WebsiteUrl string            `json:"website_url"`
}
