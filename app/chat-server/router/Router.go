package router

import (
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	logger "github.com/wmyi/gn/glog"
	"github.com/wmyi/gn/gn"
	"github.com/wmyi/gn/gnutil"
	"github.com/wmyi/gnchatdemo/app/message"
	"github.com/wmyi/gnchatdemo/app/model"
)

var (
	jsonI = jsoniter.ConfigCompatibleWithStandardLibrary
)

func InitAPIRouter(app gn.IApp) {
	app.APIRouter("addGroup", true, addGroup)
	app.APIRouter("createGroup", true, createGroup)
	app.APIRouter("chatGroup", false, chatGroup)
}

func InitRPCRouter(app gn.IApp) {
	app.RPCRouter("rpcGetAllGroups", false, rpcGetAllGroups)
}

func rpcGetAllGroups(pack gn.IPack) {
	// get all groups return
	groups, ok := pack.GetAPP().GetObjectByTag("groups")
	if ok && groups != nil {
		if mgroups, ok := groups.(map[string]*model.GroupMode); ok {
			groupList := make([]interface{}, 0, len(mgroups))
			for _, item := range mgroups {
				groupList = append(groupList, item)
			}
			pack.ResultJson(groupList)
			pack.SetRPCRespCode(0)
		} else {
			pack.SetRPCRespCode(102)
		}
		return
	}

	pack.SetRPCRespCode(101)
}

func GetRemoteUsers(pack gn.IPack) []*model.UserMode {
	app := pack.GetAPP()
	serverId, err := gnutil.RPCcalculatorServerId(pack.GetSession().GetCid(),
		app.GetServerConfig().GetServerByType("login"))
	if err == nil {
		rpcPack, err := app.RequestRPCMsg(serverId, "rpcGetAllUsers", pack.GetData())
		if rpcPack.GetRPCRespCode() == 0 && err == nil {
			var users []*model.UserMode
			err := jsonI.Unmarshal(rpcPack.GetData(), &users)
			if err == nil {
				return users
			}
		} else {
			logger.Errorf("rpc  error code  %v  error  %v ", rpcPack.GetRPCRespCode(), err)
		}
	}
	return nil
}

func addGroup(pack gn.IPack) {
	logger.Infof("chatApp  addGroup   pack  data %v \n", string(pack.GetData()))
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
	groups, _ := app.GetObjectByTag("groups")
	if mgroups, ok := groups.(map[string]*model.GroupMode); ok {
		if group, ok := mgroups[reqData.GroupID]; ok && group != nil {
			if group.Users != nil {
				group.Users = append(group.Users, &model.UserMode{
					UID:      reqData.UID,
					Nickname: reqData.Nickname,
				})
			}

		}

		groupList := make([]*model.GroupMode, 0, len(mgroups))
		for _, item := range mgroups {
			groupList = append(groupList, item)
		}
		// response  to connectors
		respon := &message.ClientRes{
			Code:     "ok",
			Router:   "addgroup",
			Date:     time.Now().Format("2006-01-02 15:04:05"),
			Msg:      reqData.Nickname + "加入了群" + reqData.GroupName,
			Nickname: reqData.Nickname,
			UID:      reqData.UID,
			Bridge:   reqData.Bridge,
			Groups:   groupList,
		}
		//rpc get groups
		userList := GetRemoteUsers(pack)
		if userList != nil {
			respon.Users = userList
		}
		pack.ResultJson(respon)
		// group session broadcast otheruser
		if g, ok := app.GetGroup(reqData.GroupID); ok && g != nil {
			g.BroadCastJson(respon) // broadcast other user
			g.AddSession(reqData.UID, pack.GetSession())
		}
	}
}

func chatGroup(pack gn.IPack) {
	logger.Infof("chatApp  chatGroup   pack  data %v \n", string(pack.GetData()))
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
	if len(reqData.GroupID) > 0 {
		group, ok := app.GetGroup(reqData.GroupID)
		if ok && group != nil {
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
			group.BroadCastJson(respon)
		}
	}
}

func createGroup(pack gn.IPack) {
	logger.Infof("chatApp  addGroup   pack  data %v \n", string(pack.GetData()))
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
	groups, ok := app.GetObjectByTag("groups")
	if !ok && groups == nil {
		groups = make(map[string]*model.GroupMode, 1<<10)
		app.SetObjectByTag("groups", groups)
	}
	if mgroups, ok := groups.(map[string]*model.GroupMode); ok {
		// group  data
		groupId := strconv.FormatInt(time.Now().Unix(), 10)
		mgroups[groupId] = &model.GroupMode{
			ID:   groupId,
			Name: reqData.GroupName,
			Users: []*model.UserMode{&model.UserMode{
				UID:      reqData.UID,
				Nickname: reqData.Nickname,
			}},
		}

		groupList := make([]*model.GroupMode, 0, len(mgroups))
		for _, item := range mgroups {
			groupList = append(groupList, item)
		}

		// group session
		var sgroup *gn.Group
		if g, ok := app.GetGroup(groupId); !ok && g == nil {
			sgroup = app.NewGroup(groupId)
			sgroup.AddSession(reqData.UID, pack.GetSession())
		}

		// response  to connectors
		respon := &message.ClientRes{
			Code:     "ok",
			Router:   "creategroup",
			Date:     time.Now().Format("2006-01-02 15:04:05"),
			Msg:      reqData.Nickname + "创建了群" + reqData.GroupName,
			Nickname: reqData.Nickname,
			UID:      reqData.UID,
			Bridge:   reqData.Bridge,
			Groups:   groupList,
		}

		serverId, err := gnutil.RPCcalculatorServerId(pack.GetSession().GetCid(),
			app.GetServerConfig().GetServerByType("login"))
		if err == nil {
			err := app.NotifyRPCJsonMsg(serverId, "notifyCreateGroup", respon)
			if err != nil {
				logger.Errorf("rpc  NotifyRPCJsonMsg error %v ", err)
			}
		}

	}

}
