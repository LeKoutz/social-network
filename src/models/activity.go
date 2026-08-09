package models

type ActivityType struct {
	Type            string
	TimestampString string
	Post            PostType
	Comment         CommentType
}
