package main

import (
	"time"

	"github.com/wmyi/gn/config"
	logger "github.com/wmyi/gn/glog"
	imaster "github.com/wmyi/gn/master"
)

func main() {
	config, err := config.NewConfig("../../config/development.json")
	if err != nil {
		logger.Infof("config  error\n ", err)
		return
	}
	master, err := imaster.DefaultMaster(config)
	if err != nil {
		logger.Infof("master error %v \n", err)
		return
	}

	master.TimeOutServerListListener(func(list []imaster.NodeInfo) {
		logger.Infof(" timeOut List  %v ", list)

		mapNode := master.GetNodeInfos()
		for key, value := range mapNode {
			logger.Infof(" timeOut kye  %s  time  %v   \n ", key, value.DiffTime)
		}

	})

	go func() {
		for {
			time.Sleep(3 * time.Second)
			mem, err := master.GetRunTimeMemStats("chat-001")
			if err != nil {
				logger.Infof("err ---  ", err)
			}
			logger.Infof("chat-001  mem  sys %d  ", mem.Sys)
		}

	}()

	err = master.Run()
	if err != nil {
		logger.Infof("master run   error   %v \n", err)
		return
	}
}
