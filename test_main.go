package main

import (
	"secKill/model"
	"secKill/proxy"
	"secKill/rabbitmq"
	"secKill/redis"
	"secKill/logic"
	"strconv"
	"sync"
	"time"
)

func main() {
	logic.InitMap()
	for _, item := range logic.ItemMap {
			item.Monitor()
		}
	redis.InitRedis()
	model.InitDB()
	rabbitmq.InitMQ()
	go rabbitmq.RabbitMqConsumer.ConsumeSimple(logic.CreateOrder)
	wg := sync.WaitGroup{}
	wg.Add(7000)

	for i:= 0; i < 7000; i++{
		go func(num int) {
			defer wg.Done()
			sec := &proxy.SecRequest{
				ProductId:   1,
				UserId:      strconv.Itoa(num),
				AccessTime:  time.Now(),
				ClientIp:    "127.0.0.1",
				CloseNotify: make(chan bool),
				ResultChan:  make(chan *proxy.SecResult, 1),
			}
			sec.SecKill()
		}(i)
	}
	wg.Wait()
}


