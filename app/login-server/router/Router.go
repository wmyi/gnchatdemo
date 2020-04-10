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
	app.APIRouter("wsclose", WsClosedHandler)

}
func InitRPCRouter(app gn.IApp) {
	app.RPCRouter("rpcGetAllUsers", rpcGetAllUsers)
	app.RPCRouter("notifyCreateGroup", notifyCreateGroup)
}

func WsClosedHandler(pack gn.IPack) {
	app := pack.GetAPP()
	if len(pack.GetBindId()) > 0 {
		userms, _ := app.GetObjectByTag("userList")
		if userMaps, ok := userms.(map[string]*model.UserMode); ok {
			if delUser, ok := userMaps[pack.GetBindId()]; ok && delUser != nil {
				// user delete
				delete(userMaps, pack.GetBindId())
				userSlice := make([]*model.UserMode, 0, len(userMaps))
				for _, item := range userMaps {
					userSlice = append(userSlice, item)
				}

				// response  to connectors
				respon := &message.ClientRes{
					Code:     "ok",
					Router:   "wsclose",
					Date:     time.Now().Format("2006-01-02 15:04:05"),
					Msg:      delUser.Nickname + "退出聊天室",
					Nickname: delUser.Nickname,
					UID:      delUser.UID,
					Users:    userSlice,
					Bridge:   []string{},
				}
				//rpc get groups
				groupMode := GetRemoteGroups(pack)
				if groupMode != nil {
					respon.Groups = groupMode
				}
				group, ok := app.GetGroup("userSession")
				if group != nil && ok {
					group.DelSession(delUser.UID)
					group.BroadCastJson(respon)
				}

			}

		}
	}

}

func notifyCreateGroup(pack gn.IPack) {
	// notify all player
	app := pack.GetAPP()
	if group, ok := app.GetGroup("userSession"); len(pack.GetData()) > 0 && ok && group != nil {
		group.BroadCast(pack.GetData())
		// finish
		pack.SetRPCRespCode(0)
	}
}

func rpcGetAllUsers(pack gn.IPack) {
	// get all groups return

	users, ok := pack.GetAPP().GetObjectByTag("userList")
	if ok && users != nil {
		if mUsers, ok := users.(map[string]*model.UserMode); ok {
			userList := make([]interface{}, 0, len(mUsers))
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

	group, ok := app.GetGroup("userSession")
	if !ok && group == nil {
		group = app.NewGroup("userSession")
	}

	if userMaps, ok := userms.(map[string]*model.UserMode); ok {
		if _, ok := userMaps[reqData.UID]; !ok {
			userMaps[reqData.UID] = &model.UserMode{
				Nickname: reqData.Nickname,
				UID:      reqData.UID,
				Status:   1,
			}
		}
		userSlice := make([]*model.UserMode, 0, len(userMaps))
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
		groupMode := GetRemoteGroups(pack)
		if groupMode != nil {
			pack.GetLogger().Infof(" groupMode     %v", groupMode)
			respon.Groups = groupMode
		}
		// broadCast other all users
		group.BroadCastJson(respon)
		// result to request  user
		pack.ResultJson(respon)

		// 绑定用户Id session
		pack.GetSession().BindId(reqData.UID)
		// 保存 在线用户的session 在 group
		group.AddSession(reqData.UID, pack.GetSession())

	}

}

func GetRemoteGroups(pack gn.IPack) []*model.GroupMode {
	app := pack.GetAPP()
	serverId, err := gnutil.RPCcalculatorServerId(pack.GetSession().GetCid(),
		app.GetServerConfig().GetServerByType("chat"))
	if err == nil {
		rpcPack, err := app.SendRPCMsg(serverId, "rpcGetAllGroups", []byte(""))
		if rpcPack.GetRPCRespCode() == 0 && err == nil {
			var groups []*model.GroupMode
			err := jsonI.Unmarshal(rpcPack.GetData(), &groups)
			if err == nil {
				return groups
			}
			app.GetLogger().Errorf("rpc  error groups  %v  error  %v ", groups, err)
		} else {
			app.GetLogger().Errorf("rpc  error code  %v  error  %v ", rpcPack.GetRPCRespCode(), err)
		}
	}
	return nil
}

// 聊天
func Chat(pack gn.IPack) {

	fmt.Printf("loginApp  Chat   pack  data %v \n", string(pack.GetData()))

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
	app := pack.GetAPP()
	if len(reqData.Bridge) > 0 {
		group, ok := app.GetGroup("userSession")
		if ok && group != nil {
			for _, uid := range reqData.Bridge {
				s, ok := group.GetSession(uid)
				if ok && s != nil {
					respon := &message.ChatRes{
						Router:   "chat",
						Date:     time.Now().Format("2006-01-02 15:04:05"),
						Msg:      reqData.Msg,
						Nickname: reqData.Nickname,
						UID:      reqData.UID,
						Bridge:   reqData.Bridge,
						GroupID:  reqData.GroupID,
						Status:   1,
					}
					// push other user msg
					app.PushJsonMsg(s, respon)
				}
			}
		}
	}
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
		if _, ok := userMaps[reqData.UID]; !ok {
			delete(userMaps, reqData.UID)
		}
		userSlice := make([]*model.UserMode, 0, len(userMaps))
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
		groupMode := GetRemoteGroups(pack)
		if groupMode != nil {

			respon.Groups = groupMode
		}
		group, ok := app.GetGroup("userSession")
		if group != nil && ok {
			group.DelSession(reqData.UID)
			group.BroadCastJson(respon)
		}
	}

}
