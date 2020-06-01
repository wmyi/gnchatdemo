package main

import (
	"github.com/wmyi/gn/config"
	logger "github.com/wmyi/gn/glog"
	"github.com/wmyi/gn/gn"
	"github.com/wmyi/gnchatdemo/app/login-server/router"
	"github.com/wmyi/gnchatdemo/app/middlerware"
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
	err = app.Run()
	if err != nil {
		logger.Infof("loginApp run   error   %v \n", err)
		return
	}
	defer app.Done()
}
