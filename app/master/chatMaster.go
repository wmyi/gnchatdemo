package main

import (
	"fmt"
	"time"

	"github.com/wmyi/gn/config"
	imaster "github.com/wmyi/gn/master"
)

func main() {
	config, err := config.NewConfig("../../config/development.json")
	if err != nil {
		fmt.Println("config  error\n ", err)
		return
	}
	master, err := imaster.DefaultMaster(config)
	if err != nil {
		fmt.Printf("master error %v \n", err)
		return
	}

	master.TimeOutServerListListener(func(list []imaster.NodeInfo) {
		master.GetLogger().Infof(" timeOut List  %v ", list)

		mapNode := master.GetNodeInfos()
		for key, value := range mapNode {
			master.GetLogger().Infof(" timeOut kye  %s  time  %v   \n ", key, value.DiffTime)
		}

	})

	go func() {
		for {
			time.Sleep(3 * time.Second)
			mem, err := master.GetRunTimeMemStats("chat-001")
			if err != nil {
				master.GetLogger().Infof("err ---  ", err)
			}
			master.GetLogger().Infof("chat-001  mem  sys %d  ", mem.Sys)
		}

	}()

	err = master.Run()
	if err != nil {
		fmt.Printf("master run   error   %v \n", err)
		return
	}
}
