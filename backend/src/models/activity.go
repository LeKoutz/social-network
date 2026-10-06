package models

type ActivityType struct {
	Type      string      `json:"Type,omitempty"`
	Timestamp string      `json:"Timestamp,omitempty"`
	Post      PostType    `json:"Post,omitzero"`
	Comment   CommentType `json:"Comment,omitzero"`
}
