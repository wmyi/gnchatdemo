package router

import (
	"fmt"

	"github.com/wmyi/gn/gn"
	"github.com/wmyi/gn/gnutil"
	"github.com/wmyi/gnchatdemo/app/message"
)

func InitAPIRouter(app gn.IApp) {
	app.APIRouter("login", Login)
	app.APIRouter("logout", Logout)
	app.APIRouter("chat", Chat)
}

func Login(pack gn.IPack) {
	fmt.Printf("loginApp  Login   pack  data %v \n", string(pack.GetData()))

	respon := &message.LoginRespon{
		Code:     "ok",
		Login:    true,
		UID:      1000,
		Nickname: "张三",
	}
	// time.Sleep(15 * time.Second)
	pack.ResultJson(respon)

}

// 单聊
func Chat(pack gn.IPack) {

	fmt.Printf("loginApp  Login   pack  data %v \n", string(pack.GetData()))

	respon := &message.LoginRespon{
		Code:     "ok",
		Login:    true,
		UID:      1000,
		Nickname: "张三",
	}
	// time.Sleep(15 * time.Second)
	pack.ResultJson(respon)

}

func Logout(pack gn.IPack) {
	fmt.Printf("loginApp  logout   pack  data %v \n", string(pack.GetData()))
	response := &message.LogoutRespon{
		Code:   "ok",
		Logout: true,
	}

	app := pack.GetAPP()

	serverId, err := gnutil.RPCcalculatorServerId(pack.GetSession().GetCid(),
		app.GetServerConfig().GetServerByType("chat"))
	fmt.Printf("loginApp  rpc   serverId     %v \n", serverId)
	if err == nil {
		rpcPack, err := app.SendRPCMsg(serverId, "kickGroup", pack.GetData())
		if err == nil {
			fmt.Printf("loginApp  rpc   kickGroup    data %v \n", string(rpcPack.GetData()))
			pack.ResultJson(response)
		} else {

			fmt.Printf("loginApp  rpc   kickGroup    error %v \n", err.Error())
			response.Code = "500"
			response.Logout = false
			pack.ResultJson(response)
		}

	}

}
func InitRPCRouter(app gn.IApp) {
}
