package controller

import (
	"github.com/gin-gonic/gin"
	"secKill/resps"
	"secKill/logic"
	"time"
)

type ItemInfo struct {
	ID        uint   `json:"id"`    // 商品id
	Name      string `json:"name"`  // 名字
	Total     int    `json:"total"` // 商品总量
	Left      int    `json:"left"`  // 商品剩余数量
	PlayTime  string `json:"play_time"`
	BeginTime string `json:"begin_time"` // 秒杀开始时间
	EndTime   string `json:"end_time"`   // 秒杀结束时间
	Image     string `json:"image"`
	Location  string `json:"location"`
	IsSoldOut int    `json:"is_sold_out"` // 是否售罄
}

func GetProductInfo(c *gin.Context) {
	itemInfos := make([]ItemInfo, 0)
	if len(logic.ItemMap) < 0 {
		resps.DefinedResp(c, resps.SecKillFailedMsg)
		return
	}
	for _, item := range logic.ItemMap {
		var isSoldOut int
		playTime := time.Unix(item.PlayTime, 0)
		beginTime := time.Unix(item.BeginTime, 0)
		endTime := time.Unix(item.EndTime, 0)
		if item.IsSoldOut == true {
			isSoldOut = 1
		}

		intemInfo := ItemInfo{
			ID:        item.ID,
			Name:      item.Name,
			Total:     item.Total,
			Left:      item.GetLeft(),
			PlayTime:  playTime.Format("2006-01-02 15:04:05"),
			BeginTime: beginTime.Format("2006-01-02 15:04:05"),
			EndTime:   endTime.Format("2006-01-02 15:04:05"),
			Image:     item.Image,
			Location:  item.Location,
			IsSoldOut: isSoldOut,
		}
		itemInfos = append(itemInfos, intemInfo)
	}

	resps.DefinedResp(c, resps.SecKillResponseMsg{
		Status: 200,
		Info:   "success",
		Data:   itemInfos,
	})
}
