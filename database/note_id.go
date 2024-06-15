package database

type NoteId struct {
	UserName    string `json:"username"`
	MessageTime string `json:"time"`
	MessageText string `json:"messageText"`
	ChannelName string `json:"channelName"`
	Owner       string `json:"-"`
}
