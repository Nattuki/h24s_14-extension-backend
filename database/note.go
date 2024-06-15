package database

import (
	"h24s_14-extension-backend/util"
	"time"
)

type Note struct {
	Owner             string    `json:"-"`
	MessageId         string    `json:"-"`
	Text              string    `json:"text"`
	Color             string    `json:"color"`
	CreationTimeStamp time.Time `json:"creationTimestamp"`
	UserName          string    `json:"username"`
	MessageTime       string    `json:"time"`
	MessageText       string    `json:"messageText"`
	ChannelName       string    `json:"channelName"`
}

var (
	_all_label_field = util.GetGormFields(Note{})
)

func CreateNote(note *Note) error {
	db := GetDBConnection()
	err := db.Create(note).Error
	return err
}

func DeleteNote(owner string, messageId string) error {
	db := GetDBConnection()
	err := db.Where("owner=? && message_id=?", owner, messageId).Delete(&Note{}).Error
	return err
}

func UpdateNote(owner string, messageId string, text string, color string) error {
	db := GetDBConnection()
	err := db.Model(Note{}).Where("owner=? && message_id=?", owner, messageId).Updates(map[string]any{"text": text, "color": color}).Error
	return err
}

func GetAllNotesByOwner(owner string) error {
	db := GetDBConnection()
	var notes []Note
	err := db.Select(_all_label_field).Where("owner=?", owner).Find(&notes).Error
	return err
}
