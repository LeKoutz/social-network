package models

import (
)

type Notification struct {
	Id        int64
	UserId    int64
	ActorId   int64
	Actor     User
	Type      string
	TargetId  int64
	Target	  any
	Timestamp string
	Read      bool
}

func (n *Notification) Add() error {
	query := `INSERT INTO notifications (user_id, actor_id, type, target_id, timestamp) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, n.UserId, n.ActorId, n.Type, n.TargetId, n.Timestamp)
	return err
}