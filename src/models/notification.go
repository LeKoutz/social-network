package models

import (
)

type Notification struct {
	Id        int64
	UserId    int64
	ActorId   int64
	Actor     User
	Type      string
	Post      Post
	PostId    int64
	Comment	  Comment
	CommentId int64
	TimestampString string
	Read      bool
}

func (n *Notification) Add() error {
	query := `INSERT INTO notifications (user_id, actor_id, type, post_id, comment_id, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, n.UserId, n.ActorId, n.Type, n.PostId, n.CommentId, n.TimestampString)
	return err
}
