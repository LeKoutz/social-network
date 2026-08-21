package models

import (
	"forum/src/db"
)

type ReactionType struct {
	db.ReactionRowType
	Post      PostType
	Comment   CommentType
}
