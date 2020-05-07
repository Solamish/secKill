package proxy

import (
	"fmt"
	"secKill/proxy/limit"
	"secKill/resps"
	"secKill/logic"
	"sync"
	"time"
)

var (
	SecReqChan = make(chan *SecRequest, 1024)
	UserCon    = &UserConn{UserConnMap: make(map[string]chan *SecResult)}
	lock       = sync.Mutex{}
)

type UserConn struct {
	UserConnMap     map[string]chan *SecResult
	UserConnMapLock sync.Mutex
}

type SecRequest struct {
	ProductId   uint            `form:"product_id" json:"product_id" binding:"required"`
	UserId      string          `form:"user_id" json:"user_id" binding:"required"` // 用户ID
	AccessTime  time.Time       `json:"access_time"`                               // 用户访问接口时间
	ClientIp    string          `json:"client_ip"`                                 // 客户端IP
	CloseNotify <-chan bool     `json:"-"`
	ResultChan  chan *SecResult `json:"-"`
}

type SecResult struct {
	ProductId uint   `json:"product_id"`
	UserId    string `json:"user_id"`
	Code      int    `json:"code"` //状态码
}

func (req *SecRequest) SecKill() (code int) {
	nowTime := req.AccessTime.Unix()
	err := limit.CheckBlack(req.UserId, req.ClientIp, nowTime)
	if err != nil {
		code = resps.SecKillUserAccessLimited
		return
	}

	// 用户有一定概率被限制秒杀
	// TODO 暂时不使用该功能
  	//if isLucky := limit.Probability(0.8); !isLucky {
  	//	code = resps.SecKillUserAccessLimited
  	//	return
	//}

	item, err := logic.GetItemByID(req.ProductId)
	if err != nil {
		code = resps.SecKillSoldOut
		return
	}

	code = item.SecKillInfo()
	if code != resps.SecKillInProcess {
		return
	}

	UserCon.UserConnMapLock.Lock()
	userKey := fmt.Sprintf("%s_%d", req.UserId, req.ProductId)
	UserCon.UserConnMap[userKey] = req.ResultChan
	UserCon.UserConnMapLock.Unlock()
	//将请求送入通道并推入到redis队列当中
	SecReqChan <- req

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		code = resps.SecKillProcessTimeOut
		err = fmt.Errorf("request timeout")
		return
	case <-req.CloseNotify:
		code = resps.ClientClosed
		err = fmt.Errorf("client already closed")
		return
	case result := <-req.ResultChan:
		code = result.Code
		return
	}
}
