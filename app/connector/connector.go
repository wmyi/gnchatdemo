package main

import (
	"fmt"
	"hash/crc32"

	"github.com/wmyi/gn/config"
	"github.com/wmyi/gn/connector"
	"github.com/wmyi/gn/gnError"
)

func main() {
	config, err := config.NewConfig("../../config/development.json")
	if err != nil {
		fmt.Println("config  error ", err)
		return
	}
	// fmt.Printf("config:      %v  \n  ", config)
	connector, err := connector.DefaultConnector(config)
	if err != nil {
		fmt.Println("new DefaultConnect  error ", err)
		return
	}
	// exception handler
	connector.AddExceptionHandler(func(exception *gnError.GnException) {
		// close handler push msg
		if exception.Exception == gnError.WS_CLOSED && len(exception.BindId) > 0 && len(exception.Id) > 0 {
			handlerName := "wsclose"
			serverAddress := connector.GetServerIdByRouter(handlerName, exception.BindId, exception.Id,
				config.GetServerByType("login"))
			connector.SendPack(serverAddress, handlerName, exception.BindId, exception.Id, nil)
		}
	})

	// set pack  route
	connector.AddRouterRearEndHandler("connector", connectorRoure)
	err = connector.Run()
	if err != nil {
		fmt.Printf("connector run   error   %v \n", err)
		return
	}

}

func connectorRoure(cid string, bindId string, serverList []*config.ServersConfig) string {
	if len(bindId) > 0 {
		index := int(crc32.ChecksumIEEE([]byte(bindId))) % len(serverList)
		if serverList[index] != nil {
			return serverList[index].ID
		}
	}
	return ""
}
