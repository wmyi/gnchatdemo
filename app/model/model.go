package model

type userMode struct {
	Nickname string `json:"nickname"`
	UID      string `json:"uid"`
	Status   int    `json:"status"`
}

type grouprMode struct {
	name  string
	ID    string
	membs []int
}
