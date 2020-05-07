package logic

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"secKill/rabbitmq"
	"secKill/resps"
	"secKill/rpc"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Item struct {
	ID        uint            // 商品id
	Name      string          // 名字
	Total     int             // 商品总量
	Left      int             // 商品剩余数量
	PlayTime  int64           // 电影放映时间
	Location  string          // 电影放映地点
	Image     string          // 图片
	BeginTime int64           // 秒杀开始时间
	EndTime   int64           // 秒杀结束时间
	IsSoldOut bool            // 是否售罄
	SaleMap   map[string]bool // 商品和用户的映射关系
	Status    int             // 商品状态
	leftCh    chan int
	sellCh    chan int
	done      chan struct{}
	Lock      sync.Mutex

}

var (
	itemMap = make(map[uint]*Item)
	itemMapLock sync.Mutex
)

var (
	SecHandleChan = make(chan *SecRequest, 1024)
	SecWriteChan  = make(chan *SecResult, 1024)
	LocalChan     = make(chan *SecResult, 1024)
	lockChan      = make(chan struct{}, 1)
    closedchan = make(chan struct{})
)

// 初始化库存
func InitStock(total uint, name string, beginTime time.Time, endTime time.Time) {
	cinemaInfos := rpc.GetCinemaInfoByRPC()
	for _, v := range cinemaInfos {
		item := Item{
			ID:        v.Id,
			Name:      v.Name,
			Total:     int(v.Num),
			Left:      int(v.Num),
			PlayTime:  v.TimePlay,
			Location:  v.Location,
			Image:     v.Image,
			BeginTime: v.TimePlay,
			EndTime:   v.TimeOut,
			IsSoldOut: false,
			SaleMap:   make(map[string]bool),
			Status:    resps.SecKillUnStarted,
			leftCh:    make(chan int),
			sellCh:    make(chan int),
		}
		itemMap[v.Id] = &item
	}
}

// TODO
func GetItemByID(itemId uint) (*Item, error) {
	item, ok := itemMap[itemId]
	if !ok {
		return new(Item), errors.New("item do not exists")
	}
	return item, nil
}

func Handle() {
	for {
		req := <-SecHandleChan
		res := &SecResult{
			ProductId: req.ProductId,
			UserId:    req.UserId,
			Code:      0,
		}
		item, err := GetItemByID(req.ProductId)
		if err != nil {
			// 商品不存在
			res.Code = resps.SecKillSoldOut
		}

		isSuccess, uid := item.SecKilling(req.UserId)
		// 重复购买
		if !isSuccess && uid != "" && uid == req.UserId {
			res.Code = resps.SecKillRepeatPurchase
		}
		// 秒杀成功
		if isSuccess {
			res.Code = resps.SecKillSuccess
			//LocalChan <- res
			data, _ := json.Marshal(res)

			go rabbitmq.RabbitMqProducer.PublishSimple(string(data))
		}

		ticker := time.NewTicker(time.Millisecond * 100)
		select {
		case <-ticker.C:
			log.Printf("send to response chan timeout, res : %v", res)
			break
		case SecWriteChan <- res:
		}
	}
}

// 抢购活动状态
func (item *Item) SecKillInfo() (status int) {

	nowTime := time.Now().Unix()
	beginTime := item.BeginTime
	endTime := item.EndTime

	lockChan <- struct{}{}
	defer func() {
		<-lockChan
	}()
	// 抢购尚未开始
	if nowTime-beginTime < 0 {
		item.Status = resps.SecKillUnStarted
		status = resps.SecKillUnStarted
		return
	}

	// 抢购进行时
	if nowTime-beginTime >= 0 {
		// 商品已售罄
		if left := item.GetLeft(); left <= 0 {
			status = resps.SecKillSoldOut
			return
		}
		item.Status = resps.SecKillInProcess
		status = resps.SecKillInProcess
		return
	}

	// 抢购已经结束
	if nowTime-endTime > 0 {
		item.Status = resps.SecKillEnded
		status = resps.SecKillEnded
		return
	}

	return
}

// 秒杀具体逻辑
// 抢购成功返回 true和uid
// 重复抢购则返回 false和uid
func (item *Item) SecKilling(userId string) (isSuccess bool, uid string) {
	item.Lock.Lock()
	defer item.Lock.Unlock()

	// 抢购已经结束，直接返回
	if status := item.SecKillInfo(); status != resps.SecKillInProcess {
		return false, ""
	}

	// 抢购数量限制,一位用户一次只能抢购一件商品
	if _, ok := item.SaleMap[userId]; ok {
		return false, userId
	}

	if _, ok := item.SaleMap[userId]; !ok {
		item.SaleMap[userId] = true
		item.BuyGoods(1)
		logMsg := userId + "_" + strconv.Itoa(int(item.ID))
		writeLog(logMsg, "./log/status.log")
		return true, userId
	}

	return false, ""
}

// 商品数量
func (item *Item) SalesGoods() {
	for {
		select {
		case num := <-item.sellCh:
			if item.Left -= num; item.Left <= 0 {
				item.IsSoldOut = true
			}

		case item.leftCh <- item.Left:
		case <-item.Done():
			return
		}
	}
}

// 保证只有一个goroutine能够操作商品的数量
func (item *Item) Monitor() {
	go item.SalesGoods()
}

// 同步goroutine的关闭
// 目前用于同步Monitor的关闭
func (item *Item) Done() <- chan struct{} {
	item.Lock.Lock()
	if item.done == nil {
		item.done = make(chan struct{})
	}
	d := item.done
	item.Lock.Unlock()
	return d
}

// 获取剩余库存
func (item *Item) GetLeft() int {
	var left int
	left = <-item.leftCh
	return left
}

// 购买商品
func (item *Item) BuyGoods(num int) {
	item.sellCh <- num
}

// 次日零点，商品定时下架
func (item *Item) OffShelve() {
	beginTime := time.Unix(item.BeginTime, 0)
	// 获取第二天时间
	nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品下架的时间
	offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
	timer := time.NewTimer(offShelveTime.Sub(beginTime))
	timer = time.NewTimer(beginTime.Add(5*time.Second).Sub(time.Now()))

	<-timer.C

	if item.done == nil {
		item.done = closedchan
	} else {
		close(item.done)
	}
	// TODO
	delete(ItemMap, item.ID)


}

func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "")
	buf := []byte(content)
	fd.Write(buf)
}

func init() {
	close(closedchan)
	for i := 0; i < 20; i++ {
		go Handle()
	}
}
