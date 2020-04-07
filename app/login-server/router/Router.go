package router

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/wmyi/gn/gn"
	"github.com/wmyi/gn/gnutil"
	"github.com/wmyi/gnchatdemo/app/message"
	"github.com/wmyi/gnchatdemo/app/model"
)

var (
	jsonI = jsoniter.ConfigCompatibleWithStandardLibrary
)

func InitAPIRouter(app gn.IApp) {
	app.APIRouter("login", Login)
	app.APIRouter("logout", Logout)
	app.APIRouter("chat", Chat)
}
func InitRPCRouter(app gn.IApp) {
	app.RPCRouter("rpcGetAllUsers", rpcGetAllUsers)
}

func rpcGetAllUsers(pack gn.IPack) {
	// get all groups return

	users, ok := pack.GetAPP().GetObjectByTag("userList")
	if ok && users != nil {
		if mUsers, ok := users.(map[string]*model.UserMode); ok {
			userList := make([]interface{}, len(mUsers))
			for _, item := range mUsers {
				userList = append(userList, item)
			}
			pack.ResultJson(userList)
			pack.SetRPCRespCode(0)
		} else {
			pack.SetRPCRespCode(102)
		}

		return
	}

	pack.SetRPCRespCode(101)
}

func Login(pack gn.IPack) {
	fmt.Printf("loginApp  Login   pack  data %v \n", string(pack.GetData()))

	//unmarshal json
	reqData := &message.LoginReq{}
	// request  data
	if len(pack.GetData()) > 0 {
		err := jsonI.Unmarshal(pack.GetData(), reqData)
		if err != nil {
			pack.ExceptionAbortJson("101", "解析前端数据失败 JSON ")
			return
		}
	}
	// logic
	app := pack.GetAPP()
	userms, ok := app.GetObjectByTag("userList")
	if !ok && userms == nil {
		userms = make(map[string]*model.UserMode, 1<<10)
		app.SetObjectByTag("userList", userms)
	}
	if userMaps, ok := userms.(map[string]*model.UserMode); ok {
		if _, ok := userMaps[reqData.UID]; !ok {
			userMaps[reqData.UID] = &model.UserMode{
				Nickname: reqData.Nickname,
				UID:      reqData.UID,
				Status:   1,
			}
		}
		userSlice := make([]*model.UserMode, len(userMaps))
		for _, item := range userMaps {
			userSlice = append(userSlice, item)
		}

		// response  to connectors
		respon := &message.ClientRes{
			Code:     "ok",
			Router:   "login",
			Date:     time.Now().Format("2006-01-02 15:04:05"),
			Msg:      reqData.Nickname + "加入聊天室",
			Nickname: reqData.Nickname,
			UID:      reqData.UID,
			Bridge:   reqData.Bridge,
			Users:    userSlice,
		}
		//rpc get groups
		group := GetRemoteGroups(pack)
		if group != nil {
			respon.Groups = group
		}
		pack.ResultJson(respon)
	}

}

func GetRemoteGroups(pack gn.IPack) []*model.GroupMode {
	app := pack.GetAPP()
	serverId, err := gnutil.RPCcalculatorServerId(pack.GetSession().GetCid(),
		app.GetServerConfig().GetServerByType("chat"))
	if err == nil {
		rpcPack, err := app.SendRPCMsg(serverId, "rpcGetAllGroups", pack.GetData())
		if rpcPack.GetRPCRespCode() == 0 && err == nil {
			groups := make([]*model.GroupMode, 10)
			err := jsonI.Unmarshal(rpcPack.GetData(), groups)
			if err == nil {
				return groups
			}
		} else {
			app.GetLoger().Errorf("rpc  error code  %v  error  %v ", rpcPack.GetRPCRespCode(), err)
		}
	}
	return nil
}

// 聊天
func Chat(pack gn.IPack) {

	fmt.Printf("loginApp  Login   pack  data %v \n", string(pack.GetData()))

	//unmarshal json
	reqData := &message.LoginReq{}
	// request  data
	if len(pack.GetData()) > 0 {
		err := jsonI.Unmarshal(pack.GetData(), reqData)
		if err != nil {
			pack.ExceptionAbortJson("101", "解析前端数据失败 JSON ")
			return
		}
	}
	// logic
	// response to connector
	respon := &message.ChatRes{
		Code:     "ok",
		Router:   "chat",
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Msg:      reqData.Msg,
		Nickname: reqData.Nickname,
		UID:      reqData.UID,
		Bridge:   reqData.Bridge,
		GroupID:  reqData.GroupID,
		Status:   1,
	}
	pack.ResultJson(respon)
}

func Logout(pack gn.IPack) {
	fmt.Printf("loginApp  logout   pack  data %v \n", string(pack.GetData()))
	//unmarshal json
	reqData := &message.LoginReq{}
	// request  data
	if len(pack.GetData()) > 0 {
		err := jsonI.Unmarshal(pack.GetData(), reqData)
		if err != nil {
			pack.ExceptionAbortJson("101", "解析前端数据失败 JSON ")
			return
		}
	}
	// logic
	app := pack.GetAPP()
	userms, _ := app.GetObjectByTag("userList")
	if userMaps, ok := userms.(map[string]*model.UserMode); ok {
		if user, ok := userMaps[reqData.UID]; !ok {
			user.Status = 0
		}
		userSlice := make([]*model.UserMode, len(userMaps))
		for _, item := range userMaps {
			userSlice = append(userSlice, item)
		}

		// response  to connectors
		respon := &message.ClientRes{
			Code:     "ok",
			Router:   "logout",
			Date:     time.Now().Format("2006-01-02 15:04:05"),
			Msg:      reqData.Nickname + "退出聊天室",
			Nickname: reqData.Nickname,
			UID:      reqData.UID,
			Bridge:   reqData.Bridge,
			Users:    userSlice,
		}
		//rpc get groups
		group := GetRemoteGroups(pack)
		if group != nil {
			respon.Groups = group
		}
		pack.ResultJson(respon)
	}

}
