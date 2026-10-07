package models

import (
	"forum/src/db"
	"forum/src/ferror"
)

type GroupMessageType struct {
	db.GroupMessageRowType

	SenderUsername string
	MemberIds      []int64 `json:"-"`
}

func (m *GroupMessageType) Add() (int64, error) {
	id, err := m.InsertGroupMessage()
	if err != nil {
		return 0, ferror.ReturnErr(err)
	}
	return id, nil
}

func (m *GroupMessageType) MarkAsRead(userId int64) error {
	if err := m.UpdateLastReadMessageId(userId); err != nil {
		return ferror.ReturnErr(err)
	}
	return nil
}
