package controller

import (
	"github.com/gin-gonic/gin"
	"secKill/proxy"
	"secKill/resps"
	"strings"
	"time"
)

type SecRequest struct {
	ProductId   uint        `form:"product_id" json:"product_id" binding:"required"`
	UserId      string      `form:"user_id" json:"user_id" binding:"required"` // 用户ID
	AccessTime  time.Time   `json:"access_time"`                               // 用户访问接口时间
	ClientIp    string      `json:"client_ip"`                                 // 客户端IP
	CloseNotify <-chan bool `json:"-"`
}

func SecKill(c *gin.Context) {
	var secRequest = SecRequest{}
	if err := c.ShouldBind(&secRequest); err != nil {
		resps.DefinedResp(c, resps.ParamError)
		return
	}
	secRequest.AccessTime = time.Now()
	if len(c.Request.RemoteAddr) > 0 {
		secRequest.ClientIp = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	secRequest.CloseNotify = c.Writer.CloseNotify()
	seq := proxy.SecRequest{
		ProductId:   secRequest.ProductId,
		UserId:      secRequest.UserId,
		AccessTime:  time.Now(),
		ClientIp:    secRequest.ClientIp,
		CloseNotify: secRequest.CloseNotify,
		ResultChan:  make(chan *proxy.SecResult, 1),
	}
	code := seq.SecKill()
	resps.DefinedResp(c, resps.GetRespMsgByCode(code))
}
