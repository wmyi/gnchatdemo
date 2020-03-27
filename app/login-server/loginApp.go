package main

import (
	"fmt"
	message "github.com/wmyi/gnchatdemo/app/mssage"
	"time"

	"github.com/wmyi/gn/config"
	"github.com/wmyi/gn/gnutil"

	"github.com/wmyi/gn/gn"
)

func main() {

	config, err := config.NewConfig("../../config/development.json")
	if err != nil {
		fmt.Println("config  error\n ", err)
		return
	}
	app, err := gn.DefaultApp(config)
	if err != nil {
		fmt.Printf("new APP error %v \n", err)
		return
	}
	app.APIRouter("login", Login)
	app.APIRouter("logout", Logout)

	err = app.Run()
	if err != nil {
		fmt.Printf("loginApp run   error   %v \n", err)
		return
	}

}

func Login(pack gn.IPack) {

	fmt.Printf("loginApp  Login   pack  data %v \n", string(pack.GetData()))
	respon := &message.LoginRespon{
		Code:     "ok",
		Login:    true,
		UID:      1000,
		Nickname: "张三",
	}
	time.Sleep(15 * time.Second)
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
