package models

import (
	"forum/src/db"
)

type ReactionType struct {
	db.ReactionRowType
	Timestamp int64
	Post      PostType
	Comment   CommentType
}
