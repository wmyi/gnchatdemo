package message

import "github.com/wmyi/gnchatdemo/app/model"

type LoginReq struct {
	UID       string   `json:"uid"`
	Router    string   `json:"router"`
	Nickname  string   `json:"nickname"`
	Msg       string   `json:"msg"`
	Bridge    []string `json:"bridge"`
	GroupID   string   `json:"groupId"`
	GroupName string   `json:"groupName"`
}

type ClientRes struct {
	Code     string             `json:"code"`
	Router   string             `json:"router"`
	Date     string             `json:"date"`
	Msg      string             `json:"msg"`
	Users    []*model.UserMode  `json:"users"`
	Groups   []*model.GroupMode `json:"groups"`
	UID      string             `json:"uid"`
	Nickname string             `json:"nickname"`
	Bridge   []string           `json:"bridge"`
}

type CmdOnlineReq struct {
	Admin string `json:"admin"`
	Msg   string `json:"msg"`
}

type CmdOnlineRes struct {
	ServerId string `json:"serverId"`
	Msg      string `json:"msg"`
}

type ChatRes struct {
	UID      string   `json:"uid"`
	Router   string   `json:"router"`
	Nickname string   `json:"nickname"`
	Msg      string   `json:"msg"`
	Bridge   []string `json:"bridge"`
	GroupID  string   `json:"groupId"`
	Date     string   `json:"date"`
	Status   int      `json:"status"`
}
