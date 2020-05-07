package model

import "time"

type Product struct {
	ID        uint      `gorm:"primary_key"`
	Name      string    `gorm:"column:name"` // 商品名字
	Total     uint      `gorm:"column:total"`       // 商品总数
	Left      uint      `gorm:"column:left"`        // 剩余数量
	IsSoldOut int       `gorm:"column:is_sold_out"`
	StartTime time.Time `gorm:"column:start_time;comment:'抢购开始的时间'"`
	EndTime   time.Time `gorm:"column:end_time;comment:'抢购结束时间'"`
}

// 添加商品
func (p *Product) AddProduct() error {
	return DB.Create(&p).Error
}

// 获取商品信息
func GetProduct(startTime time.Time) (error, Product) {
	product := Product{}
	err := DB.Where("start_time = ?", startTime).First(&product).Error
	if err != nil {
		return err, Product{}
	}
	return nil, product
}
