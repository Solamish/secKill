package logic

import (
	"time"
)

/**
   测试数据
 */
var (
	ItemMap = make(map[uint]*Item)
)

func InitMap() {
	now := time.Now()
	item := &Item{
		ID:        1,
		Name:      "电影票",
		Total:     100,
		Left:      100,
		PlayTime:  now.Add(time.Hour).Unix(),
		Location:  "cinema",
		Image:     "www.baidu.com",
		BeginTime: now.Add(-time.Hour * 1).Unix(),
		EndTime:   now.Add(time.Minute * 10).Unix(),
		IsSoldOut: false,
		SaleMap:   make(map[string]bool),
		leftCh:    make(chan int),
		sellCh:    make(chan int),
	}
	ItemMap[item.ID] = item
}


