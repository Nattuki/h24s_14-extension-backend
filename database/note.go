package database

import "time"

type Note struct {
	Id                NoteId    `json:"id" gorm:"foreignkey:Owner"`
	Owner             string    `json:"-"`
	Text              string    `json:"text"`
	Color             string    `json:"color"`
	CreationTimeStamp time.Time `json:"creationTimestamp"`
}
