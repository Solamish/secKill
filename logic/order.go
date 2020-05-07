package logic

import (
	"crypto/sha1"
	"fmt"
	"secKill/model"
	"secKill/resps"
	"time"
	"unsafe"
)

func CreateOrder(userId string, productId uint) {
	item := new(Item)
	item, _ = GetItemByID(productId)
	redId := CreateRedid(userId)
	order := &model.Order{
		ProductId:   productId,
		ProductName: item.Name,
		RedId:       redId,
		StuNum:      userId,
		Status:      resps.SecKillSuccess,
		CreatedTime: time.Now(),
		}
	order.CreateOrderList()
}

func CreateRedid(stuNum string) string {
	s := sha1.New()
	s.Write(QuickS2B(stuNum))
	b := s.Sum(nil)
	redid := fmt.Sprintf("%x", b)
	return redid
}

func QuickS2B(stuNum string) (b []byte){
	return *(*[]byte)(unsafe.Pointer(&stuNum))
}

 

