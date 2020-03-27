package message

type LoginRespon struct {
	Code     string `json:"code"`
	Login    bool   `json:"login"`
	UID      int    `json:"uid"`
	Nickname string `json:"nickname"`
}

type LogoutRespon struct {
	Code   string `json:"code"`
	Logout bool   `json:"logout"`
}

type OnChatRespon struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
