package resps

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SecKillSuccess = iota + 10000
	SecKillUnStarted
	SecKillInProcess
	SecKillEnded
	SecKillSoldOut
	SecKillRepeatPurchase
	SecKillProcessTimeOut
	SecKillUserAccessLimited
	ClientClosed
)

type SecKillResponseMsg struct {
	Status int         `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data,omitempty"`
}

var (
	SecKillUnStartedMsg = SecKillResponseMsg{
		Status: SecKillUnStarted,
		Info:   "秒杀活动还未开始",
	}
	SecKillSuccessMsg = SecKillResponseMsg{
		Status: SecKillSuccess,
		Info:   "秒杀成功",
	}
	SecKillFailedMsg = SecKillResponseMsg{
		Status: SecKillSoldOut,
		Info:   "商品已售罄或已下架",
	}
	SecKillEndedMsg = SecKillResponseMsg{
		Status: SecKillEnded,
		Info:   "秒杀活动已经结束,请您下次再来",
	}
	SecKillRepeatPurchaseMsg = SecKillResponseMsg{
		Status: SecKillRepeatPurchase,
		Info:   "不能重复购买",
	}
	SecKillProcessTimeOutMsg = SecKillResponseMsg {
		Status: SecKillProcessTimeOut,
		Info:   "请求超时,请重试",
	}
	SecKillUserAccessLimitedMsg = SecKillResponseMsg{
		Status: SecKillUserAccessLimited,
		Info:   "请求过于频繁",
	}
	ClientClosedMsg = SecKillResponseMsg{
		Status: ClientClosed,
		Info:   "客户端错误,请稍后再试",
	}
)

var (
	ParamError = SecKillResponseMsg{
		Status: 400,
		Info:   "param error",
	}
	AuthorizedError = SecKillResponseMsg{
		Status: 403,
		Info:   "Unauthorized",
	}
	UnKnownError = SecKillResponseMsg{
		Status: 404,
		Info:   "Unknown error",
	}
)

func DefinedResp(c *gin.Context, resp SecKillResponseMsg) {
	c.JSON(http.StatusOK, resp)
}

func GetRespMsgByCode(respCode int) (respMsg SecKillResponseMsg) {
	switch respCode {
	case SecKillSuccess:
		respMsg = SecKillSuccessMsg
	case SecKillUnStarted:
		respMsg = SecKillUnStartedMsg
	case SecKillEnded:
		respMsg = SecKillEndedMsg
	case SecKillSoldOut:
		respMsg = SecKillFailedMsg
	case SecKillRepeatPurchase:
		respMsg = SecKillRepeatPurchaseMsg
	case SecKillProcessTimeOut:
		respMsg = SecKillProcessTimeOutMsg
	case SecKillUserAccessLimited:
		respMsg = SecKillUserAccessLimitedMsg
	case ClientClosed:
		respMsg = ClientClosedMsg
	default:
		respMsg = UnKnownError
	}
	return
}
