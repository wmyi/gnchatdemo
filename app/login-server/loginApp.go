package main

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/wmyi/gn/config"
	logger "github.com/wmyi/gn/glog"
	"github.com/wmyi/gn/gn"
	"github.com/wmyi/gnchatdemo/app/login-server/router"
	"github.com/wmyi/gnchatdemo/app/message"
	"github.com/wmyi/gnchatdemo/app/middlerware"
)

var (
	jsonI = jsoniter.ConfigCompatibleWithStandardLibrary
)

func main() {

	config, err := config.NewConfig("../../config/development.json")
	if err != nil {
		logger.Infof("config  error\n ", err)
		return
	}
	app, err := gn.DefaultApp(config)
	if err != nil {
		logger.Infof("new APP error %v \n", err)
		return
	}
	router.InitAPIRouter(app)
	router.InitRPCRouter(app)

	app.AddConfigFile("test", "../../config/", "yaml")
	app.UseMiddleWare(&middlerware.PackTimer{})

	// master cmd  command
	app.CMDHandler("online", func(pack gn.IPack) {

		req := &message.CmdOnlineReq{}
		logger.Infof("online Cmd--    %v \n", string(pack.GetData()))
		if err = jsonI.Unmarshal(pack.GetData(), req); err == nil {
			logger.Infof("  online    req    %v \n ", req)
			if req.Admin == "test" {
				response := &message.CmdOnlineRes{
					ServerId: "login-001",
					Msg:      "onLine test  test",
				}
				logger.Infof("  online    response    %v \n ", response)
				pack.ResultJson(response)
			}
		}
	})

	err = app.Run()
	if err != nil {
		logger.Infof("loginApp run   error   %v \n", err)
		return
	}
	defer app.Done()
}
