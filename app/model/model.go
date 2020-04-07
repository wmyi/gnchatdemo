package model

type UserMode struct {
	Nickname string `json:"nickname"`
	UID      string `json:"uid"`
	Status   int    `json:"status"`
}

type GroupMode struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Users []*UserMode `json:"users"`
}
