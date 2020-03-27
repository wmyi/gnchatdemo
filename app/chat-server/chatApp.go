package main

import (
	"fmt"
	"gn/config"
	message "gnchatdemo/app/mssage"
	"time"

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
	app.APIRouter("onchat", onChat)
	app.RPCRouter("kickGroup", kickGroup)

	err = app.Run()
	if err != nil {
		fmt.Printf("ChatApp run   error   %v \n", err)
		return
	}
}

func onChat(pack gn.IPack) {
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
	time.Sleep(5 * time.Second)
	pack.ResultJson(response)

}
