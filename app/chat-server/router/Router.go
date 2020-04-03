package router

import (
	"fmt"

	"github.com/wmyi/gn/gn"
	"github.com/wmyi/gnchatdemo/app/message"
)

func InitAPIRouter(app gn.IApp) {
	app.APIRouter("groupChat", groupChat)
	app.RPCRouter("kickGroup", kickGroup)
}

func groupChat(pack gn.IPack) {
	fmt.Printf("chatApp  onChat   pack  data %v \n", string(pack.GetData()))

	respon := &message.OnChatRespon{
		Code: "ok",
		Msg:  string(pack.GetData()),
	}
	pack.ResultJson(respon)
}

func kickGroup(pack gn.IPack) {

	fmt.Printf(" RPCchatApp  kickGroup   pack  data %v \n", string(pack.GetData()))
	response := &message.OnChatRespon{
		Code: "ok",
		Msg:  "kickGroup success ",
	}
	// time.Sleep(5 * time.Second)
	pack.ResultJson(response)
}
func InitRPCRouter(app gn.IApp) {

}
