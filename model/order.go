package model

import (
	"log"
	"time"
)

type Order struct {
	ID          uint   `gorm:"primary_key"`
	ProductId   uint   `gorm:"column:product_id"`
	ProductName string `gorm:"column:product_name"` // 商品名字
	RedId       string `gorm:"column:redid"`        // 用户唯一标识
	StuNum      string `gorm:"column:stunum"`
	Status      int    `gorm:"column:status"` // 订单状态
	CreatedTime time.Time `gorm:"column:created_time"`
}

// 创建订单
func (order *Order) CreateOrderList()  {
	err := DB.Create(&order).Error
	if err != nil {
		log.Println("fail to create order, Error:",err)
	}
}

// 通过product_id查订单
func GetOrderListById(redId string, productId uint) (error, Order) {
	orderList := Order{}
	err := DB.Where("redid = ? and product_id = ?", redId, productId).First(&orderList).Error
	if err != nil {
		return err, Order{}
	}
	return nil, orderList
}

// 通过productId查抢购成功的用户
func GetAllOrderLists(productId uint) (error, []Order) {
	var orderLists []Order
	err := DB.Where("product_id = ?", productId).Find(&orderLists).Error
	if err != nil {
		return err, nil
	}
	return nil, orderLists
}
