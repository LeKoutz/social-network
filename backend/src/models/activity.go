package models

type ActivityType struct {
	Type            string
	Timestamp       string
	Post            PostType
	Comment         CommentType
}
