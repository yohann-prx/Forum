package model

type PageData struct {
	User       *User
	Posts      []*Post
	Categories []*Category
}
