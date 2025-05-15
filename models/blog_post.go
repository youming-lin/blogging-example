package models

type BlogPost struct {
	Id               int
	Title            string
	Content          string
	NumberOfComments uint
	Comments         []*Comment `json:",omitempty"`
}
