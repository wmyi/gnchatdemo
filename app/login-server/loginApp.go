package main

import (
	"fmt"

	"github.com/wmyi/gn/config"
	"github.com/wmyi/gnchatdemo/app/login-server/router"
	"github.com/wmyi/gnchatdemo/app/middlerware"

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
	router.InitAPIRouter(app)
	router.InitRPCRouter(app)

	app.AddConfigFile("test", "../../config/", "yaml")
	app.UseMiddleWare(&middlerware.PackTimer{})
	err = app.Run()
	if err != nil {
		fmt.Printf("loginApp run   error   %v \n", err)
		return
	}

}
