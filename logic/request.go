package logic

import "time"

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

 

