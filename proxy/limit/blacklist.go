package limit

import (
	"fmt"
	"log"
	"sync"
)

// 限制管理
type SecLimitMgr struct {
	UserLimitMap *sync.Map
	IpLimitMap   *sync.Map
}

var SecLimitMgrVars = &SecLimitMgr{
	UserLimitMap: new(sync.Map),
	IpLimitMap:   new(sync.Map),
}

var (
	IPBlackMap = make(map[string]bool)
	IDBlackMap = make(map[string]bool)
)

func CheckBlack(userId string, IP string, nowTime int64) (err error){
	// 初始化
	var secIdCount, secIpCount int
	//var limit *Limit
	var limit = new(Limit)
	// 用户黑名单
	_, ok := IDBlackMap[userId]
	if ok {
		// TODO
		err = fmt.Errorf("invalid request")
		log.Printf("user[%v] is block by id black", userId)
		return
	}

	// IP黑名单
	_, ok = IPBlackMap[IP]
	if ok {
		// TODO
		err = fmt.Errorf("invalid request")
		log.Printf("userId[%v] ip[%v] is block by ip black", userId, IP)
		return
	}

	// 用户ID频率控制
	user, ok := SecLimitMgrVars.UserLimitMap.Load(userId)
	if !ok {
		limit = &Limit{
			secLimit: &SecLimit{},
		}
		SecLimitMgrVars.UserLimitMap.Store(userId, limit)
	} else {
		limit = user.(*Limit)
	}
	secIdCount = limit.secLimit.Count(nowTime)

	// 客户端Ip频率控制
	ip, ok := SecLimitMgrVars.IpLimitMap.Load(IP)
	if !ok {
		limit = &Limit{
			secLimit: &SecLimit{},
		}
		SecLimitMgrVars.IpLimitMap.Store(IP,limit)
	} else {
		limit = ip.(*Limit)
	}
	secIpCount = limit.secLimit.Count(nowTime)

	// 设置频率
	if secIdCount > 50000 {
		// TODO
		err = fmt.Errorf("invalid request")
		return
	}

	if secIpCount > 50000 {
		// TODO
		err = fmt.Errorf("invalid request")
		return
	}

	return
}
