package models

type BlogPost struct {
	Id       int
	Title    string
	Content  string
	Comments []*Comment
}
