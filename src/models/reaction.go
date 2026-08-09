package models

import (
	"forum/src/db"
	"forum/src/utils"
)

type ReactionType struct {
	db.ReactionRowType
	Timestamp int64
	Post      PostType
	Comment   CommentType
}

func (r *ReactionType) ToReactionRowType() *db.ReactionRowType {
	return &db.ReactionRowType{
		Id:              r.Id,
		UserId:          r.UserId,
		PostId:          r.PostId,
		CommentId:       r.CommentId,
		TimestampString: r.TimestampString,
	}
}

func (r *ReactionType) FromReactionRowType(rr db.ReactionRowType) {
	utils.LogDebug(rr)
	r.Id = rr.Id
	r.UserId = rr.UserId
	r.PostId = rr.PostId
	r.CommentId = rr.CommentId
	r.TimestampString = rr.TimestampString
	utils.LogDebug(*r)
}
