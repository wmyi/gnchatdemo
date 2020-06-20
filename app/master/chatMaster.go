package main

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/wmyi/gn/config"
	logger "github.com/wmyi/gn/glog"
	imaster "github.com/wmyi/gn/master"
	"github.com/wmyi/gnchatdemo/app/message"
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
	// routine  master  不要阻塞主协程
	go func() {
		for {
			time.Sleep(3 * time.Second)
			req := message.CmdOnlineReq{
				Admin: "test",
				Msg:   "test test test",
			}
			results, err := master.SendCMDJson("online", "login-001", req)
			if err != nil {
				logger.Infof("master.SendCMD online error  %v   \n ", err)
				return
			}
			if len(results) > 0 {
				res := &message.CmdOnlineRes{}
				jsonI.Unmarshal(results, res)
				logger.Infof("master.SendCMD -- online  results  %v ", res)
			}

		}
	}()

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
