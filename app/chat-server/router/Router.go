package router

import (
	"fmt"
	"strconv"
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
	app.APIRouter("addGroup", addGroup)
	app.APIRouter("createGroup", createGroup)
}

func InitRPCRouter(app gn.IApp) {
	app.RPCRouter("rpcGetAllGroups", rpcGetAllGroups)
}

func rpcGetAllGroups(pack gn.IPack) {
	// get all groups return

	groups, ok := pack.GetAPP().GetObjectByTag("groups")
	if ok && groups != nil {
		if mgroups, ok := groups.(map[string]*model.GroupMode); ok {
			groupList := make([]interface{}, len(mgroups))
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
		rpcPack, err := app.SendRPCMsg(serverId, "rpcGetAllUsers", pack.GetData())
		if rpcPack.GetRPCRespCode() == 0 && err == nil {
			users := make([]*model.UserMode, 10)
			err := jsonI.Unmarshal(rpcPack.GetData(), users)
			if err == nil {
				return users
			}
		} else {
			app.GetLoger().Errorf("rpc  error code  %v  error  %v ", rpcPack.GetRPCRespCode(), err)
		}
	}
	return nil
}

func addGroup(pack gn.IPack) {
	fmt.Printf("chatApp  addGroup   pack  data %v \n", string(pack.GetData()))
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

		groupList := make([]*model.GroupMode, len(mgroups))
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

	}
}

func createGroup(pack gn.IPack) {
	fmt.Printf("chatApp  addGroup   pack  data %v \n", string(pack.GetData()))
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
		groupId := strconv.FormatInt(time.Now().Unix(), 10)
		mgroups[groupId] = &model.GroupMode{
			ID:   groupId,
			Name: reqData.GroupName,
			Users: []*model.UserMode{&model.UserMode{
				UID:      reqData.UID,
				Nickname: reqData.Nickname,
			}},
		}

		groupList := make([]*model.GroupMode, len(mgroups))
		for _, item := range mgroups {
			groupList = append(groupList, item)
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
		//rpc get groups
		userList := GetRemoteUsers(pack)
		if userList != nil {
			respon.Users = userList
		}
		pack.ResultJson(respon)

	}

}
