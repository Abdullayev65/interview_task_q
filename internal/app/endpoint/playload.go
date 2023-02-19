package endpoint

import "github.com/Abdullayev65/interview_task_q/internal/app/moduls"

type Sign struct {
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type PostInput struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type PostOutput struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type CommentInput struct {
	PostID int    `json:"postId"`
	Text   string `json:"text"`
}
type CommentOutput struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	PostID int    `json:"postId"`
	Text   string `json:"text"`
}

func (i *PostInput) toModule(userId int) moduls.Post {
	return moduls.Post{UserID: userId,
		Title: i.Title, Text: i.Text}
}
func newPostOutput(p *moduls.Post) *PostOutput {
	return &PostOutput{ID: p.ID, UserID: p.UserID,
		Title: p.Title, Text: p.Text}
}

func newCommentOutPut(c moduls.Comment) CommentOutput {
	return CommentOutput{
		ID:     c.ID,
		UserID: c.UserID,
		PostID: c.PostID,
		Text:   c.Text}
}

func (i *CommentInput) toModule(userId int) moduls.Comment {
	return moduls.Comment{UserID: userId,
		PostID: i.PostID, Text: i.Text}
}
