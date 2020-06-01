package main

import (
	"github.com/wmyi/gn/config"
	"github.com/wmyi/gnchatdemo/app/chat-server/router"

	logger "github.com/wmyi/gn/glog"
	"github.com/wmyi/gn/gn"
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
	err = app.Run()
	if err != nil {
		logger.Infof("ChatApp run   error   %v \n", err)
		return
	}
	defer app.Done()
}
